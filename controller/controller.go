package controller

import (
	"fmt"

	"github.com/DSC-Sahmyook/dscbot/api"
	"github.com/adlio/trello"
	"github.com/bwmarrin/discordgo"
)

var Message string = "value"

func Board() string {
	// *trello.Board
	board, err := api.Client.GetBoard("I8850kOn", trello.Defaults())
	if err != nil {
		fmt.Print(err)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		// Handle error
	}

	// cards, err := board.GetCards(trello.Defaults())
	// if err != nil {
	// 	// Handle error
	// }
	fmt.Println("[박기홍] lists 내용 확인")
	for _, item := range lists {
		itemCards := item.Cards
		fmt.Printf("[%s]\n", item.Name)
		for _, card := range itemCards {
			fmt.Println(card.Name)
		}
	}
	return board.Name
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		Board()
		s.ChannelMessageSend(m.ChannelID, Message)
		s.ChannelMessageSend(m.ChannelID, Board())
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
