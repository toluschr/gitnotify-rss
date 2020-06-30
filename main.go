package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"./gh"
)

type Link struct {
	Href string `xml:"href,attr"`
}

type Entry struct {
	Title string `xml:"title"`
	Link Link `xml:"link"`
	Updated string `xml:"updated"`
	Summary struct {
		Type string `xml:"type,attr,omitempty"`
		Text string `xml:",chardata"`
	} `xml:"summary"`
	Category struct {
		Term string `xml:"term,attr"`
	} `xml:"category"`
}

type XMLEntry struct {
	Entry
	XMLName xml.Name `xml:"entry"`
}

type Feed Entry
type XMLFeed struct {
	Feed
	XMLName xml.Name `xml:"feed"`
	Entries []XMLEntry
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("> incoming request")
	param := r.URL.Path[1:]
	auths := strings.Split(param, ":")
	if len(auths) != 2 { return }
	client := gh.NewClient(auths[0], auths[1])
	client.Logger = os.Stdout

	// req, err := http.NewRequest("GET", "https://api.github.com/notifications?all=true", nil)
	notifs, err := client.Notifications_GetAll(false, false, "", "")
	if err != nil { return }

	var feed XMLFeed
	feed.Entries = make([]XMLEntry, len(notifs))
	feed.Title = "Github Notifications"
	feed.Link.Href = "https://github.com/notifications"
	feed.Summary.Text = feed.Title
	if len(notifs) > 0 { feed.Updated = notifs[0].UpdatedAt }

	// This loops through every unread notification, extracts the comment url and makes an api request.
	// The api request gets parsed and the html url and body get extracted. The values will be stored in
	// the new XMLEntry object of the feed.
	for i := range notifs {
		cmt := gh.Comment{Body: "", HTMLURL: notifs[i].Subject.Url}
		curThread, curEntry := &notifs[i], &feed.Entries[i]

		// If the message was read, it should have been cached by the rss reader. No need to update it
		if curThread.Subject.LatestCommentUrl != nil {
			req, err := client.NewRequest("GET", *curThread.Subject.LatestCommentUrl, nil)
			if err != nil { return }
			res, err := client.Perform(req)
			if err != nil { return }
			defer res.Body.Close()
			bytes, err := ioutil.ReadAll(res.Body)
			if err != nil { return }

			json.Unmarshal(bytes, &cmt)
		}

		curEntry.Entry.Title = curThread.Subject.Title + " in " + curThread.Repository.FullName
		curEntry.Category.Term = curThread.Subject.Type
		fmt.Println("> Rendering markdown")
		markdown, err := gh.Markdown_Render(cmt.Body, gh.Markdown_Mode_GFM, curThread.Repository.FullName)
		if err != nil { return }
		curEntry.Entry.Summary.Type = "html"
		curEntry.Entry.Summary.Text = markdown
		curEntry.Link.Href = cmt.HTMLURL
		curEntry.Updated = curThread.UpdatedAt
	}

	// Mark notifications read, errors here may be ignored
	req, err := client.NewRequest("PUT", "https://api.github.com/notifications?last_read_at=" + time.Now().Format(time.RFC3339), nil)
	if err != nil { return }
	res, err := client.Perform(req)
	if err != nil { return }
	defer res.Body.Close()

	// Respond with the rss feed
	response, err := xml.Marshal(feed)
	fmt.Fprint(w, string(response))
}

func main() {
	port := "8092"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	fmt.Println("Server running")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
