package gh

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Notification struct {
	Subject struct {
		LatestCommentUrl *string `json:"latest_comment_url"`
		Title string `json:"title"`
		Type string `json:"type"`
		Url string `json:"url"`
	} `json:"subject"`
	UpdatedAt string `json:"updated_at"`
	Repository Repository `json:"repository"`
}

func (c *Client) notifications_get (page int, all, participating bool, since, before string) ([]byte, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf(
		"https://api.github.com/notifications?page=%v&all=%v&participating=%v&since=%v&before=%v",
		page, all, participating, since, before), nil)
	if err != nil { return nil, err }
	res, err := c.Perform(req)
	if err != nil { return nil, err }
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil { return nil, err }
	return bytes, nil
}

func (c *Client) Notifications_Get (page int, all, participating bool, since, before string) ([]Notification, error) {
	bytes, err := c.notifications_get(page, all, participating, since, before)
	if err != nil { return nil, err }
	var out []Notification
	err = json.Unmarshal(bytes, &out)
	if err != nil { return nil, err }
	return out, nil
}

// TODO make this respect the header value containing next and prev
func (c *Client) Notifications_GetAll (all, participating bool, since, before string) ([]Notification, error) {
	var out []Notification
	for i := 1; true; i++ {
		cur, err := c.Notifications_Get(i, all, participating, since, before)
		if err != nil { return nil, err }
		if len(cur) == 0 { break }
		out = append(out, cur...)
	}
	return out, nil
}
