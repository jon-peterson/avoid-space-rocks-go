package playfield

import (
	"math"
	"sync"
)

var instance *Game
var once sync.Once

// Constants for gameplay feel
const (
	shipRotateSpeed float32 = math.Pi * 3 // 1.5 rotations per second
	shipMaxSpeed    float32 = 20.0        // 20 units per second
	shipDecaySpeed  float32 = 3.0         // units per second slower
	shipFuelBoost   float32 = 10.0        // units per second added to acceleration

	bulletSpeed      float32 = 10.0 // units per second
	bulletMaxSpeed   float32 = 20.0 // units per second
	bulletLifetimeMs uint16  = 1000

	shrapnelMaxSpeed  float32 = 6.0          // units per second
	shrapnelMaxRotate float32 = math.Pi * 12 // 6 rotations per second

	rockMaxSpeed  float32 = 10.0
	rockMaxRotate float32 = math.Pi * 6 // 3 rotations per second
)

type Game struct {
	World  World
	Lives  uint8
	Level  uint8
	Score  uint64
	Paused bool
}

func GetGame() *Game {
	if instance == nil {
		panic("Game not initialized. Call InitGame first")
	}
	return instance
}

func InitGame(screenWidth, screenHeight int32) *Game {
	once.Do(func() {
		w := NewWorld(screenWidth, screenHeight)
		instance = &Game{
			World: w,
			Lives: 3,
			Level: 1,
			Score: 0,
		}
	})
	return instance
}
