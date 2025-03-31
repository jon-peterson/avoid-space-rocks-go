package playfield

import (
	"avoid_the_space_rocks/internal/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ScoreKeeper struct {
}

var _ core.EventObserver = (*ScoreKeeper)(nil)

func (sk *ScoreKeeper) eventMappings() []eventMapping {
	return []eventMapping{
		{"rock:destroyed", sk.rockScoreHandler},
		{"alien:destroyed", sk.alienScoreHandler},
	}
}
func NewScoreKeeper() *ScoreKeeper {
	return &ScoreKeeper{}
}

func (sk *ScoreKeeper) Register(game *core.Game) error {
	for _, sub := range sk.eventMappings() {
		if err := game.EventBus.SubscribeAsync(sub.event, sub.handler, false); err != nil {
			rl.TraceLog(rl.LogError, "error subscribing to %s event: %v", sub.event, err)
			return err
		}
	}
	return nil
}

func (sk *ScoreKeeper) Deregister(game *core.Game) error {
	for _, sub := range sk.eventMappings() {
		if err := game.EventBus.Unsubscribe(sub.event, sub.handler); err != nil {
			rl.TraceLog(rl.LogError, "error unsubscribing from %s event: %v", sub.event, err)
			return err
		}
	}
	return nil
}

func (sk *ScoreKeeper) Update(game *core.Game) error {
	return nil
}

func (sk *ScoreKeeper) rockScoreHandler(size core.RockSize) {
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

func (sk *ScoreKeeper) alienScoreHandler(size core.AlienSize) {
	switch size {
	case core.AlienSmall:
		core.GetGame().Score += 500
	case core.AlienBig:
		core.GetGame().Score += 1000
	}
}
