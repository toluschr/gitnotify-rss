package gh

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Markdown_Mode int
const (
	Markdown_Mode_GFM Markdown_Mode = iota
	Markdown_Mode_Markdown
)

func markdown_mtos (mode Markdown_Mode) string {
	switch mode {
	case Markdown_Mode_GFM: return "gfm"
	case Markdown_Mode_Markdown: return "markdown"
	}
	return ""
}

func Markdown_Render (text string, mode Markdown_Mode, context string) (string, error) {
	var data struct {
		Text string `json:"text"`
		Mode string `json:"mode"`
		Context string `json:"context"`
	}

	data.Text = text
	data.Mode = markdown_mtos(mode)
	data.Context = context

	post, err := json.Marshal(&data)
	if err != nil { return "", err }
	req, err := http.NewRequest("POST", "https://api.github.com/markdown", bytes.NewBuffer(post))
	if err != nil { return "", err }
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil { return "", err }
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil { return "", err }

	return string(body), nil
}
