package core

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/utils"
	"context"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"time"
)

type AlienSize int

const (
	AlienSmall AlienSize = iota
	AlienBig
)

// The sprite file specifies filename, rows, and columns
var alienSpriteFile = []struct {
	filename string
	row, col int32
}{
	{"alien_small.png", 3, 3},
	{"alien_big.png", 2, 2},
}

// Alien spaceships
type Alien struct {
	gameobjects.Rigidbody
	spritesheet *gameobjects.SpriteSheet
	isAlive     bool
	size        AlienSize
}

var _ gameobjects.Collidable = (*Alien)(nil)
var _ gameobjects.Destructible = (*Alien)(nil)
var _ gameobjects.GameObject = (*Alien)(nil)

func NewAlien(size AlienSize, position rl.Vector2) Alien {
	spriteFile := alienSpriteFile[size]
	sheet := gameobjects.LoadSpriteSheet(spriteFile.filename, spriteFile.row, spriteFile.col)
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
	row, col, err := a.spritesheet.FrameLocation(a.frameIndex())
	if err != nil {
		rl.TraceLog(rl.LogError, "Error getting alien frame location: %v", err)
		row = 0
		col = 0
	}
	return a.spritesheet.Draw(row, col, a.Position, a.Rotation)
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

// OnCollision handles the collision with another Collidable object. Aliens can blow up
// spaceships only; they are in turn destroyed by rocks.
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

// frameIndex returns the index of the correct frame to use given the current time
func (a *Alien) frameIndex() int {
	// AlienBig is 2x2; AlienSmall is 3x3 but only 7 frames
	frameCount := 4
	if a.size == AlienSmall {
		frameCount = 7
	}
	halfSeconds := int(math.Floor(rl.GetTime() * 2))
	return halfSeconds % frameCount
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
			size := AlienBig
			if game.Level > 2 && utils.RndInt32InRange(0, 10) < game.Level {
				size = AlienSmall
			}

			position := game.World.RandomBorderPosition()
			spawnedAlien := NewAlien(size, position)

			// Point the alien towards the target
			target := game.World.RandomPosition()
			spawnedAlien.Velocity = rl.Vector2Normalize(rl.Vector2Subtract(target, spawnedAlien.Position))

			if size == AlienSmall {
				sp := utils.RndFloat32InRange(alienSmallMaxSpeed/2, alienSmallMaxSpeed)
				spawnedAlien.Velocity = rl.Vector2Scale(spawnedAlien.Velocity, sp)
			} else {
				sp := utils.RndFloat32InRange(alienBigMaxSpeed/2, alienBigMaxSpeed)
				spawnedAlien.Velocity = rl.Vector2Scale(spawnedAlien.Velocity, sp)
			}

			alien = &spawnedAlien
			game.World.Objects.Add(alien)
			sinceLastSpawn = 0.0
		}
	}
}
