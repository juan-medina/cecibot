package bot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/juan-medina/cecibot/prototype"
	"testing"
)

type fakeCfg struct {
}

func (f fakeCfg) GetOwner() string {
	return "12345"
}

func (f fakeCfg) GetToken() string {
	return "12345"
}

var fakeError = errors.New("fake error")

type FakeDiscordClientSpy struct {
	failOnOpen               bool
	failOnClose              bool
	failOnChannelMessageSend bool
	failOnAddHandler         bool
	failure                  bool
	lastError                error
	lastMethod               string
	lastMessage              string
	lastChannelTo            string
}

func (f *FakeDiscordClientSpy) recordError(method string, err error) error {
	f.failure = true
	f.lastError = err
	f.lastMethod = method
	return err
}

func (f *FakeDiscordClientSpy) recordSuccess(method string) {
	f.failure = false
	f.lastError = nil
	f.lastMethod = method
}

func (f *FakeDiscordClientSpy) Open() error {
	if f.failOnOpen {
		return f.recordError("Open()", fakeError)
	}
	f.recordSuccess("Open()")
	return nil
}

func (f *FakeDiscordClientSpy) Close() error {
	if f.failOnClose {
		return f.recordError("Close()", fakeError)
	}
	f.recordSuccess("Close()")
	return nil
}

func (f *FakeDiscordClientSpy) AddHandler(interface{}) func() {
	if f.failOnAddHandler {
		_ = f.recordError("AddHandler()", fakeError)
		return nil
	}
	f.recordSuccess("AddHandler()")
	return nil
}

func (f *FakeDiscordClientSpy) ChannelMessageSend(channelID string, content string) (*discordgo.Message, error) {
	if f.failOnChannelMessageSend {
		return nil, f.recordError("ChannelMessageSend()", fakeError)
	}
	f.recordSuccess("ChannelMessageSend()")
	f.lastMessage = content
	f.lastChannelTo = channelID
	return nil, nil
}

func assertSpySuccess(t *testing.T, spy *FakeDiscordClientSpy, method string) bool {
	t.Helper()
	if method != spy.lastMethod {
		t.Errorf("want spy last method to be %q, got %q", method, spy.lastMethod)
		return false
	}
	if spy.failure != false {
		t.Errorf("want spy sucess what was failure")
		return false
	}

	return true
}

func assertSpyFailure(t *testing.T, spy *FakeDiscordClientSpy, method string, err error) bool {
	t.Helper()
	if method != spy.lastMethod {
		t.Errorf("want spy last method to be %q, got %q", method, spy.lastMethod)
		return false
	}
	if spy.failure != true {
		t.Errorf("want spy failure but was sucess")
		return false
	}
	if spy.lastError != err {
		t.Errorf("want spy last error to be %q, got %q", err, spy.lastError)
		return false
	}

	return true
}

type fakeProcessor struct {
	failOnInit bool
}

func (f *fakeProcessor) IsOwner(userId string) bool {
	return false
}

func (f *fakeProcessor) GetCommandHelp(key string) string {
	return ""
}

func (f *fakeProcessor) GetHelp() string {
	return ""
}

func (f *fakeProcessor) Init(bot prototype.Bot) error {
	if f.failOnInit {
		return fakeError
	}
	return nil
}

func (f fakeProcessor) End() {
}

func (f fakeProcessor) ProcessMessage(text string, author string) string {
	return author + " told me : " + text
}

func TestNew(t *testing.T) {
	cfg := fakeCfg{}
	got, err := New(cfg)

	if err != nil {
		t.Errorf("want not error, got %v", err)
		return
	}

	if got == nil {
		t.Errorf("want new bot, got nil")
		return
	}

	if got.GetConfig() != cfg {
		t.Errorf("want config %v, got %v", cfg, got.GetConfig())
		return
	}
}

