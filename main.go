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

func quizGame(csvFileName string) (uint, uint) {
	questions := readCsvFile(csvFileName)

	var correctAnswers uint = 0
	var wrongAnswers uint = 0
	answer := ""
	start := time.Now()
	duration := time.Since(start).Seconds()

	for i := 0; i < len(questions) && duration <= 10; i++ {
		q := questions[i]
		fmt.Printf("What is %s equal to: ", q.question)
		fmt.Scan(&answer)

		if answer != q.answer {
			wrongAnswers++
		} else {
			correctAnswers++
		}

		duration = time.Since(start).Seconds() // need to execute this in parallel with the quiz.
	}

	return correctAnswers, correctAnswers
}

func main() {
	var defaultFile string = "problems.csv"
	// var defaultTimer int = 5
	var correctAnswers uint = 0
	var wrongAnswers uint = 0
	var fileFlag = flag.String("file-name", defaultFile, "Name of the CSV file. Default is "+defaultFile)
	// var timerFlag = flag.Int("timer", defaultTimer, "Duration of quiz in seconds. Default value is "+strconv.Itoa(defaultTimer))

	flag.Parse()

	if fileExists(*fileFlag) {
		correctAnswers, wrongAnswers = quizGame(*fileFlag)
		fmt.Printf("You answered %v correct answers & %v wrong answers\n", correctAnswers, wrongAnswers)
	} else {
		fmt.Println("File doesn't exists")
	}

}
