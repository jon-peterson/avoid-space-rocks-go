package core

import (
	"avoid_the_space_rocks/internal/utils"
	"context"
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
	shipMaxSpeed    float32 = 15.0
	shipDecaySpeed  float32 = 2.0
	shipFuelBoost   float32 = 40.0

	bulletSpeed      float32 = 500.0
	bulletLifetimeMs uint16  = 1250

	shrapnelMaxSpeed  float32 = 500.0
	shrapnelMaxRotate float32 = math.Pi * 12 // 6 rotations per second

	rockMaxSpeed  float32 = 200.0
	rockMaxRotate float32 = math.Pi * 6 // 3 rotations per second

	alienMaxSpeed       float32 = 300.0
	alienMaxBulletDrift float32 = math.Pi / 2
	alienMinShootDelay          = 500
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

	levelOver context.CancelFunc
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

// StartLevel kicks off a new level. Run this as a goroutine so that the physics
// engine and everything keeps running.
func (g *Game) StartLevel() {
	g.Level += 1

	// Display the level number for a few seconds
	rl.TraceLog(rl.LogInfo, fmt.Sprintf("Starting level %d", g.Level))
	g.Overlay = func() {
		utils.CenterText(fmt.Sprintf("Level %d", g.Level), rl.Vector2{X: g.World.Width / 2, Y: g.World.Height / 3}, 60)
	}
	time.Sleep(time.Second * 2)
	g.Overlay = nil
	time.Sleep(time.Millisecond * 500)

	// Kick off the alien spawner
	ctx, cancel := context.WithCancel(context.Background())
	g.levelOver = cancel
	go AlienSpawner(ctx)

	// Spawn the appropriate number of rocks
	for range g.Level + 3 {
		rock := NewRock(RockBig, g.World.RandomBorderPosition())
		g.World.Objects.Add(&rock)
	}
}

// StopLevel runs the end of level logic
func (g *Game) StopLevel() {
	g.levelOver()
}

// GameOver is called when the player has no more lives.
func (g *Game) GameOver() {
	g.Overlay = func() {
		utils.CenterText("Game Over", rl.Vector2{X: g.World.Width / 2, Y: g.World.Height / 3}, 60)
	}
	time.Sleep(time.Second * 5)
	g.Overlay = nil
}
