package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type question struct {
	question string
	answer   int
	correct  bool
}

type quiz struct {
	questions []question
	score     int
}

// extension method to calculate the score
func (q *quiz) calcScore() {
	numCorrect := 0
	for i := 0; i < len(q.questions); i++ {
		if q.questions[i].correct {
			numCorrect++
		}
	}
	q.score = numCorrect
}

// basic error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// File handle operations
func fileHandle() *quiz {
	// Open the file
	file, err := os.Open("problems.csv")
	defer file.Close()

	check(err)

	quiz := quiz{}
	quiz.score = 0

	reader := csv.NewReader(file)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		answer, err := strconv.Atoi(row[1])
		check(err)
		question := question{row[0], answer, false}
		quiz.questions = append(quiz.questions, question)
	}

	return &quiz
}

// running the quiz questions and returns the score
func runQuiz(q quiz) int {
	numQuestions := len(q.questions)

	var answer string
	var convertedAnswer int
	for i := 0; i < numQuestions; i++ {
		fmt.Print("Question: ", q.questions[i].question, " = ")
		fmt.Scanln(&answer)
		convertedAnswer, _ = strconv.Atoi(answer)
		q.questions[i].correct = (convertedAnswer == q.questions[i].answer)
	}
	q.calcScore()
	return q.score
}

func main() {

	// open the file
	quiz := fileHandle()

	// run the quiz and collect the score
	score := runQuiz(*quiz)
	fmt.Println("Your score is: ", score, " out of ", len(quiz.questions))

}
