package playfield

import (
	"avoid_the_space_rocks/internal/core"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

// Constants for score stuff
const (
	shipExtraLife = 10_000
)

type ScoreKeeper struct {
	game *core.Game
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
	sk.game = game
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

func (sk *ScoreKeeper) Update(_ *core.Game) error {
	return nil
}

func (sk *ScoreKeeper) rockScoreHandler(size core.RockSize) {
	switch size {
	case core.RockTiny:
		sk.addPoints(100)
	case core.RockSmall:
		sk.addPoints(75)
	case core.RockMedium:
		sk.addPoints(50)
	case core.RockBig:
		sk.addPoints(25)
	}
}

func (sk *ScoreKeeper) alienScoreHandler(size core.AlienSize) {
	switch size {
	case core.AlienSmall:
		sk.addPoints(250)
	case core.AlienBig:
		sk.addPoints(500)
	}
}

func (sk *ScoreKeeper) addPoints(points int) {
	rewardLevel := uint64(math.Floor(float64(uint64(core.GetGame().Score/shipExtraLife))) + 1)
	pointsForNewLife := rewardLevel * shipExtraLife
	core.GetGame().Score += uint64(points)
	if core.GetGame().Score >= pointsForNewLife && sk.game.Lives < 20 {
		sk.game.Lives += 1
		sk.game.EventBus.Publish("spaceship:extra_life")
	}
}
