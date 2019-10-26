package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/juan-medina/cecibot/config"
	"github.com/juan-medina/cecibot/processor"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Bot interface {
	Run() error
}

var errInvalidDiscordClient = errors.New("invalid discord client")

type bot struct {
	cfg     config.Config
	discord *discordgo.Session
	prc     processor.Processor
}

func New(cfg config.Config) (Bot, error) {
	bot := &bot{cfg: cfg, prc: processor.New(cfg)}
	return bot, nil
}

func (b *bot) connect() error {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Bot is connecting.")

	log.Info("Creating discord client.")
	var discord, err = discordgo.New("Bot " + b.cfg.GetToken())
	if err != nil {
		return err
	}

	log.Info("Open connection to discord.")
	err = discord.Open()
	if err != nil {
		return err
	}

	b.discord = discord
	log.Info("Bot is connected", zap.String("user name", discord.State.User.Username))

	return nil
}

func (b *bot) disconnect() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Bot is disconnecting.")

	var err error

	if b.discord == nil {
		err = errInvalidDiscordClient
	} else {
		log.Info("Closing connection to discord.")
		err = b.discord.Close()
	}

	if err != nil {
		log.Error("Error disconnecting.", zap.Error(err))
		return
	}

	log.Info("Bot disconnected.")
}

func (b bot) waitToSignalClose(){
	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func (b bot) Run() error {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Bot starting.")

	var err error

	err = b.connect()
	if err != nil {
		log.Error("Error connecting bot", zap.Error(err))
		return err
	}

	defer b.disconnect()

	b.discord.AddHandler(b.onChannelMessage)

	log.Info("Bot started.")

	b.waitToSignalClose()

	log.Info("Bot ending.")

	return nil
}

func (b bot) sendMessage(channelID string, text string) {

	log, _ := zap.NewProduction()
	defer log.Sync()

	_, err := b.discord.ChannelMessageSend(channelID, text)

	if err != nil {
		log.Error("Error sending message", zap.Error(err))
		return
	}
}

func (b bot) isSelfMessage(m *discordgo.MessageCreate) bool {
	return m.Author.ID == b.discord.State.User.ID
}

func (b bot) removeBotMention(m *discordgo.MessageCreate) string {
	text := strings.Replace(m.Content, b.discord.State.User.Mention()+" ", "", -1)
	return strings.TrimSpace(text)
}

func (b bot) getMessageToBoot(m *discordgo.MessageCreate) string {
	for _, element := range m.Mentions {
		if element.ID == b.discord.State.User.ID {
			return b.removeBotMention(m)
		}
	}
	return ""
}

func (b bot) replyToMessage(m *discordgo.MessageCreate, text string) {
	b.sendMessage(m.ChannelID, fmt.Sprintf("%s %s", m.Author.Mention(), text))
}

func (b bot) getResponseToMessage(text string, author string) string {
	return b.prc.ProcessMessage(text, author)
}

func (b bot) onChannelMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !b.isSelfMessage(m) {
		if text := b.getMessageToBoot(m); text != "" {
			if response := b.getResponseToMessage(text, m.Author.ID); response != "" {
				b.replyToMessage(m, response)
			}
		}
	}
}
