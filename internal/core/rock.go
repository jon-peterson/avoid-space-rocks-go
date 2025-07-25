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
	spritesheet   *gameobjects.SpriteSheet
	rotationSpeed float32 // rotations per second
	isAlive       bool
	size          RockSize
}

var _ gameobjects.Collidable = (*Rock)(nil)
var _ gameobjects.Destructible = (*Rock)(nil)
var _ gameobjects.GameObject = (*Rock)(nil)

func NewRock(size RockSize, position rl.Vector2) Rock {
	sheet := gameobjects.LoadSpriteSheet(rockSpriteFile[size], 1, 1)
	rock := Rock{
		spritesheet: sheet,
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
	maxSpeed := rockMaxSpeed / float32(size+2)
	rock.Velocity = rl.Vector2{
		X: utils.RndFloat32InRange(-maxSpeed, maxSpeed),
		Y: utils.RndFloat32InRange(-maxSpeed, maxSpeed),
	}
	return rock
}

// Update applies physics to the rock so it moves per its velocity and rotation speed.
func (r *Rock) Update(delta float32) error {
	game := GetGame()
	r.Rotation = rl.Vector2Rotate(r.Rotation, r.rotationSpeed*delta)
	r.Rigidbody.ApplyPhysics(delta)
	r.Position = game.World.Wraparound(r.Position)
	return nil
}

// Draw renders the rock to the screen.
func (r *Rock) Draw() error {
	return r.spritesheet.Draw(0, 0, r.Position, r.Rotation)
}

// IsAlive returns whether the rock is alive or not.
func (r *Rock) IsAlive() bool {
	return r.isAlive
}

func (r *Rock) IsEnemy() bool {
	return true
}

// GetHitbox returns the hitbox of the rock, used for basic collision detection.
func (r *Rock) GetHitbox() rl.Rectangle {
	return r.spritesheet.GetRectangle(r.Position)
}

// OnCollision handles the collision of the rock with another Collidable object. If object
// is destructible it destroys it -- unless it's a rock; rocks don't destroy rocks in this game.
func (r *Rock) OnCollision(other gameobjects.Collidable) error {
	if destructible, ok := other.(gameobjects.Destructible); ok {
		if _, ok := other.(*Rock); ok {
			return nil
		}
		return destructible.OnDestruction(r.Velocity)
	}
	return nil
}

// OnDestruction handles the destruction of the rock, spawning smaller rocks if applicable.
// This is called by the bullet's OnCollision method when it hits this rock.
func (r *Rock) OnDestruction(bulletVelocity rl.Vector2) error {
	game := GetGame()
	r.isAlive = false
	// Spawn smaller rocks at same location as appropriate for level
	if int(r.size) > max(0, 4-game.Level) {
		// Span more rocks at higher levels, but if we've hit our cap, replace one for one
		toSpawn := utils.RndIntInRange(2, max(3, int(game.Level/2)))
		if game.Rocks >= rockMaxCount {
			toSpawn = 1
		}
		for range toSpawn {
			// Spawn a new rock at the same position as the old one but a bit away from dir of the bullet
			newRock := NewRock(r.size-1, r.Position)
			spriteWidth := newRock.spritesheet.GetRectangle(newRock.Position).Width / 2
			scaledBulletVelocity := rl.Vector2Scale(rl.Vector2Normalize(bulletVelocity), spriteWidth)
			newRock.Position = rl.Vector2Add(newRock.Position, scaledBulletVelocity)
			game.World.Objects.Add(&newRock)
			game.EventBus.Publish("rock:spawned", r.size)
			// Add a bit of bullet velocity to each new rock so more likely moving away
			newRock.Velocity = rl.Vector2Add(newRock.Velocity, rl.Vector2Scale(bulletVelocity, 0.1))
		}
	}
	// Spawn shrapnel in random directions and lifespans
	sheet := gameobjects.LoadSpriteSheet("shrapnel.png", 5, 1)
	for range utils.RndIntInRange(int(r.size)+2, int(r.size*2)+4) {
		frame := int(utils.RndIntInRange(0, 4))
		shrapnel := NewShrapnel(r.Position, sheet, uint(utils.RndIntInRange(300, 600)), frame)
		game.World.Objects.Add(&shrapnel)
	}
	// Notify other services
	game.EventBus.Publish("rock:destroyed", r.size)

	return nil
}
