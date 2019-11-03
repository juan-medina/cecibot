package processor

import (
	"github.com/juan-medina/cecibot/config"
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

	proc := *New()
	_ = proc.Init(bot)

	help := proc.(*processorImpl).help

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
			"help command",
			"help",
			help,
			"6789",
		},
		{
			"help a concrete command",
			"help ping",
			"Command **ping** : \nThis is a test command for the *bot* that will reply with a pong message",
			"6789",
		},
		{
			"help a unknown command",
			"help zzz",
			"Unknown command in help. " + help,
			"6789",
		},
		{
			"invalid command",
			"zzz",
			"Unknown command. " + help,
			"6789",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := proc.ProcessMessage(tt.text, tt.user)
			if got != tt.want {
				t.Errorf("processor error want %q, got %q", tt.want, got)
			}
		})

	}
	proc.End()
}

func TestDefaultProcessor_parseCommand(t *testing.T) {
	proc := processorImpl{}

	type testCase struct {
		name string
		text string
		cmd  string
		args []string
	}

	cases := []testCase{
		{
			name: "empty command",
			text: "",
			cmd:  "",
			args: []string{},
		},
		{
			name: "one word command",
			text: "hello",
			cmd:  "hello",
			args: []string{},
		},
		{
			name: "two words commands",
			text: "hello world",
			cmd:  "hello",
			args: []string{"world"},
		},
		{
			name: "thee words commands",
			text: "hello world command",
			cmd:  "hello",
			args: []string{"world", "command"},
		},
		{
			name: "command with single quotes",
			text: "hello 'big world'",
			cmd:  "hello",
			args: []string{"big world"},
		},
		{
			name: "command with double quotes",
			text: "hello \"big world\"",
			cmd:  "hello",
			args: []string{"big world"},
		},
		{
			name: "command with mixed quotes",
			text: "hello \"big world\" 'with some more quotes'",
			cmd:  "hello",
			args: []string{"big world", "with some more quotes"},
		},
		{
			name: "command with broken quotes",
			text: "hello 'big world",
			cmd:  "",
			args: []string{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotArgs := proc.parseCommand(tt.text)

			if gotCmd != tt.cmd {
				t.Errorf("want %q, got %q", tt.cmd, gotCmd)
				return
			}

			if reflect.DeepEqual(gotArgs, tt.args) == false {
				t.Errorf("want %v, got %v", tt.args, gotArgs)
				return
			}
		})
	}

}
