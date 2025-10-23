package main

import (
	"strings"
)

const (
	squareUnknown = ":black_large_square:"
	squareAbsent  = ":black_large_square:"
	squareMaybe   = ":yellow_square:"
	squarePresent = ":green_square:"
	squareNoGuess = ":blue_square:"
)

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

// Guess submits a guess to this EwegnSession, adding it to the RevealBoard and GuessBoard, and incrementing the
// RoundNumber. The caller is expected to validate the input, which should be a string of precisely length 5 containing
// only lowercase ASCII letters from 'a' to 'z'. Does not check that the guess is a valid guess string, the caller
// should also do this itself. Returns true if the guess was correct, false otherwise (or if the RoundNumber was > 5).
func (s *EwegnSession) Guess(guess string) bool {
	if s.RoundNumber > 5 {
		return false
	}

	// Green guesses (and setting letters in the GuessBoard)
	for i := 0; i < 5; i++ {
		s.GuessBoard[s.RoundNumber][i] = guess[i]

		if guess[i] == s.Answer[i] {
			s.RevealBoard[s.RoundNumber][i] = revealedPresent
		}
	}

	// Yellow guesses
	for i := 0; i < 5; i++ {
		if s.RevealBoard[s.RoundNumber][i] == revealedPresent {
			continue
		}

		for j := 0; j < 5; j++ {
			if j == i {
				continue
			}

			if guess[i] == s.Answer[j] && s.RevealBoard[s.RoundNumber][j] == revealedUnknown {
				s.RevealBoard[s.RoundNumber][i] = revealedMaybe
				break
			}
		}
	}

	// Absent guesses
	for i := 0; i < 5; i++ {
		if s.RevealBoard[s.RoundNumber][i] == revealedUnknown {
			s.RevealBoard[s.RoundNumber][i] = revealedAbsent
		}
	}

	s.RoundNumber++

	return guess == s.Answer
}

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
