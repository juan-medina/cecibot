package processor

import (
	"github.com/juan-medina/cecibot/config"
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

type fakeBot struct {
	cfg config.Config
}

func (f fakeBot) GetConfig() config.Config {
	return f.cfg
}

func (f fakeBot) Run() error {
	return nil
}

func TestDefaultProcessor_ProcessMessage(t *testing.T) {
	cfg := fakeCfg{}
	bot := fakeBot{cfg: cfg}

	type testCase struct {
		name string
		text string
		want string
		user string
	}
	cases := []testCase{
		{
			"ping command",
			"ping",
			"pong!",
			"6789",
		},
		{
			"pong command",
			"pong",
			"ping!",
			"6789",
		},
		{
			"hello command",
			"hello",
			"hello!",
			"6789",
		},
		{
			"hello command by owner",
			"hello",
			"hello master!",
			cfg.GetOwner(),
		},
		{
			"invalid command",
			"zzz",
			"",
			"6789",
		},
	}

	proc := defaultProcessor{}
	_ = proc.Init(bot)

	for _, tt := range cases {
		got := proc.ProcessMessage(tt.text, tt.user)

		t.Run(tt.name, func(t *testing.T) {
			if got != tt.want {
				t.Errorf("processor error got %q, want %q", got, tt.want)
			}
		})

	}
	proc.End()
}
