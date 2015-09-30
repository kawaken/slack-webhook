package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

type Tag struct {
	Name     string
	HookUrl  string
	Channel  string
	UserName string
	IconUrl  string
}

type Config struct {
	Tags []*Tag
}

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type attachment struct {
	Fallback string `json:"fallback"`
	Pretext  string `json:"pretext"`
	Color    string `json:"color"`

	AuthorIcon string `json:"author_icon"`
	AuthorLink string `json:"author_link"`
	AuthorName string `json:"author_name"`

	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`

	Fields []*field `json:"fields"`

	ImageURL string `json:"image_url"`
	ThumbURL string `json:"thumb_url"`
}

type payload struct {
	Channel     string
	UserName    string        `json:"username"`
	IconUrl     string        `json:"icon_url"`
	Attachments []*attachment `json:"attachments"`
}

func loadConfiguration() Config {

	var err error
	var conf Config
	_, err = toml.DecodeFile("conf.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf

}

func Post(tag *Tag, text string) error {
	p := &payload{
		Channel:  tag.Channel,
		UserName: tag.UserName,
		IconUrl:  tag.IconUrl,
		Attachments: []*attachment{
			&attachment{
				Text: text,
			},
		},
	}
	json, err := json.Marshal(&p)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("payload", string(json))
	_, err = http.PostForm(tag.HookUrl, params)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	optTag := flag.String("tag", "", "slackboard tag name")
	flag.Parse()

	if *optTag == "" {
		log.Fatal("tag is required")
	}

	conf := loadConfiguration()
	var tag *Tag
	for _, t := range conf.Tags {
		if *optTag == t.Name {
			tag = t
		}
	}
	if tag == nil {
		log.Fatal("tag is not matched")
	}

	var text bytes.Buffer
	io.Copy(&text, os.Stdin)

	Post(tag, text.String())
}
