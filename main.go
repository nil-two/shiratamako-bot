package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	port          = os.Getenv("PORT")
	channelID     = os.Getenv("ChannelID")
	channelSecret = os.Getenv("ChannelSecret")
	mid           = os.Getenv("MID")
)

const (
	apiURL = "https://trialbot-api.line.me/v1/events"
)

const (
	responseChannelID = "1383378250"
)

const (
	sendMessageEvent       = "138311608800106203"
	receivedMessageEvent   = "138311609000106303"
	receivedOperationEvent = "138311609100106403"
)

const (
	textMessage     = 1
	imageMessage    = 2
	videoMessage    = 3
	audioMessage    = 4
	locationMessage = 7
	stickerMessage  = 8
	contactMessage  = 10
)

const (
	toUser = 1
)

type Content struct {
	ID          string   `json:"id"`
	ContentType int      `json:"contentType"`
	From        string   `json:"from"`
	ToType      int      `json:"toType"`
	To          []string `json:"to"`
	Text        string   `json:"text"`
}

type Result struct {
	ID          string   `json:"id"`
	EventType   string   `json:"eventType"`
	From        string   `json:"from"`
	FromChannel int      `json:"fromChannel"`
	To          []string `json:"to"`
	ToChannel   int      `json:"toChannel"`
	Content     *Content `json:"content"`
}

type Request struct {
	Result []*Result `json:"result"`
}

type Response struct {
	To        []string `json:"to"`
	ToChannel string   `json:"toChannel"`
	EventType string   `json:"eventType"`
	Content   *Content `json:"content"`
}

func toResponseText(s string) string {
	return fmt.Sprintf("「%s」", s)
}

func responseOne(result *Result) error {
	res := &Response{
		To:        []string{result.Content.From},
		ToChannel: responseChannelID,
		EventType: sendMessageEvent,
		Content: &Content{
			ContentType: textMessage,
			ToType:      toUser,
			Text:        toResponseText(result.Content.Text),
		},
	}

	b, err := json.Marshal(res)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("X-Line-ChannelID", channelID)
	req.Header.Add("X-Line-ChannelSecret", channelSecret)
	req.Header.Add("X-Line-Trusted-User-With-ACL", mid)

	c := &http.Client{}
	if _, err := c.Do(req); err != nil {
		return err
	}
	return nil
}

func response(req *Request) error {
	for _, result := range req.Result {
		if err := responseOne(result); err != nil {
			return err
		}
	}
	return nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a LINE bot.")
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	req := new(Request)

	jd := json.NewDecoder(r.Body)
	if err := jd.Decode(req); err != nil {
		log.Println(err)
		return
	}

	if err := response(req); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/callback", handleCallback)

	addr := ":" + port
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
