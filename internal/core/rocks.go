package core

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RockSize int

const (
	RockTiny RockSize = iota
	RockSmall
	RockMedium
	RockBig
)

// Create a constant array of four string elements
var rockSpriteFile = []string{
	"rock_tiny.png",
	"rock_small.png",
	"rock_medium.png",
	"rock_big.png",
}

// Rock is a game object that has a consistent rotation speed and constant velocity.
type Rock struct {
	gameobjects.Rigidbody
	gameobjects.SpriteSheet
	rotationSpeed float32 // rotations per second
	isAlive       bool
	size          RockSize
}

var _ gameobjects.Collidable = (*Rock)(nil)
var _ gameobjects.GameObject = (*Rock)(nil)

func NewRock(size RockSize, position rl.Vector2) Rock {
	sheet, _ := gameobjects.LoadSpriteSheet(rockSpriteFile[size], 1, 1)
	rock := Rock{
		SpriteSheet: *sheet,
		Rigidbody: gameobjects.Rigidbody{
			Transform: gameobjects.Transform{
				Position: position,
				Rotation: rl.Vector2{X: 1, Y: 0},
			},
		},
		rotationSpeed: utils.RndFloat32(rockMaxRotate) / 4,
		isAlive:       true,
		size:          size,
	}
	// Half of 'em rotate counterclockwise
	if utils.Chance(0.5) {
		rock.rotationSpeed = -rock.rotationSpeed
	}
	// Randomize the speed and direction
	rock.MaxVelocity = rockMaxSpeed / float32(size+2)
	rock.Velocity = rl.Vector2{
		X: utils.RndFloat32InRange(-rock.MaxVelocity, rock.MaxVelocity),
		Y: utils.RndFloat32InRange(-rock.MaxVelocity, rock.MaxVelocity),
	}
	return rock
}

// Update applies physics to the rock so it moves per its velocity and rotation speed.
func (r *Rock) Update() error {
	game := GetGame()
	delta := rl.GetFrameTime()
	r.Rotation = rl.Vector2Rotate(r.Rotation, r.rotationSpeed*delta)
	r.Rigidbody.ApplyPhysics()
	r.Position = game.World.Wraparound(r.Position)
	return nil
}

// Draw renders the rock to the screen.
func (r *Rock) Draw() error {
	return r.SpriteSheet.Draw(0, 0, r.Position, r.Rotation)
}

// IsAlive returns whether the rock is alive or not.
func (r *Rock) IsAlive() bool {
	return r.isAlive
}

// GetHitbox returns the hitbox of the rock, used for basic collision detection.
func (r *Rock) GetHitbox() rl.Rectangle {
	return r.SpriteSheet.GetRectangle(r.Position)
}

// OnCollision handles the collision with another Collidable object.
func (r *Rock) OnCollision(_ gameobjects.Collidable) error {
	// TODO: Check for collision with the spaceship
	return nil
}

// OnDestruction handles the destruction of the rock, spawning smaller rocks if applicable.
// This is called by the bullet's OnCollision method when it hits this rock.
func (r *Rock) OnDestruction(bulletVelocity rl.Vector2) error {
	game := GetGame()
	r.isAlive = false
	// So long as it isn't a tiny rock, spawn more smaller rocks at same loc
	if r.size > RockTiny {
		for range utils.RndInt32InRange(2, 4) {
			// Spawn a new rock at the same position as the old one but a bit back
			newRock := NewRock(r.size-1, rl.Vector2Add(r.Position, bulletVelocity))
			game.World.Objects.Add(&newRock)
			// Add a bit of bullet velocity to each new rock so more likely moving away
			newRock.Velocity = rl.Vector2Add(newRock.Velocity, rl.Vector2Scale(bulletVelocity, 0.1))
		}
	}
	// Spawn shrapnel in random directions and lifespans
	for range utils.RndInt32InRange(5, 10) {
		shrapnel := NewShrapnel(r.Position, uint16(utils.RndInt32InRange(300, 600)))
		game.World.Objects.Add(&shrapnel)
	}
	// Notify other services
	game.EventBus.Publish("rock:destroyed", r.size)

	return nil
}
