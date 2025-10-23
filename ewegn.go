package main

import "strings"

const (
	squareUnknown  = ":black_large_square:"
	squareAbsent   = ":black_large_square:"
	squareMaybe    = ":yellow_square:"
	squarePresent  = ":green_square:"
	squareNoGuess  = ":blue_square:"
	zeroWidthSpace = "\u200b"
)

var (
	Sessions = make(map[string]*EwegnSession)
)

type EwegnSession struct {
	Owner       string     // ID of the user who is playing this session
	RevealBoard [6][5]byte // The colors of spaces that have been guessed
	GuessBoard  [6][5]byte // The letters guessed
	RoundNumber byte       // The current round number (0 to 5)
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
			// Discord will concatenate regional indicator emotes into country flags if they adjacent. Adding a
			// zero-width space after each of them will prevent this.
			str.WriteString(zeroWidthSpace)
		}

		str.WriteByte('\n')
	}

	return str.String()
}
