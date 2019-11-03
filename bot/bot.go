package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/juan-medina/cecibot/config"
	"github.com/juan-medina/cecibot/processor"
	"github.com/juan-medina/cecibot/prototype"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type discordClient interface {
	Open() error
	Close() error
	AddHandler(interface{}) func()
	ChannelMessageSend(channelID string, content string) (*discordgo.Message, error)
}

var errInvalidDiscordClient = errors.New("invalid discord client")

type waitFunc func()

type bot struct {
	cfg     config.Config
	discord discordClient
	prc     processor.Processor
	wait    waitFunc
}

func (b *bot) GetConfig() config.Config {
	return b.cfg
}

func New(cfg config.Config) (prototype.Bot, error) {
	log, _ := zap.NewProduction()
	defer log.Sync()

	bot := &bot{cfg: cfg, prc: *processor.New()}

	log.Info("Creating discord client.")
	var discord, err = discordgo.New("Bot " + cfg.GetToken())

	if err != nil {
		return nil, err
	}

	bot.discord = discord
	bot.wait = bot.waitToSignalClose
	return bot, nil
}

func (b *bot) connect() error {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var err error

	log.Info("Init processor.")
	err = b.prc.Init(b)
	if err != nil {
		return err
	}

	log.Info("Bot is connecting.")
	if b.discord == nil {
		err = errInvalidDiscordClient
		return err
	}

	log.Info("Open connection to discord.")
	err = b.discord.Open()
	if err != nil {
		return err
	}

	log.Info("Bot is connected.")

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

	log.Info("End processor")
	b.prc.End()
	log.Info("Bot disconnected.")
}

func (b bot) waitToSignalClose() {
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

	b.wait()

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

func (b bot) isSelfMessage(m *discordgo.MessageCreate, botUser *discordgo.User) bool {
	return m.Author.ID == botUser.ID
}

func (b bot) removeBotMention(m *discordgo.MessageCreate, botUser *discordgo.User) string {
	text := strings.Replace(m.Content, botUser.Mention()+" ", "", -1)
	return strings.TrimSpace(text)
}

func (b bot) getMessageToBoot(m *discordgo.MessageCreate, botUser *discordgo.User) string {
	for _, element := range m.Mentions {
		if element.ID == botUser.ID {
			return b.removeBotMention(m, botUser)
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
	if !b.isSelfMessage(m, s.State.User) {
		if text := b.getMessageToBoot(m, s.State.User); text != "" {
			if response := b.getResponseToMessage(text, m.Author.ID); response != "" {
				b.replyToMessage(m, response)
			}
		}
	}
}
