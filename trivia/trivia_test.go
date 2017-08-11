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
