package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type Bullet struct {
	gameobjects.Rigidbody
	gameobjects.SpriteSheet
	isAlive bool
	born    time.Time
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
		born:    time.Now(),
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

// IsAlive returns true if the bullet is still alive. Always dead after its lifetime.
func (b *Bullet) IsAlive() bool {
	return b.isAlive && time.Since(b.born).Milliseconds() < bulletLifetimeMs
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
		// Both the rock and the bullet die
		rock.isAlive = false
		b.isAlive = false
		// So long as it isn't a tiny rock, spawn more smaller rocks at same loc
		if rock.size > RockTiny {
			game := GetGame()
			for range utils.RndInt32InRange(2, 4) {
				// Spawn a new rock at the same position as the old one but a bit back
				newRock := NewRock(rock.size-1, rl.Vector2Add(rock.Position, b.Velocity))
				game.World.Objects.Add(&newRock)
				// Add a bit of bullet velocity to each new rock so more likely moving away
				newRock.Velocity = rl.Vector2Add(newRock.Velocity, rl.Vector2Scale(b.Velocity, 0.1))
			}
		}
	}
	return nil
}
