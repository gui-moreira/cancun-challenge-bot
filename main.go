package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"flag"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tb "gopkg.in/tucnak/telebot.v2"
)

const databaseName = "heroku_szwnf11w"

func main() {
	var local bool
	flag.BoolVar(&local, "local", false, "should it be run client side")
	flag.Parse()

	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TELEGRAM_TOKEN")
		mongoURI  = os.Getenv("MONGODB_URI")
	)

	client := setUpMongo(mongoURI)
	bot := setUpBot(local, token, port, publicURL)
	
	scores := &scoreRepo{collection: client.Database(databaseName).Collection("scores")}
	phrases := &phraseRepo{collection: client.Database(databaseName).Collection("phrases")}

	c := &command {
		bot: bot,
		scores: scores,
		phrases: phrases,
	}

	c.setUpCommands()
	bot.Start()
}

func setUpMongo(mongoURI string) (*mongo.Client) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

func setUpBot(local bool, token, port, publicURL string) *tb.Bot {
	var poller tb.Poller
	poller = &tb.LongPoller{Timeout: 10 * time.Second}
	if !local {
		fmt.Println("Setting up webhook")
		poller = &tb.Webhook{
			Listen:   ":" + port,
			Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
		}
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: poller,
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return b
}