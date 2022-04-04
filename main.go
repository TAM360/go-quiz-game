package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/eiannone/keyboard"
	"io"
	"math/rand"
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

func readCsvFile(fileName string, shuffle bool) []QuizQuestion {
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

	if shuffle {
		fmt.Println("Shuffling the list of questions")

		rand.Shuffle(len(quizQuestions), func(i int, j int) {
			quizQuestions[i], quizQuestions[j] = quizQuestions[j], quizQuestions[i]
		})
	}

	return quizQuestions
}

func getUserInput(ch chan string, q QuizQuestion) {
	answer := ""
	fmt.Printf("What is %s equal to: ", q.question)
	fmt.Scan(&answer)

	ch <- answer
}

func quizGame(csvFileName string, quizTime int, shuffle bool) (uint, uint, int) {
	questions := readCsvFile(csvFileName, shuffle)
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

	close(ch)
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
	var shuffleFlag = flag.Bool("shuffle", false, "Shuffle the quiz questions.")

	flag.Parse()

	fmt.Println("========== Welcome Go Quiz Game =========")
	fmt.Println("    Press Enter key to start the game    ")
	fmt.Println("=========================================")

	_, key, err := keyboard.GetSingleKey()

	if err != nil {
		panic(err)
	}

	if fileExists(*fileFlag) && key == keyboard.KeyEnter {
		correctAnswers, wrongAnswers, remainingQuestions = quizGame(*fileFlag, *timerFlag, *shuffleFlag)
		fmt.Println("========== Quiz Stats =========")
		fmt.Printf("Right answers: %v\n", correctAnswers)
		fmt.Printf("Wrong answers: %v\n", wrongAnswers)
		fmt.Printf("Not answered: %v\n", remainingQuestions)

	} else {
		fmt.Println("File doesn't exists")
	}

}
