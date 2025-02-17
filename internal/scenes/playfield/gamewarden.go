package playfield

import (
	"avoid_the_space_rocks/internal/core"
	rl "github.com/gen2brain/raylib-go/raylib"
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
	gw.game = game
	return nil
}

func (gw *GameWarden) Deregister(game *core.Game) error {
	if err := game.EventBus.Unsubscribe("rock:destroyed", gw.EnemyDestroyedWatcher); err != nil {
		rl.TraceLog(rl.LogError, "error unsubscribing from rock:destroyed event: %v", err)
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
