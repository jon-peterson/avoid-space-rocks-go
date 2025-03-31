package core

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"time"
)

type Spaceship struct {
	gameobjects.Rigidbody
	Spritesheet  *gameobjects.SpriteSheet
	FuelBurning  bool // Is the user burning fuel to accelerate?
	Alive        bool
	InHyperspace bool
}

var _ gameobjects.Collidable = (*Spaceship)(nil)
var _ gameobjects.Destructible = (*Spaceship)(nil)
var _ gameobjects.GameObject = (*Spaceship)(nil)

func NewSpaceship() Spaceship {
	sheet := gameobjects.LoadSpriteSheet("spaceship.png", 7, 1)
	ship := Spaceship{
		Spritesheet: sheet,
		Rigidbody: gameobjects.Rigidbody{
			MaxVelocity: shipMaxSpeed,
		},
		FuelBurning:  false,
		InHyperspace: false,
	}
	return ship
}

// Update the status of the spaceship
func (s *Spaceship) Update(delta float32) error {
	game := GetGame()
	if s.FuelBurning {
		s.Acceleration = rl.Vector2Scale(s.Rotation, shipFuelBoost)
	} else {
		// Decrease the magnitude of the velocity vector by shipDecaySpeed per second
		s.Acceleration = rl.Vector2{}
		s.Velocity = rl.Vector2Scale(s.Velocity, 1-shipDecaySpeed*delta)
	}
	s.Rigidbody.ApplyPhysics(delta)
	s.Position = game.World.Wraparound(s.Position)
	return nil
}

// Spawn the spaceship at the center of the playfield at the start of level. This function may take a long
// time to return because it waits until the spaceship can spawn at a safe place; call appropriately.
func (s *Spaceship) Spawn() {
	game := GetGame()
	s.Alive = true
	s.Position = rl.Vector2{
		X: game.World.Width / 2,
		Y: game.World.Height / 2,
	}
	s.Velocity = rl.Vector2{}
	s.Acceleration = rl.Vector2{}
	s.Rotation = rl.Vector2{X: 0, Y: -1}

	// Wait until spawning won't make the ship explode immediately
	extendedLocation := gameobjects.ExtendRectangle(s.GetHitbox(), 0.5)
	dangerous := game.World.Objects.IsRectangleOccupied(extendedLocation)
	if dangerous {
		time.Sleep(100 * time.Millisecond)
		dangerous = game.World.Objects.IsRectangleOccupied(extendedLocation)
	}
	game.World.Objects.Add(s)
}

// Draw the spaceship at its current position and rotation
func (s *Spaceship) Draw() error {
	if !s.InHyperspace {
		frame := s.frameIndex()
		return s.Spritesheet.Draw(frame, 0, s.Position, s.Rotation)
	}
	return nil
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
	// Create the starting position of the bullet so it's outside of the hitbox
	hitbox := s.GetHitbox()
	bulletOffset := float32(math.Max(float64(hitbox.Width), float64(hitbox.Height))) / 2
	startPos := rl.Vector2Add(s.Position, rl.Vector2Scale(s.Rotation, bulletOffset))

	b := NewBullet(startPos, s.Rotation, true)
	b.Velocity = rl.Vector2Add(rl.Vector2Scale(s.Rotation, bulletSpeed), s.Velocity)

	game := GetGame()
	game.World.Objects.Add(&b)
	game.EventBus.Publish("spaceship:fire")
}

// EnterHyperspace causes the spaceship to jump to a random location on the playfield
func (s *Spaceship) EnterHyperspace() {
	game := GetGame()
	game.EventBus.Publish("spaceship:enter_hyperspace")
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
	if s.InHyperspace {
		return nil
	}
	game := GetGame()
	s.Alive = false
	// Spawn the pieces flying away
	for i := range 4 {
		piece := NewShrapnel(s.Position, s.Spritesheet, uint16(utils.RndInt32InRange(1000, 2000)), i+3)
		game.World.Objects.Add(&piece)
	}
	// Notify other services
	game.EventBus.Publish("spaceship:destroyed")
	return nil
}
