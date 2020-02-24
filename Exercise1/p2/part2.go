package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
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
func runQuiz(q quiz, t int) int {
	numQuestions := len(q.questions)
	fmt.Println("Quiz will contain ", numQuestions, " questions")
	fmt.Println("You will have: ", t, "seconds")

	var answer string
	var convertedAnswer int
	timer := time.NewTimer(time.Duration(t) * time.Second)

quiz:

	for i := 0; i < numQuestions; i++ {
		fmt.Print("Question: ", q.questions[i].question, " = ")
		aChan := make(chan int)

		// collect answers as a goroutine
		go func() {
			fmt.Scanln(&answer)
			convertedAnswer, _ = strconv.Atoi(answer)
			aChan <- convertedAnswer
		}()

		// Check the timer before collecting the next answer
		select {
		case <-timer.C:
			fmt.Println("\nTIME IS UP")
			break quiz
		case convertedAnswer = <-aChan:
			q.questions[i].correct = (convertedAnswer == q.questions[i].answer)
		}
	}

	q.calcScore()
	return q.score
}

func main() {

	// get the time in seconds from the t=xxx flag
	timePtr := flag.Int("t", 30, "quiz time limit in seconds")
	flag.Parse()

	// open the file
	quiz := fileHandle()

	// run the quiz and collect the score
	score := runQuiz(*quiz, *timePtr)
	fmt.Println("Your score is: ", score, " out of ", len(quiz.questions))

}
