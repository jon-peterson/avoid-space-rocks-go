package playfield

import (
	"avoid_the_space_rocks/internal/core"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type GameWarden struct {
	game *core.Game
}

func (gw *GameWarden) eventMappings() []eventMapping {
	return []eventMapping{
		{"rock:destroyed", gw.EnemyDestroyedWatcher},
		{"spaceship:destroyed", gw.SpaceshipDestroyedWatcher},
		{"spaceship:enter_hyperspace", gw.SpaceshipHyperspaceWatcher},
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

// EnemyDestroyedWatcher is called when an enemy is destroyed. If the level no longer has any live
// enemies, then it runs the next level.
func (gw *GameWarden) EnemyDestroyedWatcher(_ core.RockSize) {
	// If all rocks are done we can launch the next level
	if !gw.game.World.Objects.HasRemainingEnemies() {
		go gw.game.StartLevel()
	}
}

// SpaceshipDestroyedWatcher is called when the spaceship is destroyed. It decrements the lives
// remaining, waits a moment, and then respawns the spaceship. If the player is out of lives it goes to
// game over state.
func (gw *GameWarden) SpaceshipDestroyedWatcher() {
	gw.game.Lives--
	time.Sleep(4 * time.Second)
	if gw.game.Lives > 0 {
		gw.game.World.Spaceship.Spawn()
	} else {
		rl.TraceLog(rl.LogInfo, "Game over")
		gw.game.Over = true
	}
}

// SpaceshipHyperspaceWatcher moves the spaceship to a random location with some graphic flair
func (gw *GameWarden) SpaceshipHyperspaceWatcher() {
	s := &gw.game.World.Spaceship
	s.InHyperspace = true
	time.Sleep(4 * time.Second)
	s.Position = gw.game.World.RandomPosition()
	s.InHyperspace = false
}
