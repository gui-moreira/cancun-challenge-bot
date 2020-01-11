package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type phraseRepo struct {
	collection *mongo.Collection
}

func (p *phraseRepo) insertPhrase(text string, t string) {
	phrase := &phrase {
		Text: text,
		Type: t,
	}

	insertResult, err := p.collection.InsertOne(context.TODO(), phrase)
	if err != nil {
		fmt.Printf("Error inserting document %v\n", err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
}

func (p *phraseRepo) getAll(t string) []*phrase {

	filter := bson.D{{"type", t}}
	cur, err := p.collection.Find(context.TODO(), filter)

	if err != nil {
		fmt.Printf("Error finding documents %v\n", err)
	}

	var results []*phrase
	for cur.Next(context.TODO()) {
    
		// create a value into which the single document can be decoded
		var ph phrase
		err := cur.Decode(&ph)
		if err != nil {
			fmt.Printf("Error decoding documents %v\n", err)
		}
	
		results = append(results, &ph)
	}

	return results
}
