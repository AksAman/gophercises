package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	csvFileName := flag.String("filename", "./problems.csv", "Name of the csv file to be read")
	shouldShuffle := flag.Bool("shuffle", false, "Shuffle the questions?")
	timerLimit := flag.Int("timelimit", 30, "Time Limit?")

	flag.Parse()

	f, err := os.Open(*csvFileName)
	if err != nil {
		log.Fatalf("Error %v while opening file %v\n", err, csvFileName)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatalf("Error %v while reading csv records from %v", err, csvFileName)
	}

	var (
		correct int = 0
	)

	// shuffle records
	if *shouldShuffle {
		rand.Shuffle(len(records), func(i, j int) { records[i], records[j] = records[j], records[i] })
	}

	fmt.Printf("Score: %v/%v\n", correct, len(records))

	timer := time.NewTimer(time.Duration(*timerLimit) * time.Second)

	var userAnswer string

	for index, record := range records {
		fmt.Printf("Question #%d: %v? ", index, record[0])
		expectedAnswer := record[1]
		expectedAnswer = strings.TrimSpace(expectedAnswer)

		answerChan := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%v", &answer)
			answerChan <- answer
		}()

		select {
		case timerResult := <-timer.C:
			{
				fmt.Printf("Score: %v/%v\n", correct, len(records))
				fmt.Println("Time limit exceeded", timerResult)
				return
			}
		case userAnswer = <-answerChan:
			{
				if userAnswer == expectedAnswer {
					correct++
				}
			}
		}
	}

	fmt.Printf("Score: %v/%v\n", correct, len(records))
}
