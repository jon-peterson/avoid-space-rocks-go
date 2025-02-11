package core

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitGame(800, 600)
	code := m.Run()
	os.Exit(code)
}

func TestInitGame(t *testing.T) {
	game := GetGame()

	if game.Lives != 3 {
		t.Errorf("Expected Lives to be 3, got %d", game.Lives)
	}

	if game.Level != 1 {
		t.Errorf("Expected Level to be 1, got %d", game.Level)
	}

	if game.Score != 0 {
		t.Errorf("Expected Score to be 0, got %d", game.Score)
	}

	if game.World.Width != 800 {
		t.Errorf("Expected World width to be 800, got %f", game.World.Width)
	}

	if game.World.Height != 600 {
		t.Errorf("Expected World height to be 600, got %f", game.World.Height)
	}
}

func TestGetGame(t *testing.T) {
	game := GetGame()
	if game == nil {
		t.Fatal("Expected game to be initialized, got nil")
	}
}
