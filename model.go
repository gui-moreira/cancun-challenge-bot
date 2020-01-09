package main

type score struct {
	ID int `"bson":"id"`
	Name string `"bson":"name"`
	ActualScore int `"bson":"actualscore"`
}

type phrase struct {
	Text string
	Type string
}