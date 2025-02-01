package playfield

import (
	"testing"
)

func TestInitGame(t *testing.T) {
	world := World{} // Assuming World is a struct that can be initialized like this
	game := InitGame(world)

	if game.Lives != 3 {
		t.Errorf("Expected Lives to be 3, got %d", game.Lives)
	}

	if game.Level != 1 {
		t.Errorf("Expected Level to be 1, got %d", game.Level)
	}

	if game.Score != 0 {
		t.Errorf("Expected Score to be 0, got %d", game.Score)
	}
}

func TestGetGame(t *testing.T) {
	world := World{}
	InitGame(world)

	game := GetGame()
	if game == nil {
		t.Fatal("Expected game to be initialized, got nil")
	}
}
