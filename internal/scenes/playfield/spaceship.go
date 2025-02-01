package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

// Constants for gameplay feel
const (
	rotateSpeed float32 = math.Pi * 3 // 1.5 rotations per second
	maxSpeed    float32 = 20.0        // 20 units per second
	decaySpeed  float32 = 2.0         // 5 units per second slower
	fuelBoost   float32 = 10.0        // 10 units per second
)

type Spaceship struct {
	gameobjects.Rigidbody
	gameobjects.SpriteSheet
	FuelBurning bool // Is the user burning fuel to accelerate?
}

func NewSpaceship() Spaceship {
	sheet, _ := gameobjects.NewSpriteSheet("spaceship.png", 3, 1)
	ship := Spaceship{
		SpriteSheet: sheet,
	}
	// Traditionally starts pointing straight up
	ship.Rotation = rl.Vector2{
		X: 0.0,
		Y: -1.0,
	}
	return ship
}

// Update the status of the spaceship
func (s *Spaceship) Update() error {
	delta := rl.GetFrameTime()
	if s.FuelBurning {
		s.Velocity = rl.Vector2Add(s.Velocity, rl.Vector2Scale(s.Rotation, fuelBoost*delta))
		s.Velocity = rl.Vector2ClampValue(s.Velocity, 0, maxSpeed)
	} else {
		// Decrease the magnitude of the velocity vector by decaySpeed per second
		s.Velocity = rl.Vector2Scale(s.Velocity, 1-decaySpeed*delta)
	}
	// Position is updated after velocity is applied, so that the velocity is applied to the new position
	game := GetGame()
	s.Position = game.World.Wraparound(rl.Vector2Add(s.Position, s.Velocity))
	return nil
}

// Draw the spaceship at its current position and rotation
func (s *Spaceship) Draw() error {
	frame := s.frameIndex()
	return s.SpriteSheet.Draw(frame, 0, s.Position, s.Rotation)
}

// frameIndex returns the index of the correct frame to use in the sprite sheet. There are two
// fuel burning frames, so the index is either 0, 1, or 2.
func (s *Spaceship) frameIndex() int {
	if s.FuelBurning {
		t := rl.GetTime()
		if t-math.Floor(t) < 0.5 {
			return 1
		} else {
			return 2
		}
	}
	return 0
}
