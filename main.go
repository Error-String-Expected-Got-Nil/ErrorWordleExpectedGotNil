package main

import (
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
)

func main() {

	// Make sure the file with your token doesn't have a BOM.
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		fmt.Printf("Failed to read bot token: %s\n", err.Error())
		return
	}

	token := string(data)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("Failed to initialize bot: %s\n", err.Error())
		return
	}

	dg.AddHandler(onMessageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Printf("Failed to open connection to Discord: %s\n", err.Error())
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
	case "test":
		test := EwegnSession{
			RevealBoard: [6][5]byte{{revealedAbsent, revealedAbsent, revealedAbsent, revealedMaybe, revealedAbsent}},
			GuessBoard:  [6][5]byte{{'c', 'r', 'a', 't', 'e'}},
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, test.ToString())
	}
}
