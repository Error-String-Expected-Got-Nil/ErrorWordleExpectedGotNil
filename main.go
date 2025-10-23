package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	tokenFile = "token"
	prefix    = "::"

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

func main() {

	// Make sure the file with your token doesn't have a BOM.
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		fmt.Println("Failed to read bot token:", err.Error())
		return
	}

	token := string(data)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Failed to initialize bot:", err.Error())
		return
	}

	dg.AddHandler(onMessageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("Failed to open connection to Discord:", err.Error())
		return
	}

	fmt.Println("EWEGN bot has started up, press CTRL-C to terminate.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	_ = dg.Close()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))

	if len(args) < 1 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No command given.")
		return
	}

	switch args[0] {

	}
}
