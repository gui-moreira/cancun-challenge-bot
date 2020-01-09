package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type scoreRepo struct {
	mongo *mongo.Client
	collection *mongo.Collection
}

func (s *scoreRepo) getAll() []*score {
	options := options.Find()

	// Sort by `_id` field descending
	options.SetSort(bson.D{{"actualscore", -1}})

	filter := bson.D{}
	cur, err := s.collection.Find(context.TODO(), filter, options)

	if err != nil {
		fmt.Printf("Error finding documents %v\n", err)
	}

	var results []*score
	for cur.Next(context.TODO()) {
    
		// create a value into which the single document can be decoded
		var sc score
		err := cur.Decode(&sc)
		if err != nil {
			fmt.Printf("Error decoding documents %v\n", err)
		}
	
		results = append(results, &sc)
	}

	return results
}
func (s *scoreRepo) getScore(ID int) *score {
	var sc score

	filter := bson.D{{"id", ID}}

	err := s.collection.FindOne(context.TODO(), filter).Decode(&sc)
	if err != nil {
		fmt.Printf("Could not get document %v\n", err)
		return nil
	}

	fmt.Printf("Found a single document: %+v\n", sc)
	return &sc
}

func (s *scoreRepo) insertScore(newScore *score) {
	insertResult, err := s.collection.InsertOne(context.TODO(), newScore)
	if err != nil {
		fmt.Printf("Error inserting document %v\n", err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
}

func (s *scoreRepo) updateScore(updatedScore *score) {
	filter := bson.D{{"id", updatedScore.ID}}

	update := bson.D{
		{"$set", bson.D{
			{"actualscore", updatedScore.ActualScore},
		}},
	}

	updateResult, err := s.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Printf("Error updating document %v\n", err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}