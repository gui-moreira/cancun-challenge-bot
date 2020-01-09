package main

import (
	"log"
	"context"
	"fmt"
	"strings"
	"os"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	database := client.Database("heroku_szwnf11w")

	scores := &scoreRepo { mongo: client, collection: database.Collection("scores")}
	phrases := &phraseRepo { mongo: client, collection: database.Collection("phrases")}

	var (
        port      = os.Getenv("PORT")
        publicURL = os.Getenv("PUBLIC_URL")
        token     = os.Getenv("TELEGRAM_TOKEN")
    )

    webhook := &tb.Webhook{
        Listen:   ":" + port,
        Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
    }

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: webhook,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/up", func(m *tb.Message) {
		sc := scores.getScore(m.Sender.ID)
		if sc == nil {
			scores.insertScore(&score{
				ID: m.Sender.ID,
				Name: m.Sender.FirstName,
				ActualScore: 0,
			})
		}
		
		sc.ActualScore = sc.ActualScore + 1
		scores.updateScore(sc)

		msg := fmt.Sprintf("É isso aí %s! Continue nesse ritmo.", m.Sender.FirstName)
		b.Send(m.Chat, msg)
	})

	b.Handle("/status", func(m *tb.Message) {
		sc := scores.getScore(m.Sender.ID)
		if sc == nil {
			scores.insertScore(&score{
				ID: m.Sender.ID,
				Name: m.Sender.FirstName,
				ActualScore: 0,
			})
		}

		score := scores.getAll()
		maxScore := 0
		minScore := 0
		if len(score) > 0 {
			maxScore = score[0].ActualScore
			minScore = score[len(score)-1].ActualScore
		}

		lines := []string{}
		for _, s := range score {
			l := fmt.Sprintf("%s: %v pontos", s.Name, s.ActualScore)

			if s.ActualScore == maxScore {
				l = fmt.Sprintf("*%s* \xF0\x9F\x92\xAA", l)
			}

			if maxScore != minScore && s.ActualScore == minScore {
				l = fmt.Sprintf("%s \xF0\x9F\x8D\x95", l)
			}
			

			lines = append(lines, l)
		}

		b.Send(m.Chat, strings.Join(lines, "\n"), &tb.SendOptions{ParseMode:tb.ModeMarkdown})
	})

	b.Handle("/down", func(m *tb.Message) {
		sc := scores.getScore(m.Sender.ID)
		if sc == nil {
			scores.insertScore(&score{
				ID: m.Sender.ID,
				Name: m.Sender.FirstName,
				ActualScore: 0,
			})
		}

		sc.ActualScore = sc.ActualScore - 1
		scores.updateScore(sc)

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

	b.Handle("/addbirl", func(m *tb.Message) {
		text := strings.Replace(m.Text, "/addbirl", "", 1)
		text = strings.Replace(text, "@cancun_challenge_bot", "", 1)
		text = strings.TrimSpace(text)

		if len(text) == 0 {
			b.Send(m.Chat, "Escreve uma frase, PORRA!")
			return
		}

		phrases.insertPhrase(text, "birl")
		b.Send(m.Chat, "Isso aí, frase de monstrão!")
	})

	b.Handle("/birl", func(m *tb.Message) {
		allPhrases := phrases.getAll("birl")

		rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
		p := allPhrases[rand.Intn(len(allPhrases))]

		b.Send(m.Chat, p.Text)
	})

	b.Start()
}