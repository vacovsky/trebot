package trivia

import (
	"testing"

	"github.com/go-chat-bot/bot"
)

func Test_loadScores(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "testload",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadScores()
		})
	}
}

func Test_saveScores(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "testsave"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveScores()
		})
	}
}

func Test_renderScores(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "nothing",
		},
	}
	loadScores()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := renderScores()
			if (err != nil) != tt.wantErr {
				t.Errorf("renderScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_trivia(t *testing.T) {
	type args struct {
		command *bot.Cmd
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				command: &bot.Cmd{
					User: &bot.User{
						ID:   "fniaodaw",
						Nick: "fesnjfis",
					},
					Command: "!trivia answer stuff'",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := trivia(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("trivia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("trivia() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deepCheckAnswer(t *testing.T) {
	type args struct {
		providedAnswer string
		realAnswer     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deepCheckAnswer(tt.args.providedAnswer, tt.args.realAnswer); got != tt.want {
				t.Errorf("deepCheckAnswer() = %v, want %v", got, tt.want)
			}
		})
	}
}
