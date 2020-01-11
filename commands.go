package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type command struct {
	bot     *tb.Bot
	scores  *scoreRepo
	phrases *phraseRepo
}

func (c *command) setUpCommands() {
	c.bot.Handle("/up", c.up)
	c.bot.Handle("/down", c.down)
	c.bot.Handle("/status", c.status)
	c.bot.Handle("/regras", c.regras)
	c.bot.Handle("/addbirl", c.addBirl)
	c.bot.Handle("/birl", c.birl)
}

func (c *command) regras(m *tb.Message) {
	c.bot.Send(m.Chat, `
		Regras:
		+1 ponto se for na academia (up)
		-1 ponto se comer junk food (down)
	`)
}

func (c *command) up(m *tb.Message) {
	sc := c.scores.getScore(m.Sender.ID)
	if sc == nil {
		c.scores.insertScore(&score{
			ID:          m.Sender.ID,
			Name:        m.Sender.FirstName,
			ActualScore: 0,
		})
	}

	sc.ActualScore = sc.ActualScore + 1
	c.scores.updateScore(sc)

	msg := fmt.Sprintf("É isso aí %s! Continue nesse ritmo.", m.Sender.FirstName)
	c.bot.Send(m.Chat, msg)
}

func (c *command) down(m *tb.Message) {
	sc := c.scores.getScore(m.Sender.ID)
	if sc == nil {
		c.scores.insertScore(&score{
			ID:          m.Sender.ID,
			Name:        m.Sender.FirstName,
			ActualScore: 0,
		})
	}

	sc.ActualScore = sc.ActualScore - 1
	c.scores.updateScore(sc)

	msg := fmt.Sprintf("Que pena %s! Mas o importante é continuar.", m.Sender.FirstName)
	c.bot.Send(m.Chat, msg)
}

func (c *command) status(m *tb.Message) {
	sc := c.scores.getScore(m.Sender.ID)
	if sc == nil {
		c.scores.insertScore(&score{
			ID:          m.Sender.ID,
			Name:        m.Sender.FirstName,
			ActualScore: 0,
		})
	}

	score := c.scores.getAll()
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

	c.bot.Send(m.Chat, strings.Join(lines, "\n"), &tb.SendOptions{ParseMode: tb.ModeMarkdown})
}

func (c *command) addBirl(m *tb.Message) {
	text := strings.Replace(m.Text, "/addbirl", "", 1)
	text = strings.Replace(text, "@cancun_challenge_bot", "", 1)
	text = strings.TrimSpace(text)

	if len(text) == 0 {
		c.bot.Send(m.Chat, "Escreve uma frase, PORRA!")
		return
	}

	c.phrases.insertPhrase(text, "birl")
	c.bot.Send(m.Chat, "Isso aí, frase de monstrão!")
}

func (c *command) birl(m *tb.Message) {
	allPhrases := c.phrases.getAll("birl")

	sendables := []interface{}{}
	for _, phrase := range allPhrases {
		sendables = append(sendables, phrase.Text)
	}

	rand.Seed(time.Now().Unix())
	_, err := c.bot.Send(m.Chat, sendables[rand.Intn(len(sendables))])
	if err != nil {
		log.Printf("Error sending %s", err.Error())
	}
}
