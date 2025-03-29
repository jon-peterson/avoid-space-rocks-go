package core

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	gameobjects.Rigidbody
	spritesheet   *gameobjects.SpriteSheet
	isAlive       bool
	isPlayerFired bool
	ageMs         uint16
}

var _ gameobjects.Collidable = (*Bullet)(nil)
var _ gameobjects.GameObject = (*Bullet)(nil)

// NewBullet creates a new bullet with a given position and velocity.
func NewBullet(position, velocity rl.Vector2, isPlayerFired bool) Bullet {
	sheet := gameobjects.LoadSpriteSheet("bullet.png", 1, 1)
	bullet := Bullet{
		spritesheet: sheet,
		Rigidbody: gameobjects.Rigidbody{
			Velocity: velocity,
			Transform: gameobjects.Transform{
				Position: position,
			},
		},
		isPlayerFired: isPlayerFired,
		isAlive:       true,
		ageMs:         0,
	}
	return bullet
}

// Update applies physics to the bullet so it moves per its velocity.
func (b *Bullet) Update() error {
	game := GetGame()
	b.Rigidbody.ApplyPhysics()
	b.Position = game.World.Wraparound(b.Position)
	b.ageMs += uint16(rl.GetFrameTime() * 1000)
	return nil
}

// Draw renders the bullet to the screen.
func (b *Bullet) Draw() error {
	return b.spritesheet.Draw(0, 0, b.Position, b.Rotation)
}

// IsAlive returns true if the bullet is still alive. Always dead after its lifetime.
func (b *Bullet) IsAlive() bool {
	return b.isAlive && b.ageMs < bulletLifetimeMs
}

func (b *Bullet) IsEnemy() bool {
	return false
}

func (b *Bullet) IsPlayerFired() bool {
	return b.isPlayerFired
}

// GetHitbox returns the hitbox of the bullet, used for basic collision detection.
func (b *Bullet) GetHitbox() rl.Rectangle {
	return rl.Rectangle{
		X:      b.Position.X,
		Y:      b.Position.Y,
		Width:  1,
		Height: 1,
	}
}

// OnCollision handles the collision of the bullet with another object.
func (b *Bullet) OnCollision(other gameobjects.Collidable) error {
	if destructible, ok := other.(gameobjects.Destructible); ok {
		if _, ok := other.(*Spaceship); ok {
			// Spaceship bullets don't destroy the spaceship
			if b.isPlayerFired {
				return nil
			}
		}
		if _, ok := other.(*Alien); ok {
			// Alien bullets don't destroy the alien
			if !b.isPlayerFired {
				return nil
			}
		}
		b.isAlive = false
		return destructible.OnDestruction(b.Velocity)
	}
	return nil
}
