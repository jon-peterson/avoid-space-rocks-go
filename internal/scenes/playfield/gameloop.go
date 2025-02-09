package playfield

import (
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func GameLoop() {

	for !rl.WindowShouldClose() {
		handleInput()
		update()
		render()
	}
}

// Handle player input
func handleInput() {
	spaceship := &GetGame().World.Spaceship
	if rl.IsKeyDown(rl.KeyLeft) {
		spaceship.RotateLeft()
	}
	if rl.IsKeyDown(rl.KeyRight) {
		spaceship.RotateRight()
	}
	if rl.IsKeyPressed(rl.KeySpace) {
		spaceship.Fire()
	}
	if rl.IsKeyPressed(rl.KeyEscape) {
		GetGame().Paused = !GetGame().Paused
	}
	spaceship.FuelBurning = rl.IsKeyDown(rl.KeyUp)
}

// Update all game state since last time through game loop
func update() {
	game := GetGame()
	if !game.Paused {
		if err := game.World.Spaceship.Update(); err != nil {
			rl.TraceLog(rl.LogError, "error updating spaceship: %v", err)
		}
		game.World.Objects.Update()
	}
}

// Draw all game state
func render() {
	game := GetGame()
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	if err := game.World.Spaceship.Draw(); err != nil {
		rl.TraceLog(rl.LogError, "error drawing spaceship: %v", err)
	}
	game.World.Objects.Draw()

	if game.Paused {
		utils.CenterText("PAUSED", rl.Vector2{X: game.World.width / 2, Y: game.World.height / 2}, 40)
	}

	rl.EndDrawing()
}
