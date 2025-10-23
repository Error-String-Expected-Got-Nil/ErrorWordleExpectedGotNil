package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	tokenFile = "token"
	prefix    = "::"

	helpMessage = "This is Error: Wordle Expected, Got Nil (EWEGN), a Wordle-clone Discord bot made as an excuse to learn the " +
		"programming language Go.\n\n" +
		"This bot uses `::` as a command prefix. Commands are:\n" +
		"- `::help` Displays this message.\n" +
		"- `::play` Starts a new game.\n" +
		"- `::quit` Ends an active game.\n" +
		"- `::guess <word>` Guess 'word' in your current game. Must be a 5-letter word of only 'a' to 'z'.\n" +
		"- `::view` Prints your current game board."

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

	send := func(str string) {
		_, _ = s.ChannelMessageSend(m.ChannelID, str)
	}

	if len(args) < 1 {
		send("No command given.")
		return
	}

	switch args[0] {
	case "help":
		send(helpMessage)

	case "play":
		if _, ok := Sessions[m.Author.ID]; ok {
			send("You already have an active game.")
			return
		}

		session := new(EwegnSession)
		session.Owner = m.Author.ID
		session.Answer = ValidAnswers[rand.IntN(len(ValidAnswers))]
		Sessions[m.Author.ID] = session

		send("New game started.")
		send(session.ToString())

	case "quit":
		session, ok := Sessions[m.Author.ID]
		if !ok {
			send("You don't have an active game.")
			return
		}

		send("Your game has been ended. The answer was: `" + session.Answer + "`")
		delete(Sessions, m.Author.ID)

	case "view":
		session, ok := Sessions[m.Author.ID]
		if !ok {
			send("You don't have an active game.")
			return
		}

		send(session.ToString())

	default:
		send("Unrecognized command `" + args[0] + "`")
	}
}
