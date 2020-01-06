package main

type score struct {
	ID int `"bson":"id"`
	Name string `"bson":"name"`
	ActualScore int `"bson":"actualscore"`
}

type repository interface {
	getAll()
	getScore(ID int)
	insertScore(s *score)
	updateScore(s *score)
}