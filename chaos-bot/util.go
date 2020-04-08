package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func sendError(s *discordgo.Session, m *discordgo.MessageCreate, err error) {
	log.Printf("ERROR: %v", err)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("ERROR!"))
}

func sendBlock(s *discordgo.Session, m *discordgo.MessageCreate, text string) (*discordgo.Message, error) {
	return s.ChannelMessageSend(m.ChannelID, "```\n"+text+"\n```")
}

func sendMarkdown(s *discordgo.Session, m *discordgo.MessageCreate, markdown string) (*discordgo.Message, error) {
	return s.ChannelMessageSend(m.ChannelID, "```markdown\n"+markdown+"\n```")
}
