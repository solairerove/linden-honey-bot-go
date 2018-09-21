package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/telegram-bot-api.v4"
)

// Song ... tbd
type Song struct {
	Title  string  `json:"title,omitempty"`
	Link   string  `json:"link,omitempty"`
	Author string  `json:"author,omitempty"`
	Album  string  `json:"album,omitempty"`
	Verses []Verse `json:"verses,omitempty"`
}

// Verse ... tbd
type Verse struct {
	Ordinal int    `json:"ord"`
	Data    string `json:"data,omitempty"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI("164424204:AAFlNSTSwpQMOLme2t-0GSkvYJY6Gd8yOfA")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery.Query == "" {
			continue
		}

		log.Printf("Inline query is %s", update.InlineQuery.Query)

		response, err := http.Get("http://127.0.0.1:8000/songs?name=" + update.InlineQuery.Query)

		if err != nil {
			log.Printf("The HTTP request failed with error %s\n", err)
		}

		data, _ := ioutil.ReadAll(response.Body)

		var dat map[string]interface{}

		json.Unmarshal(data, &dat)

		results := make([]interface{}, 0)

		for key, value := range dat {

			response, err := http.Get("http://127.0.0.1:8000/songs/" + key)
			if err != nil {
				log.Printf("The HTTP request failed with error %s\n", err)
			}

			data, _ := ioutil.ReadAll(response.Body)
			var song Song

			unmarshaledSong := json.Unmarshal(data, &song)
			if unmarshaledSong != nil {
				log.Println("Smth wrong with song unmarshaling")
			}

			var verses string
			for _, v := range song.Verses {
				verses = verses + v.Data
			}

			article := tgbotapi.NewInlineQueryResultArticleHTML(key, value.(string), verses)

			results = append(results, article)
		}

		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			CacheTime:     0,
			Results:       results,
		}

		if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
			log.Println(err)
		}
	}
}
