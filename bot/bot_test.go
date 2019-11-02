package bot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/juan-medina/cecibot/processor"
	"reflect"
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
	failOnOpen  bool
	failOnClose bool
	failure     bool
	lastError   error
	lastMethod  string
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
	return nil
}

func (f *FakeDiscordClientSpy) ChannelMessageSend(channelID string, content string) (*discordgo.Message, error) {
	return nil, nil
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
}

func Test_bot_connect(t *testing.T) {
	cfg := fakeCfg{}
	discord := &FakeDiscordClientSpy{}
	prc := processor.New(cfg)

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

		want := &FakeDiscordClientSpy{
			failOnOpen:  false,
			failOnClose: false,
			lastMethod:  "Open()",
			lastError:   nil,
			failure:     false,
		}

		if reflect.DeepEqual(discord, want) != true {
			t.Errorf("spy fail want %v, got %v", want, discord)
			return
		}
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

		want := &FakeDiscordClientSpy{
			failOnOpen:  true,
			failOnClose: false,
			lastMethod:  "Open()",
			lastError:   fakeError,
			failure:     true,
		}

		if reflect.DeepEqual(discord, want) != true {
			t.Errorf("spy fail want %v, got %v", want, discord)
			return
		}
	})
}

func Test_bot_disconnect(t *testing.T) {
	cfg := fakeCfg{}
	discord := &FakeDiscordClientSpy{}
	prc := processor.New(cfg)

	b := &bot{
		cfg:     cfg,
		discord: discord,
		prc:     prc,
	}

	t.Run("it should disconnect correctly", func(t *testing.T) {
		b.disconnect()

		want := &FakeDiscordClientSpy{
			failOnOpen:  false,
			failOnClose: false,
			lastMethod:  "Close()",
			lastError:   nil,
			failure:     false,
		}

		if reflect.DeepEqual(discord, want) != true {
			t.Errorf("spy fail want %v, got %v", want, discord)
			return
		}
	})

	t.Run("it should not fail on disconnect failure", func(t *testing.T) {
		discord.failOnClose = true
		b.disconnect()

		want := &FakeDiscordClientSpy{
			failOnOpen:  false,
			failOnClose: true,
			lastMethod:  "Close()",
			lastError:   fakeError,
			failure:     true,
		}

		if reflect.DeepEqual(discord, want) != true {
			t.Errorf("spy fail want %v, got %v", want, discord)
			return
		}
	})
}
