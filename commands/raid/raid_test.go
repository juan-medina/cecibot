package raid

import (
	"github.com/juan-medina/cecibot/command/provider"
	"github.com/juan-medina/cecibot/commands/raid/data/memory"
	"github.com/juan-medina/cecibot/prototype"
	"testing"
)

type fakeProcessor struct {
}

func (f fakeProcessor) ProcessMessage(text string, author string) string {
	return ""
}

func (f fakeProcessor) Init(bot prototype.Bot) error {
	return nil
}

func (f fakeProcessor) End() {
}

func (f fakeProcessor) IsOwner(userId string) bool {
	return userId == "123"
}

func (f fakeProcessor) GetCommandHelp(key string) string {
	return ""
}

func (f fakeProcessor) GetHelp() string {
	return ""
}

func TestNew(t *testing.T) {
	prc := fakeProcessor{}
	got := New(prc)

	gotCommands := got.GetCommands()

	if gotCommands == nil {
		t.Errorf("want commands, got nil")
		return
	}

	gotNumCommands := len(*gotCommands)

	if gotNumCommands != 1 {
		t.Errorf("want 1 command, got %d", gotNumCommands)
		return
	}

	cmd, found := (*gotCommands)["raid"]

	if !found {
		t.Errorf("want raid command, got not found")
		return
	}

	if cmd.Key != "raid" {
		t.Errorf("invalid command key want \"raid\", got %q", cmd.Key)
	}
}

func Test_raidCommands_raid(t *testing.T) {
	prc := fakeProcessor{}
	base := provider.New(prc)
	data := memory.New()
	rc := raidCommands{
		BaseProvider: base,
		data:         data,
	}

	t.Run("should return empty string with not sub command", func(t *testing.T) {
		got := rc.raid([]string{}, "123")
		want := ""
		if got != want {
			t.Errorf("want %q, got %q", want, got)
			return
		}
	})

	t.Run("should return officers", func(t *testing.T) {
		data.AddOfficer("123")
		data.AddOfficer("456")
		got := rc.raid([]string{"officers"}, "123")
		want := "raid officers:\n\t<@123>\n\t<@456>\n"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
		data.DeleteOfficer("123")
		data.DeleteOfficer("456")
	})

	t.Run("should return empty string without officer action", func(t *testing.T) {
		got := rc.raid([]string{"officer"}, "123")
		want := ""
		if got != want {
			t.Errorf("want %q, got %q", want, got)
			return
		}
	})

	t.Run("should delete an officer", func(t *testing.T) {
		data.AddOfficer("123")
		data.AddOfficer("456")

		var got = rc.raid([]string{"officer", "delete", "456"}, "123")
		var want = "officer <@456> deleted"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		got = rc.raid([]string{"officers"}, "123")
		want = "raid officers:\n\t<@123>\n"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		data.DeleteOfficer("123")
		data.DeleteOfficer("456")
	})

	t.Run("should return empty string deleting without id", func(t *testing.T) {
		got := rc.raid([]string{"officer", "delete"}, "123")
		want := ""
		if got != want {
			t.Errorf("want %q, got %q", want, got)
			return
		}
	})

	t.Run("should add an officer", func(t *testing.T) {
		var got = rc.raid([]string{"officer", "add", "456"}, "123")
		var want = "officer <@456> added"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		got = rc.raid([]string{"officers"}, "123")
		want = "raid officers:\n\t<@456>\n"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		data.DeleteOfficer("456")
	})

	t.Run("should return empty string adding without id", func(t *testing.T) {
		got := rc.raid([]string{"officer", "add"}, "123")
		want := ""
		if got != want {
			t.Errorf("want %q, got %q", want, got)
			return
		}
	})

	t.Run("couldn't delete or add officer if is not officer", func(t *testing.T) {
		var got = rc.raid([]string{"officer", "delete", "456"}, "231")
		var want = "this command is for **officers** only"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		got = rc.raid([]string{"officer", "add", "456"}, "231")
		want = "this command is for **officers** only"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})

	t.Run("could delete or add officer if is officer", func(t *testing.T) {
		data.AddOfficer("456")

		var got = rc.raid([]string{"officer", "add", "678"}, "456")
		var want = "officer <@678> added"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		got = rc.raid([]string{"officer", "delete", "678"}, "456")
		want = "officer <@678> deleted"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		got = rc.raid([]string{"officers"}, "123")
		want = "raid officers:\n\t<@456>\n"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		data.DeleteOfficer("456")
	})
}
