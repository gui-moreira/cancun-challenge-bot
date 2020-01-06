package main

import (
	"time"
	"log"
	"context"
	"fmt"
	"strings"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	clientOptions := options.Client().ApplyURI(os.Getenv("CONN_STR"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("heroku_szwnf11w").Collection("scores")

	repo := &scoreRepo { mongo: client, collection: collection}

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/up", func(m *tb.Message) {
		sc := repo.getScore(m.Sender.ID)
		if sc == nil {
			repo.insertScore(&score{
				ID: m.Sender.ID,
				Name: m.Sender.FirstName,
				ActualScore: 0,
			})
		}
		
		sc.ActualScore = sc.ActualScore + 1
		repo.updateScore(sc)

		msg := fmt.Sprintf("É isso aí %s! Continue nesse ritmo.", m.Sender.FirstName)
		b.Send(m.Chat, msg)
	})

	b.Handle("/status", func(m *tb.Message) {
		sc := repo.getScore(m.Sender.ID)
		if sc == nil {
			repo.insertScore(&score{
				ID: m.Sender.ID,
				Name: m.Sender.FirstName,
				ActualScore: 0,
			})
		}

		score := repo.getAll()

		lines := []string{}
		for i, s := range score {
			l := fmt.Sprintf("%s: %v pontos", s.Name, s.ActualScore)
			if i == 0 {
				l = fmt.Sprintf("*%s* \xF0\x9F\x92\xAA", l)
			}
			

			lines = append(lines, l)
		}

		lines[len(lines)-1] = fmt.Sprintf("%s \xF0\x9F\x8D\x95", lines[len(lines)-1])

		b.Send(m.Chat, strings.Join(lines, "\n"), &tb.SendOptions{ParseMode:tb.ModeMarkdown})
	})

	b.Handle("/down", func(m *tb.Message) {
		sc := repo.getScore(m.Sender.ID)
		if sc == nil {
			repo.insertScore(&score{
				ID: m.Sender.ID,
				Name: m.Sender.FirstName,
				ActualScore: 0,
			})
		}

		sc.ActualScore = sc.ActualScore - 1
		repo.updateScore(sc)

		msg := fmt.Sprintf("Que pena %s! Mas o importante é continuar.", m.Sender.FirstName)
		b.Send(m.Chat, msg)
	})

	b.Handle("/regras", func(m *tb.Message) {
		b.Send(m.Chat, `
			Regras:
			+1 ponto se for na academia (up)
			-1 ponto se comer junk food (down)
		`)
	})

	b.Start()
}