func Test_bot_connect(t *testing.T) {
	cfg := fakeCfg{}
	discord := &FakeDiscordClientSpy{}
	prc := &fakeProcessor{}

	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	t.Run("it should connect correctly", func(t *testing.T) {
		err := b.connect()

		if err != nil {
			t.Errorf("want not error, got %v", err)
			return
		}

		assertSpySuccess(t, discord, "Open()")
	})

	t.Run("it should fail on connect", func(t *testing.T) {
		discord.failOnOpen = true
		err := b.connect()

		if err == nil {
			t.Errorf("want error, got nil")
			return
		}

		if err != fakeError {
			t.Errorf("want fake error, got %v", err)
			return
		}

		assertSpyFailure(t, discord, "Open()", fakeError)
	})

	t.Run("it should fail with failing processor", func(t *testing.T) {
		prc.failOnInit = true
		err := b.connect()

		if err != fakeError {
			t.Errorf("want fake error, got %v", err)
			return
		}
		prc.failOnInit = false
	})

	t.Run("it should fail without client", func(t *testing.T) {
		b.discord = nil
		err := b.connect()

		if err != errInvalidDiscordClient {
			t.Errorf("want invalid discord client, got %v", err)
			return
		}
	})
}

func Test_bot_disconnect(t *testing.T) {
	cfg := fakeCfg{}
	discord := &FakeDiscordClientSpy{}
	prc := &fakeProcessor{}

	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	t.Run("it should disconnect correctly", func(t *testing.T) {
		b.disconnect()

		assertSpySuccess(t, discord, "Close()")
	})

	t.Run("it should not fail on disconnect failure", func(t *testing.T) {
		discord.failOnClose = true
		b.disconnect()

		assertSpyFailure(t, discord, "Close()", fakeError)
	})

	t.Run("it should not fail without client", func(t *testing.T) {
		b.discord = nil
		b.disconnect()
	})
}

func Test_bot_sendMessage(t *testing.T) {
	cfg := fakeCfg{}
	discord := &FakeDiscordClientSpy{}
	prc := &fakeProcessor{}

	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	t.Run("it should send message correctly", func(t *testing.T) {
		b.sendMessage("chanel", "text")

		assertSpySuccess(t, discord, "ChannelMessageSend()")
	})

	t.Run("it should fail sending message"+
		"e", func(t *testing.T) {
		discord.failOnChannelMessageSend = true
		b.sendMessage("chanel", "text")

		assertSpyFailure(t, discord, "ChannelMessageSend()", fakeError)
	})
}

func Test_bot_Run(t *testing.T) {
	noop := func() {}
	cfg := fakeCfg{}
	prc := &fakeProcessor{}

	t.Run("it should not fail", func(t *testing.T) {
		discord := &FakeDiscordClientSpy{}

		b := &bot{
			cfg:     cfg,
			discord: discord,
			prc:     prc,
		}

		b.wait = noop
		err := b.Run()

		if err != nil {
			t.Errorf("want not error, got %v", err)
			return
		}
	})

	t.Run("it should fail on failure on open", func(t *testing.T) {
		discord := &FakeDiscordClientSpy{}
		discord.failOnOpen = true
		b := &bot{
			cfg:     cfg,
			discord: discord,
			prc:     prc,
		}

		b.wait = noop
		err := b.Run()

		if err != fakeError {
			t.Errorf("want fake error, got %v", err)
			return
		}

		assertSpyFailure(t, discord, "Open()", fakeError)
	})

	t.Run("it should not fail on failure on close", func(t *testing.T) {
		discord := &FakeDiscordClientSpy{}
		discord.failOnClose = true
		b := &bot{
			cfg:     cfg,
			discord: discord,
			prc:     prc,
		}

		b.wait = noop
		err := b.Run()

		if err != nil {
			t.Errorf("want not error, got %v", err)
			return
		}

		assertSpyFailure(t, discord, "Close()", fakeError)
	})

	t.Run("it should not fail on failure on addHandler", func(t *testing.T) {
		discord := &FakeDiscordClientSpy{}
		discord.failOnAddHandler = true
		b := &bot{
			cfg:     cfg,
			discord: discord,
			prc:     prc,
		}

		b.wait = noop
		err := b.Run()

		if err != nil {
			t.Errorf("want not error, got %v", err)
			return
		}

		assertSpySuccess(t, discord, "Close()")
	})
}

func Test_bot_isSelfMessage(t *testing.T) {

	cfg := fakeCfg{}
	prc := &fakeProcessor{}

	discord := &FakeDiscordClientSpy{}
	discord.failOnClose = true
	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	u := &discordgo.User{ID: "123"}
	t.Run("should get a self message", func(t *testing.T) {
		m := &discordgo.MessageCreate{
			Message: &discordgo.Message{Author: u},
		}
		got := b.isSelfMessage(m, u)
		if got != true {
			t.Errorf("is should be self message got %v", got)
		}
	})

	t.Run("should not get a self message", func(t *testing.T) {
		m := &discordgo.MessageCreate{
			Message: &discordgo.Message{Author: &discordgo.User{ID: "456"}},
		}
		got := b.isSelfMessage(m, u)
		if got == true {
			t.Errorf("is should not be self message got %v", got)
		}
	})
}

