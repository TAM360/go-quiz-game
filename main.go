package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type QuizQuestion struct {
	question string
	answer   string
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)

	return err == nil
}

func readCsvFile(fileName string) []QuizQuestion {
	file, fileError := os.Open(fileName)
	quizQuestions := make([]QuizQuestion, 0)

	if fileError != nil {
		fmt.Println(fileError)
	}

	csv := csv.NewReader(file)

	defer file.Close()

	for {
		rec, csvError := csv.Read()
		if csvError == io.EOF {
			fmt.Println("End of CSV file.")
			break
		}

		if csvError != nil {
			fmt.Println(csvError)
		}

		q := QuizQuestion{
			question: rec[0],
			answer:   rec[1],
		}
		quizQuestions = append(quizQuestions, q)

	}

	return quizQuestions
}

func getUserInput(ch chan string, q QuizQuestion) {
	answer := ""
	fmt.Printf("What is %s equal to: ", q.question)
	fmt.Scan(&answer)

	ch <- answer
}

func quizGame(csvFileName string, quizTime int) (uint, uint, int) {
	questions := readCsvFile(csvFileName)
	ch := make(chan string, 1)
	breakLoop := false
	var correctAnswers uint = 0
	var wrongAnswers uint = 0
	var remainingQuestions int = len(questions)

	for i := 0; i < len(questions) && !breakLoop; i++ {
		q := questions[i]

		go getUserInput(ch, q)

		select {
		case answer := <-ch:
			fmt.Printf("Your answer: %s\n", answer)
			if answer != q.answer {
				wrongAnswers++
			} else {
				correctAnswers++
			}
			remainingQuestions -= 1
		case <-time.After(time.Duration(quizTime) * time.Second):
			fmt.Println("Time complete.")
			breakLoop = true
		}
	}

	return correctAnswers, wrongAnswers, remainingQuestions
}

func main() {
	defaultFile := "problems.csv"
	defaultTime := 30
	var correctAnswers uint = 0
	var wrongAnswers uint = 0
	var remainingQuestions int = 0
	var fileFlag = flag.String("file-name", defaultFile, "Name of the CSV file. Default is "+defaultFile+".")
	var timerFlag = flag.Int("timer", defaultTime, "Time of the quiz. Default value is 30 seconds.")

	flag.Parse()

	if fileExists(*fileFlag) {
		correctAnswers, wrongAnswers, remainingQuestions = quizGame(*fileFlag, *timerFlag)
		fmt.Println("Quiz Stats")
		fmt.Printf("Right answers: %v\n", correctAnswers)
		fmt.Printf("Wrong answers: %v\n", wrongAnswers)
		fmt.Printf("Not answered: %v\n", remainingQuestions)

	} else {
		fmt.Println("File doesn't exists")
	}

}
