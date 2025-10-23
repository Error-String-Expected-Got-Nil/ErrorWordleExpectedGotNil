package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	squareUnknown = ":black_large_square:"
	squareAbsent  = ":black_large_square:"
	squareMaybe   = ":yellow_square:"
	squarePresent = ":green_square:"
	squareNoGuess = ":blue_square:"

	validAnswersFile = "valid_answers.txt"
	validGuessesFile = "valid_guesses.txt"

	// Based on the number of lines in each corresponding file
	answersCount = 2315
	guessesCount = 10657
)

var (
	Sessions     = make(map[string]*EwegnSession)
	ValidWords   = make(map[string]struct{}, answersCount+guessesCount)
	ValidAnswers = make([]string, 0, answersCount)
)

func init() {
	answersFile, err := os.Open(validAnswersFile)
	if err != nil {
		fmt.Println("ERROR: Failed to open answers file!", err)
		panic(err)
	}
	defer answersFile.Close()

	answersScanner := bufio.NewScanner(answersFile)
	for answersScanner.Scan() {
		ValidAnswers = append(ValidAnswers, answersScanner.Text())
		ValidWords[answersScanner.Text()] = struct{}{}
	}

	err = answersScanner.Err()
	if err != nil {
		fmt.Println("ERROR: Failed to read answers file!", err)
		panic(err)
	}

	guessesFile, err := os.Open(validGuessesFile)
	if err != nil {
		fmt.Println("ERROR: Failed to open guesses file!", err)
		panic(err)
	}
	defer guessesFile.Close()

	guessesScanner := bufio.NewScanner(guessesFile)
	for guessesScanner.Scan() {
		ValidWords[guessesScanner.Text()] = struct{}{}
	}

	err = guessesScanner.Err()
	if err != nil {
		fmt.Println("ERROR: Failed to read guesses file!", err)
		panic(err)
	}
}

type EwegnSession struct {
	Owner       string     // ID of the user who is playing this session
	RevealBoard [6][5]byte // The colors of spaces that have been guessed
	GuessBoard  [6][5]byte // The letters guessed
	RoundNumber byte       // The current round number (0 to 5)
	Answer      string     // The correct answer for this session
}

const (
	revealedUnknown byte = iota
	revealedAbsent  byte = iota
	revealedMaybe   byte = iota
	revealedPresent byte = iota
)

func (s *EwegnSession) ToString() string {
	str := strings.Builder{}

	for i := 0; i < 6; i++ {
		for j := 0; j < 5; j++ {
			var icon string

			switch s.RevealBoard[i][j] {
			case revealedUnknown:
				icon = squareUnknown
			case revealedAbsent:
				icon = squareAbsent
			case revealedMaybe:
				icon = squareMaybe
			case revealedPresent:
				icon = squarePresent
			}

			str.WriteString(icon)
		}

		// Assumes all bytes in s.GuessBoard are lowercase ASCII letters, or null bytes
		for j := 0; j < 5; j++ {
			c := s.GuessBoard[i][j]

			if c == 0 {
				str.WriteString(squareNoGuess)
				continue
			}

			str.WriteString(":regional_indicator_")
			str.WriteByte(c)
			str.WriteByte(':')
		}

		str.WriteByte('\n')
	}

	return str.String()
}
