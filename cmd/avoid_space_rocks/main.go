package main

import (
	"avoid_the_space_rocks/internal/core"
	"avoid_the_space_rocks/internal/scenes/playfield"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Avoid the Space Rocks")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)

	// For now there's only one screen, jump right into it
	game := core.InitGame(screenWidth, screenHeight)
	game.World.InitializeLevel(1)

	playfield.RegisterScoreKeeper(game)
	playfield.GameLoop()
	playfield.DeregisterScoreKeeper(game)
}
