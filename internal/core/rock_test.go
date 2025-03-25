package core

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
	"testing"
)

func TestNewRock(t *testing.T) {
	position := rl.NewVector2(100, 100)
	rock := NewRock(RockMedium, position)

	if rock.Position != position {
		t.Errorf("Expected position %v, got %v", position, rock.Position)
	}

	if rock.size != RockMedium {
		t.Errorf("Expected size %v, got %v", RockMedium, rock.size)
	}

	if !rock.isAlive {
		t.Errorf("Expected rock to be alive, got %v", rock.isAlive)
	}
}

func TestRock_GetHitbox(t *testing.T) {
	position := rl.NewVector2(100, 100)
	rock := NewRock(RockMedium, position)
	expectedHitbox := rock.spritesheet.GetRectangle(position)

	hitbox := rock.GetHitbox()
	if hitbox != expectedHitbox {
		t.Errorf("Expected hitbox %v, got %v", expectedHitbox, hitbox)
	}
}

func TestRock_OnSpaceshipCollision(t *testing.T) {
	rock := NewRock(RockMedium, rl.NewVector2(100, 100))
	spaceship := NewSpaceship()
	spaceship.Position = rl.NewVector2(100, 100)

	err := rock.OnCollision(&spaceship)
	if err != nil {
		t.Errorf("Unexpected error during collision: %v", err)
	}

	if spaceship.IsAlive() {
		t.Errorf("Expected spaceship to be destroyed after collision")
	}
}

func TestRock_OnAlienCollision(t *testing.T) {
	rock := NewRock(RockMedium, rl.NewVector2(100, 100))
	alien := NewAlien(AlienBig, rl.NewVector2(100, 100))

	err := rock.OnCollision(&alien)
	if err != nil {
		t.Errorf("Unexpected error during collision: %v", err)
	}

	if alien.IsAlive() {
		t.Errorf("Expected alien to be destroyed after collision")
	}
}

func TestRock_OnDestruction(t *testing.T) {

	// To start assert the game has no small rocks
	game := GetGame()
	if game.World.Objects.Any(
		func(obj gameobjects.GameObject) bool {
			return obj.(*Rock).size == RockSmall
		}) {
		t.Errorf("Expected game.World.Objects to not contain small rocks")
	}

	rock := NewRock(RockMedium, rl.NewVector2(100, 100))
	bulletVelocity := rl.NewVector2(1, 1)

	err := rock.OnDestruction(bulletVelocity)
	if err != nil {
		t.Errorf("Unexpected error during destruction: %v", err)
	}

	if rock.IsAlive() {
		t.Errorf("Expected rock to be destroyed")
	}

	// After a tick there should be some small rocks
	game.World.Objects.Update()
	if !game.World.Objects.Any(
		func(obj gameobjects.GameObject) bool {
			rock, ok := obj.(*Rock)
			if !ok {
				return true // not a rock so we don't care, continue
			}
			return rock.size == RockSmall
		}) {
		t.Errorf("Expected game.World.Objects to contain only small rocks")
	}
}
