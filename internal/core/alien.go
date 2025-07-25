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
	row, col int
}{
	{"alien_small.png", 3, 3},
	{"alien_big.png", 2, 2},
}

// Alien spaceships
type Alien struct {
	gameobjects.Rigidbody
	spritesheet  *gameobjects.SpriteSheet
	isAlive      bool
	size         AlienSize
	bulletDrift  float32
	runnerCancel context.CancelFunc
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
func (a *Alien) Update(delta float32) error {
	game := GetGame()
	a.Rigidbody.ApplyPhysics(delta)
	if game.World.IsOutsideEdges(a.Position) {
		// If the alien goes outside the edges, we remove it from the game sometimes
		if utils.Chance(0.2) {
			a.isAlive = false
			game.EventBus.Publish("alien:left_playfield", a.size)
		} else {
			// Wrap around the edges of the playfield
			a.Position = game.World.Wraparound(a.Position)
		}
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
func (a *Alien) OnDestruction(_ rl.Vector2) error {
	game := GetGame()
	a.isAlive = false
	// Spawn shrapnel in random directions and lifespans
	sheet := gameobjects.LoadSpriteSheet("shrapnel.png", 5, 1)
	for range 6 {
		frame := int(utils.RndIntInRange(0, 4))
		shrapnel := NewShrapnel(a.Position, sheet, uint(utils.RndIntInRange(200, 400)), frame)
		game.World.Objects.Add(&shrapnel)
	}
	// Cancel the runner goroutine if it exists
	if a.runnerCancel != nil {
		a.runnerCancel()
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

// randomizeAlienTarget sets the alien's target to a random position on the playfield, at a random speed.
func (a *Alien) randomizeAlienTarget() {
	game := GetGame()
	target := game.World.RandomPosition()
	a.Velocity = rl.Vector2Normalize(rl.Vector2Subtract(target, a.Position))
	if a.size == AlienBig {
		sp := utils.RndFloat32InRange(alienMaxSpeed/3, alienMaxSpeed) / 2
		a.Velocity = rl.Vector2Scale(a.Velocity, sp)
	} else {
		sp := utils.RndFloat32InRange(alienMaxSpeed/3, alienMaxSpeed)
		a.Velocity = rl.Vector2Scale(a.Velocity, sp)
	}
}

// AlienSpawner adds new aliens to the playfield at an appropriate rate
func AlienSpawner(ctx context.Context) {
	rl.TraceLog(rl.LogDebug, "AlienSpawner starting")
	// Decide how frequently we should spawn aliens
	game := GetGame()
	var alien *Alien = nil
	var runnerCtx context.Context
	spawnDelay := time.Second * max(1, time.Duration(10.0-(float32(game.Level)*1.25)))

	ticker := time.NewTicker(spawnDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			rl.TraceLog(rl.LogDebug, "AlienSpawner exiting")
			if alien != nil && alien.runnerCancel != nil {
				alien.runnerCancel()
			}
			return

		case <-ticker.C:

			if game.Paused {
				continue
			}

			// Handle case where there's already an alien on the playfield
			if alien != nil {
				if alien.IsAlive() {
					// There's already an active alien in the level; let it run
					continue
				}
				rl.TraceLog(rl.LogInfo, "Alien no longer on playfield; stopping")
				alien.runnerCancel()
				alien = nil
				// Don't spawn another right away
				continue
			}

			// Try to spawn a new alien, but if the position is occupied just skip this time around
			position := game.World.RandomBorderPosition()
			alien = newSpawnedAlien(game, position)
			if game.World.Objects.IsRectangleOccupied(gameobjects.ExtendRectangle(alien.GetHitbox(), 0.25)) {
				continue
			}
			rl.TraceLog(rl.LogInfo, "Spawning new alien")
			runnerCtx, alien.runnerCancel = context.WithCancel(context.Background())
			go AlienRunner(runnerCtx, alien)
			game.World.Objects.Add(alien)
			game.EventBus.Publish("alien:spawned", alien.size)
		}
	}
}

func AlienRunner(ctx context.Context, alien *Alien) {
	rl.TraceLog(rl.LogDebug, "AlienRunner starting")
	game := GetGame()

	// Small aliens do things more frequently
	actionDelay := 3000 - int(300*game.Level)
	if alien.size == AlienSmall {
		actionDelay /= 2
	}
	if actionDelay < alienMinActionDelay {
		actionDelay = alienMinActionDelay
	}

	ticker := time.NewTicker(time.Millisecond * time.Duration(actionDelay))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			rl.TraceLog(rl.LogDebug, "AlienRunner exiting")
			return

		case <-ticker.C:
			if game.Paused {
				continue
			}
			if !alien.IsAlive() {
				// The alien associated with this runner is no longer alive; stop shooting permanently
				return
			}
			if !game.World.Spaceship.IsAlive() {
				continue
			}

			// Alien does one of three things: fire, change direction, nothing
			if utils.Chance(0.3) {
				// Change direction
				rl.TraceLog(rl.LogInfo, "Alien changing direction")
				alien.randomizeAlienTarget()
			} else if utils.Chance(0.5) {
				// Fire a bullet roughly towards the spaceship
				drift := utils.RndFloat32InRange(-alien.bulletDrift, alien.bulletDrift)
				shootDirection := rl.Vector2Normalize(rl.Vector2Subtract(game.World.Spaceship.Position, alien.Position))
				shootDirection = rl.Vector2Rotate(shootDirection, drift)
				bullet := NewBullet(alien.Position, rl.Vector2Scale(shootDirection, bulletSpeed), false)
				game.World.Objects.Add(&bullet)
				game.EventBus.Publish("alien:fire")
			}
		}
	}
}

// newSpawnedAlien returns a new alien at the specified position, moving in a random direction at the
// appropriate speed for that alien type. Smaller aliens are more common at higher levels.
func newSpawnedAlien(game *Game, position rl.Vector2) *Alien {
	// Spawn a new alien
	size := AlienBig
	if game.Level > 2 && utils.RndIntInRange(0, 10) < game.Level {
		size = AlienSmall
	}
	spawnedAlien := NewAlien(size, position)
	spawnedAlien.randomizeAlienTarget()

	if size == AlienBig {
		spawnedAlien.bulletDrift = utils.RndFloat32(alienMaxBulletDrift)
	} else {
		spawnedAlien.bulletDrift = utils.RndFloat32(alienMaxBulletDrift) / 3
	}
	return &spawnedAlien
}
