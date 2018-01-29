package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/adampointer/go-slackbot"
	"github.com/essentialkaos/slack"
	"golang.org/x/net/context"
)

func startFakeWebServer() {
	port := os.Getenv("PORT")
	if port == "" {
		return
	}
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
		})
		println("Starting listener on port", port)
		http.ListenAndServe(":"+port, nil)
	}()
}

func main() {
	startFakeWebServer()
	bot := slackbot.New(os.Getenv("SLACK_TOKEN"))
	bot.Client.SetDebug(true)

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	toMe.Hear("(?i)(hi|hello|hey).*").MessageHandler(HelloHandler)
	toMe.Hear("(?i)how are you(.*)").MessageHandler(HowAreYouHandler)
	toMe.Hear("(?)attachment").MessageHandler(AttachmentsHandler)

	bot.Run()
}

func HelloHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, "Oh hello!", slackbot.WithTyping)
}

func HowAreYouHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, "A bit tired. You get it? A bit?", slackbot.WithTyping)
}

func AttachmentsHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	txt := "Beep Beep Boop is a ridiculously simple hosting platform for your Slackbots."
	attachment := slack.Attachment{
		Pretext:   "We bring bots to life. :sunglasses: :thumbsup:",
		Title:     "Host, deploy and share your bot in seconds.",
		TitleLink: "https://beepboophq.com/",
		Text:      txt,
		Fallback:  txt,
		ImageURL:  "https://storage.googleapis.com/beepboophq/_assets/bot-1.22f6fb.png",
		Color:     "#7CD197",
	}

	// supports multiple attachments
	attachments := []slack.Attachment{attachment}

	bot.ReplyWithAttachments(evt, attachments, slackbot.WithTyping)
}
