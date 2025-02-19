package core

import (
	"avoid_the_space_rocks/internal/utils"
	"fmt"
	evbus "github.com/asaskevich/EventBus"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"os"
	"sync"
	"time"
)

var instance *Game
var once sync.Once

// Constants for gameplay feel
const (
	shipRotateSpeed float32 = math.Pi * 3 // 1.5 rotations per second
	shipMaxSpeed    float32 = 7.50        // units per second max speed
	shipDecaySpeed  float32 = 1.0         // units per second slower
	shipFuelBoost   float32 = 3.0         // units per second added to acceleration

	bulletSpeed      float32 = 10.0  // units per second
	bulletMaxSpeed   float32 = 100.0 // units per second
	bulletLifetimeMs uint16  = 1000

	shrapnelMaxSpeed  float32 = 6.0          // units per second
	shrapnelMaxRotate float32 = math.Pi * 12 // 6 rotations per second

	rockMaxSpeed  float32 = 10.0
	rockMaxRotate float32 = math.Pi * 6 // 3 rotations per second
)

type Game struct {
	World *World

	Lives int32
	Level int32
	Score uint64

	Paused    bool
	DebugMode bool

	EventBus  evbus.Bus
	Observers []EventObserver

	Overlay func()
}

type EventObserver interface {
	Register(game *Game) error
	Deregister(game *Game) error
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
			World:     w,
			Lives:     3,
			Level:     0,
			Score:     0,
			Paused:    false,
			EventBus:  evbus.New(),
			Observers: make([]EventObserver, 0, 10),
		}
		if os.Getenv("DEBUG") != "" {
			instance.DebugMode = true
		}
	})
	return instance
}

// StartLevel initializes the level.
func (g *Game) StartLevel() {
	g.Level += 1
	rl.TraceLog(rl.LogInfo, fmt.Sprintf("Starting level %d", g.Level))
	g.Overlay = func() {
		utils.CenterText(fmt.Sprintf("Level %d", g.Level), rl.Vector2{X: g.World.Width / 2, Y: g.World.Height / 3}, 60)
	}
	time.Sleep(time.Second * 2)
	g.Overlay = nil
	time.Sleep(time.Millisecond * 500)
	for range g.Level + 3 {
		rock := NewRock(RockBig, g.World.RandomBorderLocation())
		g.World.Objects.Add(&rock)
	}
}
