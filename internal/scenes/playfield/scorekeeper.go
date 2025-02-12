package playfield

import (
	"avoid_the_space_rocks/internal/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func RegisterScoreKeeper(game *core.Game) error {
	if err := game.EventBus.Subscribe("rock:destroyed", RockScoreHandler); err != nil {
		rl.TraceLog(rl.LogError, "error subscribing to rock:destroyed event: %v", err)
		return err
	}
	return nil
}

func DeregisterScoreKeeper(game *core.Game) error {
	if err := game.EventBus.Unsubscribe("rock:destroyed", RockScoreHandler); err != nil {
		rl.TraceLog(rl.LogError, "error unsubscribing from rock:destroyed event: %v", err)
		return err
	}
	return nil
}

func RockScoreHandler(size core.RockSize) {
	switch size {
	case core.RockTiny:
		core.GetGame().Score += 25
	case core.RockSmall:
		core.GetGame().Score += 50
	case core.RockMedium:
		core.GetGame().Score += 100
	case core.RockBig:
		core.GetGame().Score += 250
	}
}
