package playfield

import (
	"avoid_the_space_rocks/internal/core"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type GameWarden struct {
	game *core.Game
}

func NewGameWarden() *GameWarden {
	return &GameWarden{}
}

func (gw *GameWarden) Register(game *core.Game) error {
	if err := game.EventBus.Subscribe("rock:destroyed", gw.EnemyDestroyedWatcher); err != nil {
		rl.TraceLog(rl.LogError, "error subscribing to rock:destroyed event: %v", err)
		return err
	}
	if err := game.EventBus.SubscribeAsync("spaceship:destroyed", gw.SpaceshipDestroyedWatcher, true); err != nil {
		rl.TraceLog(rl.LogError, "error subscribing to spaceship:destroyed event: %v", err)
		return err
	}
	gw.game = game
	return nil
}

func (gw *GameWarden) Deregister(game *core.Game) error {
	if err := game.EventBus.Unsubscribe("rock:destroyed", gw.EnemyDestroyedWatcher); err != nil {
		rl.TraceLog(rl.LogError, "error unsubscribing from rock:destroyed event: %v", err)
		return err
	}
	if err := game.EventBus.Unsubscribe("spaceship:destroyed", gw.SpaceshipDestroyedWatcher); err != nil {
		rl.TraceLog(rl.LogError, "error unsubscribing from spaceship:destroyed event: %v", err)
		return err
	}
	gw.game = game
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
	time.Sleep(5 * time.Second)
	if gw.game.Lives > 0 {
		gw.game.World.Spaceship.Spawn()
	} else {
		gw.game.GameOver()
	}
}
