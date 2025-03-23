package core

import (
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
