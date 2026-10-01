package core

import (
	"context"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestAlien_OnCollision(t *testing.T) {
	alien := NewAlien(AlienSmall, rl.NewVector2(100, 100))
	spaceship := NewSpaceship()
	spaceship.Position = rl.NewVector2(100, 100)

	err := alien.OnCollision(&spaceship)
	if err != nil {
		t.Errorf("Unexpected error during collision: %v", err)
	}
	if spaceship.IsAlive() {
		t.Errorf("Expected spaceship to be destroyed after collision")
	}
}

func TestAlien_OnDestruction(t *testing.T) {
	alien := NewAlien(AlienBig, rl.NewVector2(1, 1))
	bulletVelocity := rl.NewVector2(1, 1)

	err := alien.OnDestruction(bulletVelocity)
	if err != nil {
		t.Errorf("Unexpected error during destruction: %v", err)
	}

	if alien.IsAlive() {
		t.Errorf("Expected alien to be destroyed")
	}
}

func TestAlienSpawner_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	// Should exit cleanly and not panic or hang
	AlienSpawner(ctx)
}

func TestNewSpawnedAlien(t *testing.T) {
	game := GetGame()
	pos := rl.NewVector2(50, 50)
	alien := newSpawnedAlien(game, pos)
	if alien == nil {
		t.Fatal("expected alien to not be nil")
	}
	if !alien.IsAlive() {
		t.Errorf("expected new spawned alien to be alive")
	}
	if alien.Position != pos {
		t.Errorf("expected position %v, got %v", pos, alien.Position)
	}
}
