package trivia

import (
	"testing"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderScores()
			if (err != nil) != tt.wantErr {
				t.Errorf("renderScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("renderScores() = %v, want %v", got, tt.want)
			}
		})
	}
}
