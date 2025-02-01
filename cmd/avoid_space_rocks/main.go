package main

import (
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

	// For now there's only one screen
	game := playfield.InitGame(screenWidth, screenHeight)
	game.World.InitializeLevel(1)
	playfield.GameLoop()
}
