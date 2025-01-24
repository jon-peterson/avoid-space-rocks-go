package main

import (
	"avoid_the_space_rocks/internal/scenes/playfield"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Avoid the Space Rocks")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// For now there's only one screen
	world := playfield.MakeWorld(screenWidth, screenHeight)
	playfield.GameLoop(world)
}