func Test_bot_removeBotMention(t *testing.T) {

	cfg := fakeCfg{}
	prc := &fakeProcessor{}

	discord := &FakeDiscordClientSpy{}
	discord.failOnClose = true
	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	botUser := &discordgo.User{ID: "123"}

	t.Run("we should remove the mention", func(t *testing.T) {
		m := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Author:  &discordgo.User{ID: "456"},
				Content: "<@123> this is a message",
			},
		}

		got := b.removeBotMention(m, botUser)
		want := "this is a message"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("we should not remove the mention", func(t *testing.T) {
		m := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Author:  &discordgo.User{ID: "456"},
				Content: "<@456> this is a another",
			},
		}

		got := b.removeBotMention(m, botUser)
		want := "<@456> this is a another"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("there is not mention", func(t *testing.T) {
		m := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Author:  &discordgo.User{ID: "456"},
				Content: "there is no mention",
			},
		}

		got := b.removeBotMention(m, botUser)
		want := "there is no mention"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}

func Test_bot_getMessageToBoot(t *testing.T) {

	cfg := fakeCfg{}
	prc := &fakeProcessor{}

	discord := &FakeDiscordClientSpy{}
	discord.failOnClose = true
	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	botUser := &discordgo.User{ID: "123"}

	t.Run("we should get the message in a mention", func(t *testing.T) {
		m := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Author:   &discordgo.User{ID: "456"},
				Content:  "<@123> this is a message",
				Mentions: []*discordgo.User{botUser},
			},
		}

		got := b.getMessageToBoot(m, botUser)
		want := "this is a message"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("we should not get the message without mention", func(t *testing.T) {
		m := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Author:   &discordgo.User{ID: "456"},
				Content:  "this is a message",
				Mentions: []*discordgo.User{},
			},
		}

		got := b.getMessageToBoot(m, botUser)
		want := ""

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}

func Test_bot_replyToMessage(t *testing.T) {

	cfg := fakeCfg{}
	prc := &fakeProcessor{}

	discord := &FakeDiscordClientSpy{}
	discord.failOnClose = true
	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	m := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: "chanel1",
			Author:    &discordgo.User{ID: "456"},
			Content:   "this is a message",
			Mentions:  []*discordgo.User{},
		},
	}

	b.replyToMessage(m, "hello world")

	wantChannel := "chanel1"
	gotChannel := discord.lastChannelTo
	if wantChannel != gotChannel {
		t.Errorf("want message reply to %q, got %q", wantChannel, gotChannel)
	}

	wantMessage := "<@456> hello world"
	gotMessage := discord.lastMessage
	if wantMessage != gotMessage {
		t.Errorf("want message %q, got %q", wantMessage, gotMessage)
	}
}

func Test_getResponseToMessage(t *testing.T) {
	cfg := fakeCfg{}
	discord := &FakeDiscordClientSpy{}
	prc := &fakeProcessor{}

	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	got := b.getResponseToMessage("hello", "user1")
	want := "user1 told me : hello"

	if got != want {
		t.Errorf("want message %q, got %q", want, got)
	}
}

func Test_bot_onChannelMessage(t *testing.T) {

	cfg := fakeCfg{}
	prc := &fakeProcessor{}

	discord := &FakeDiscordClientSpy{}
	discord.failOnClose = true
	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	botUser := &discordgo.User{ID: "123"}
	sta := discordgo.NewState()
	sta.User = botUser
	ses := &discordgo.Session{State: sta}

	m := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: "chanel1",
			Author:    &discordgo.User{ID: "456"},
			Content:   "<@123> this is a message",
			Mentions:  []*discordgo.User{botUser},
		},
	}

	b.onChannelMessage(ses, m)

	wantChannel := "chanel1"
	gotChannel := discord.lastChannelTo
	if wantChannel != gotChannel {
		t.Errorf("want message reply to %q, got %q", wantChannel, gotChannel)
	}

	wantMessage := "<@456> 456 told me : this is a message"
	gotMessage := discord.lastMessage
	if wantMessage != gotMessage {
		t.Errorf("want message %q, got %q", wantMessage, gotMessage)
	}
}
