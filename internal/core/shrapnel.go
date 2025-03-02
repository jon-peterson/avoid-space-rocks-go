package core

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Shrapnel struct {
	gameobjects.Rigidbody
	spritesheet   *gameobjects.SpriteSheet
	rotationSpeed float32 // rotations per second
	lifespanMs    uint16  // How long the shrapnel lives in ms
	ageMs         uint16  // How long the shrapnel has been alive
	frame         int
}

var _ gameobjects.GameObject = (*Shrapnel)(nil)

// NewShrapnel creates a new piece of shrapnel with random direction and lifetime
func NewShrapnel(position rl.Vector2, sheet *gameobjects.SpriteSheet, lifespan uint16, frame int) Shrapnel {
	shrapnel := Shrapnel{
		spritesheet: sheet,
		Rigidbody: gameobjects.Rigidbody{
			Velocity: rl.Vector2{
				X: utils.RndFloat32InRange(-shrapnelMaxSpeed, shrapnelMaxSpeed),
				Y: utils.RndFloat32InRange(-shrapnelMaxSpeed, shrapnelMaxSpeed),
			},
			MaxVelocity: shrapnelMaxSpeed,
			Transform: gameobjects.Transform{
				Position: position,
				Rotation: rl.Vector2{
					X: utils.RndFloat32InRange(-1.0, 1.0),
					Y: utils.RndFloat32InRange(-1.0, 1.0),
				},
			},
		},
		rotationSpeed: utils.RndFloat32(shrapnelMaxRotate),
		lifespanMs:    lifespan,
		ageMs:         0,
		frame:         frame,
	}
	// Half of 'em rotate counterclockwise
	if utils.Chance(0.5) {
		shrapnel.rotationSpeed = -shrapnel.rotationSpeed
	}
	return shrapnel
}

// Update applies physics to the bullet so it moves per its velocity.
func (s *Shrapnel) Update() error {
	game := GetGame()
	delta := rl.GetFrameTime()
	s.Rotation = rl.Vector2Rotate(s.Rotation, s.rotationSpeed*delta)
	s.Rigidbody.ApplyPhysics()
	s.Position = game.World.Wraparound(s.Position)
	s.ageMs += uint16(rl.GetFrameTime() * 1000)
	return nil
}

// Draw renders the bullet to the screen.
func (s *Shrapnel) Draw() error {
	return s.spritesheet.Draw(s.frame, 0, s.Position, s.Rotation)
}

// IsAlive returns true if the bullet is still alive. Always dead after its lifetime.
func (s *Shrapnel) IsAlive() bool {
	return s.ageMs < s.lifespanMs
}

func (b *Shrapnel) IsEnemy() bool {
	return false
}
