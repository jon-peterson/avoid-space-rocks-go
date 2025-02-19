package core

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type Spaceship struct {
	gameobjects.Rigidbody
	Spritesheet *gameobjects.SpriteSheet
	FuelBurning bool // Is the user burning fuel to accelerate?
	Alive       bool
}

var _ gameobjects.Collidable = (*Spaceship)(nil)
var _ gameobjects.GameObject = (*Spaceship)(nil)

func NewSpaceship() Spaceship {
	sheet := gameobjects.LoadSpriteSheet("spaceship.png", 7, 1)
	ship := Spaceship{
		Spritesheet: sheet,
		Rigidbody: gameobjects.Rigidbody{
			MaxVelocity: shipMaxSpeed,
			Transform: gameobjects.Transform{
				Rotation: rl.Vector2{X: 0, Y: -1},
			},
		},
		Alive:       true,
		FuelBurning: false,
	}
	return ship
}

// Update the status of the spaceship
func (s *Spaceship) Update() error {
	delta := rl.GetFrameTime()
	if s.FuelBurning {
		s.Acceleration = rl.Vector2Scale(s.Rotation, shipFuelBoost*delta)
	} else {
		// Decrease the magnitude of the velocity vector by shipDecaySpeed per second
		s.Acceleration = rl.Vector2{}
		s.Velocity = rl.Vector2Scale(s.Velocity, 1-shipDecaySpeed*delta)
	}
	s.Rigidbody.ApplyPhysics()
	// Position is updated after velocity is applied, so that the velocity is applied to the new position
	game := GetGame()
	s.Position = game.World.Wraparound(rl.Vector2Add(s.Position, s.Velocity))
	return nil
}

// Draw the spaceship at its current position and rotation
func (s *Spaceship) Draw() error {
	frame := s.frameIndex()
	return s.Spritesheet.Draw(frame, 0, s.Position, s.Rotation)
}

// RotateLeft rotates the spaceship to the left the standard amount
func (s *Spaceship) RotateLeft() {
	delta := rl.GetFrameTime()
	s.Rotation = rl.Vector2Rotate(s.Rotation, -shipRotateSpeed*delta)
}

// RotateRight rotates the spaceship to the right the standard amount
func (s *Spaceship) RotateRight() {
	delta := rl.GetFrameTime()
	s.Rotation = rl.Vector2Rotate(s.Rotation, shipRotateSpeed*delta)
}

// Fire creates a new bullet with the spaceship's current position and rotation
func (s *Spaceship) Fire() {
	b := NewBullet(s.Position, s.Rotation)
	b.Velocity = rl.Vector2Add(rl.Vector2Scale(s.Rotation, bulletSpeed), s.Velocity)
	game := GetGame()
	game.World.Objects.Add(&b)
	game.EventBus.Publish("spaceship:fire")
}

func (s *Spaceship) IsAlive() bool {
	return s.Alive
}

func (s *Spaceship) IsEnemy() bool {
	return false
}

func (s *Spaceship) OnCollision(_ gameobjects.Collidable) error {
	// Spaceship is always considered the anvil so we always do nothing here
	return nil
}

func (s *Spaceship) GetHitbox() rl.Rectangle {
	return s.Spritesheet.GetRectangle(s.Position)
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

// OnDestruction handles the destruction of the spaceship, causing pieces to fly around.
// This is called by the rock's OnCollision method when it hits this spaceship.
func (s *Spaceship) OnDestruction(_ rl.Vector2) error {
	game := GetGame()
	game.Lives -= 1
	s.Alive = false
	// Spawn the pieces flying away
	for i := range 4 {
		piece := NewShrapnel(s.Position, s.Spritesheet, uint16(utils.RndInt32InRange(1000, 2000)))
		piece.frame = i + 3
		game.World.Objects.Add(&piece)
	}
	// Notify other services
	game.EventBus.Publish("spaceship:destroyed")
	return nil
}
