package playfield

import (
	"avoid_the_space_rocks/internal/core"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type GameWarden struct {
	game *core.Game
}

var _ core.EventObserver = (*GameWarden)(nil)

func (gw *GameWarden) eventMappings() []eventMapping {
	return []eventMapping{
		{"rock:spawned", gw.rockSpawnedWatcher},
		{"rock:destroyed", gw.rockDestroyedWatcher},
		{"alien:destroyed", gw.alienRemovedWatcher},
		{"alien:left_playfield", gw.alienRemovedWatcher},
		{"spaceship:destroyed", gw.spaceshipDestroyedWatcher},
		{"spaceship:enter_hyperspace", gw.spaceshipHyperspaceWatcher},
	}
}

func NewGameWarden() *GameWarden {
	return &GameWarden{}
}

func (gw *GameWarden) Register(game *core.Game) error {
	gw.game = game
	for _, sub := range gw.eventMappings() {
		if err := game.EventBus.SubscribeAsync(sub.event, sub.handler, false); err != nil {
			rl.TraceLog(rl.LogError, "error subscribing to %s event: %v", sub.event, err)
			return err
		}
	}
	return nil
}

func (gw *GameWarden) Deregister(game *core.Game) error {
	for _, sub := range gw.eventMappings() {
		if err := game.EventBus.Unsubscribe(sub.event, sub.handler); err != nil {
			rl.TraceLog(rl.LogError, "error unsubscribing from %s event: %v", sub.event, err)
			return err
		}
	}
	return nil
}

func (gw *GameWarden) Update(_ *core.Game) error {
	return nil
}

// rockSpawnedWatcher is called when a new rock is added to the level.
func (gw *GameWarden) rockSpawnedWatcher(_ core.RockSize) {
	gw.game.Rocks += 1
}

// rockDestroyedWatcher is called when a rock is destroyed. Calls the end-of-level check.
func (gw *GameWarden) rockDestroyedWatcher(_ core.RockSize) {
	gw.game.Rocks -= 1
	gw.checkEndOfLevel()
}

// alienDestroyedWatcher is called when an alien is destroyed or leaves the playfield. Calls the end-of-level check.
func (gw *GameWarden) alienRemovedWatcher(_ core.AlienSize) {
	gw.checkEndOfLevel()
}

// checkEndOfLevel sees if there are remaining enemies; if not, it starts the next level.
func (gw *GameWarden) checkEndOfLevel() {
	if !gw.game.World.Objects.HasRemainingEnemies() {
		gw.game.StopLevel()
		go gw.game.StartLevel()
	}
}

// SpaceshipDestroyedWatcher is called when the spaceship is destroyed. It decrements the lives
// remaining, waits a moment, and then respawns the spaceship. If the player is out of lives it goes to
// game over state.
func (gw *GameWarden) spaceshipDestroyedWatcher() {
	gw.game.Lives--
	time.Sleep(4 * time.Second)
	if gw.game.Lives > 0 {
		gw.game.World.Spaceship.Spawn()
	} else {
		rl.TraceLog(rl.LogInfo, "Game over")
		gw.game.Over = true
	}
}

// SpaceshipHyperspaceWatcher moves the spaceship to a random location with some graphic flair.
// The audio clip is two seconds but the re-entry clack is at 1.9 seconds, so adjust accordingly.
func (gw *GameWarden) spaceshipHyperspaceWatcher() {
	s := &gw.game.World.Spaceship
	// Stop the spaceship and put it in hyperspace
	s.InHyperspace = true
	s.Velocity = rl.Vector2{}
	s.Acceleration = rl.Vector2{}
	// Send four pieces of the spaceship to random locations
	pieces := make([]*core.Shrapnel, 4)
	for i := 0; i < 4; i++ {
		piece := core.NewShrapnel(s.Position, s.Spritesheet, 1800, i+3)
		pos := gw.game.World.RandomPosition()
		piece.Velocity = rl.Vector2Normalize(rl.Vector2Subtract(pos, piece.Position))
		piece.Velocity = rl.Vector2Scale(piece.Velocity, 300)
		pieces[i] = &piece
		gw.game.World.Objects.Add(&piece)
	}
	time.Sleep(900 * time.Millisecond)
	// Place the spaceship at a random location, and send the four pieces to its new location
	s.Position = gw.game.World.RandomPosition()
	for i := 0; i < 4; i++ {
		// Create a vector that points from piece.Position to s.Position
		piece := &pieces[i]
		(*piece).Velocity = rl.Vector2Normalize(rl.Vector2Subtract(s.Position, (*piece).Position))
		// Scale the velocity so it arrives at the new space position
		distance := rl.Vector2Distance(s.Position, (*piece).Position)
		(*piece).Velocity = rl.Vector2Scale((*piece).Velocity, distance)
	}
	time.Sleep(time.Second)
	s.InHyperspace = false
}
