package core

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/utils"
	"context"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type AlienSize int

const (
	AlienSmall AlienSize = iota
	AlienBig
)

// Create a constant array of four string elements
var alienSpriteFile = []string{
	"alien_small.png",
	"alien_big.png",
}

// Alien spaceships
type Alien struct {
	gameobjects.Rigidbody
	spritesheet *gameobjects.SpriteSheet
	isAlive     bool
	size        AlienSize
}

var _ gameobjects.Collidable = (*Alien)(nil)
var _ gameobjects.GameObject = (*Alien)(nil)

func NewAlien(size AlienSize, position rl.Vector2) Alien {
	sheet := gameobjects.LoadSpriteSheet(alienSpriteFile[size], 1, 1)
	alien := Alien{
		spritesheet: sheet,
		Rigidbody: gameobjects.Rigidbody{
			Transform: gameobjects.Transform{
				Position: position,
				Rotation: rl.Vector2{X: 1, Y: 0},
			},
		},
		isAlive: true,
		size:    size,
	}
	return alien
}

// Update applies physics to the alien so it moves along its current direction
func (a *Alien) Update() error {
	game := GetGame()
	a.Rigidbody.ApplyPhysics()
	if game.World.IsOutsideEdges(a.Position) {
		// If the alien goes outside the edges, we just remove it from the game
		a.isAlive = false
	}
	return nil
}

// Draw renders the alien  to the screen
func (a *Alien) Draw() error {
	// TODO: add some sort of animation
	return a.spritesheet.Draw(0, 0, a.Position, a.Rotation)
}

// IsAlive returns whether the alien is alive or not
func (a *Alien) IsAlive() bool {
	return a.isAlive
}

// IsEnemy returns true; always true for Aliens
func (a *Alien) IsEnemy() bool {
	return true
}

// GetHitbox returns the hitbox of the alien, used for basic collision detection.
func (a *Alien) GetHitbox() rl.Rectangle {
	return a.spritesheet.GetRectangle(a.Position)
}

// OnCollision handles the collision with another Collidable object.
// TODO: Handle collision with rocks, which should destroy the alien
func (a *Alien) OnCollision(other gameobjects.Collidable) error {
	s, ok := other.(*Spaceship)
	if ok {
		return s.OnDestruction(a.Velocity)
	}
	return nil
}

// OnDestruction handles the destruction of the alien.
func (a *Alien) OnDestruction(bulletVelocity rl.Vector2) error {
	game := GetGame()
	a.isAlive = false
	// Spawn shrapnel in random directions and lifespans
	sheet := gameobjects.LoadSpriteSheet("shrapnel.png", 5, 1)
	for range 6 {
		frame := int(utils.RndInt32InRange(0, 4))
		shrapnel := NewShrapnel(a.Position, sheet, uint16(utils.RndInt32InRange(200, 400)), frame)
		game.World.Objects.Add(&shrapnel)
	}
	// Notify other services
	game.EventBus.Publish("alien:destroyed", a.size)
	return nil
}

// AlienSpawner adds new aliens to the playfield at an appropriate rate
func AlienSpawner(ctx context.Context) {
	rl.TraceLog(rl.LogInfo, "AlienSpawner starting")
	tickerStep := time.Millisecond * 250
	ticker := time.NewTicker(tickerStep)
	defer ticker.Stop()

	// Decide how frequently we should spawn aliens
	game := GetGame()
	var alien *Alien = nil
	var sinceLastSpawn time.Duration = 0.0
	spawnDelay := time.Second * max(1, time.Duration(3-game.Level))

	for {
		select {
		case <-ctx.Done():
			rl.TraceLog(rl.LogInfo, "AlienSpawner exiting")
			return

		case <-ticker.C:

			if game.Paused {
				continue
			}

			if alien != nil {
				if alien.IsAlive() {
					// There's already an active alien in the level; let it run
					continue
				}
				rl.TraceLog(rl.LogInfo, "Alien no longer on playfield; eligible to spawn a new one")
				alien = nil
				sinceLastSpawn = 0.0
			}

			sinceLastSpawn += tickerStep
			if sinceLastSpawn < spawnDelay {
				// Wait a bit longer
				continue
			}

			// Spawn a new alien
			rl.TraceLog(rl.LogInfo, "Spawning new alien")
			position := game.World.RandomBorderPosition()
			spawnedAlien := NewAlien(AlienBig, position)

			// Point the alien towards the target
			target := game.World.RandomPosition()
			spawnedAlien.MaxVelocity = alienBigMaxSpeed
			spawnedAlien.Velocity = rl.Vector2Normalize(rl.Vector2Subtract(target, spawnedAlien.Position))
			spawnedAlien.Velocity = rl.Vector2Scale(spawnedAlien.Velocity, alienBigMaxSpeed)

			alien = &spawnedAlien
			game.World.Objects.Add(alien)
			sinceLastSpawn = 0.0
		}
	}
}
