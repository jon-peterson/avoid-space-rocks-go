package playfield

import (
	"testing"
)

func TestInitGame(t *testing.T) {
	game := InitGame(200, 500)

	if game.Lives != 3 {
		t.Errorf("Expected Lives to be 3, got %d", game.Lives)
	}

	if game.Level != 1 {
		t.Errorf("Expected Level to be 1, got %d", game.Level)
	}

	if game.Score != 0 {
		t.Errorf("Expected Score to be 0, got %d", game.Score)
	}

	if game.World.width != 200 {
		t.Errorf("Expected World width to be 200, got %f", game.World.width)
	}

	if game.World.height != 500 {
		t.Errorf("Expected World height to be 500, got %f", game.World.height)
	}
}

func TestGetGame(t *testing.T) {
	InitGame(800, 600)

	game := GetGame()
	if game == nil {
		t.Fatal("Expected game to be initialized, got nil")
	}
}
