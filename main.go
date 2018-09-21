package main

import (
	"log"

	"github.com/satori/go.uuid"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("Bot token")
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

		generatedUUID, _ := uuid.NewV4()
		id := generatedUUID.String()

		article := tgbotapi.NewInlineQueryResultArticle(id, "Here are", update.InlineQuery.Query)
		article.Description = update.InlineQuery.Query

		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			CacheTime:     0,
			Results:       []interface{}{article},
		}

		if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
			log.Println(err)
		}
	}
}
