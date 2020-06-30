package gh

import (
	"io"
	"fmt"
	"net/http"
	"encoding/base64"
)

// Client : Client for authenticated Github api requests
// Logger : file the simple debug messages should be written to
type Client struct {
	auth string
	Logger io.Writer
}

// NewClient : Creates a new GitHub api client from user and pass (token)
func NewClient (user, pass string) *Client {
	return &Client{ auth: base64.StdEncoding.EncodeToString([]byte(user + ":" + pass)), Logger: nil }
}

// NewRequest : Just like http.NewRequest, but adds the authentication from the client object
func (c *Client) NewRequest (method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil { return nil, err }
	req.Header.Add("Authorization", "Basic " + c.auth)
	return req, nil
}

func fprintf (writer io.Writer, format string, a interface{}) {
	if writer != nil { fmt.Fprintf(writer, format, a) }
}

func fprintln (writer io.Writer, a interface{}) {
	if writer != nil { fmt.Fprintln(writer, a) }
}

// Perform : Simple wrapper arround client.Do(), but with a debug line
func (c *Client) Perform (req *http.Request) (*http.Response, error) {
	fprintf(c.Logger, "> Requesting %v\n", req.URL.String())
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil { return nil, err }
	return res, nil
}
