package playfield

import (
	"avoid_the_space_rocks/internal/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ScoreKeeper struct{}

func NewScoreKeeper() *ScoreKeeper {
	return &ScoreKeeper{}
}

func (sk *ScoreKeeper) Register(game *core.Game) error {
	if err := game.EventBus.Subscribe("rock:destroyed", sk.RockScoreHandler); err != nil {
		rl.TraceLog(rl.LogError, "error subscribing to rock:destroyed event: %v", err)
		return err
	}
	return nil
}

func (sk *ScoreKeeper) Deregister(game *core.Game) error {
	if err := game.EventBus.Unsubscribe("rock:destroyed", sk.RockScoreHandler); err != nil {
		rl.TraceLog(rl.LogError, "error unsubscribing from rock:destroyed event: %v", err)
		return err
	}
	return nil
}

func (sk *ScoreKeeper) RockScoreHandler(size core.RockSize) {
	switch size {
	case core.RockTiny:
		core.GetGame().Score += 100
	case core.RockSmall:
		core.GetGame().Score += 75
	case core.RockMedium:
		core.GetGame().Score += 50
	case core.RockBig:
		core.GetGame().Score += 20
	}
}
