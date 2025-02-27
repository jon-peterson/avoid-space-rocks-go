package core

import (
	"avoid_the_space_rocks/internal/utils"
	"fmt"
	evbus "github.com/asaskevich/EventBus"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"os"
	"time"
)

var instance *Game

// Constants for gameplay feel
const (
	shipRotateSpeed float32 = math.Pi * 3 // 1.5 rotations per second
	shipMaxSpeed    float32 = 10.0        // 20 units per second
	shipDecaySpeed  float32 = 1.5         // units per second slower
	shipFuelBoost   float32 = 10.0        // units per second added to acceleration

	bulletSpeed      float32 = 20.0 // units per second
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
	Over      bool

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

func InitGame(screenWidth, screenHeight float32) *Game {
	w := NewWorld(screenWidth, screenHeight)
	instance = &Game{
		World:     w,
		Lives:     3,
		EventBus:  evbus.New(),
		Observers: make([]EventObserver, 0, 10),
	}
	if os.Getenv("DEBUG") != "" {
		instance.DebugMode = true
	}
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

// GameOver is called when the player has no more lives.
func (g *Game) GameOver() {
	g.Overlay = func() {
		utils.CenterText("Game Over", rl.Vector2{X: g.World.Width / 2, Y: g.World.Height / 3}, 60)
	}
	time.Sleep(time.Second * 5)
	g.Overlay = nil
}
