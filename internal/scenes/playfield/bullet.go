package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type Bullet struct {
	gameobjects.Rigidbody
	gameobjects.SpriteSheet
	born time.Time
}

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
		born: time.Now(),
	}
	return bullet
}

// Update applies physics to the bullet so it moves per its velocity.
func (b *Bullet) Update() error {
	game := GetGame()
	b.Rigidbody.ApplyPhysics()
	b.Position = game.World.Wraparound(b.Position)
	return nil
}

// Draw renders the bullet to the screen.
func (b *Bullet) Draw() error {
	return b.SpriteSheet.Draw(0, 0, b.Position, b.Rotation)
}

func (b *Bullet) IsAlive() bool {
	return time.Since(b.born).Milliseconds() < bulletLifetimeMs
}
