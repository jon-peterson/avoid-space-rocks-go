package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	gameobjects.Rigidbody
	gameobjects.SpriteSheet
	isAlive bool
	ageMs   uint16
}

var _ gameobjects.Collidable = (*Bullet)(nil)
var _ gameobjects.GameObject = (*Bullet)(nil)

// NewBullet creates a new bullet with a given position and velocity.
func NewBullet(position, velocity rl.Vector2) Bullet {
	sheet, _ := gameobjects.NewSpriteSheet("bullet.png", 1, 1)
	bullet := Bullet{
		SpriteSheet: sheet,
		Rigidbody: gameobjects.Rigidbody{
			Velocity:    velocity,
			MaxVelocity: bulletMaxSpeed,
			Transform: gameobjects.Transform{
				Position: position,
			},
		},
		isAlive: true,
		ageMs:   0,
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
	return b.SpriteSheet.Draw(0, 0, b.Position, b.Rotation)
}

// IsAlive returns true if the bullet is still alive. Always dead after its lifetime.
func (b *Bullet) IsAlive() bool {
	return b.isAlive && b.ageMs < bulletLifetimeMs
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
	// Bullets can only destroy rocks
	rock, ok := other.(*Rock)
	if ok {
		b.isAlive = false
		return rock.OnDestruction(b.Velocity)
	}
	return nil
}
