package playfield

import rl "github.com/gen2brain/raylib-go/raylib"

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
	spaceship.FuelBurning = rl.IsKeyDown(rl.KeyUp)
}

// Update all game state since last time through game loop
func update() {
	game := GetGame()
	if err := game.World.Spaceship.Update(); err != nil {
		rl.TraceLog(rl.LogError, "error updating spaceship: %v", err)
	}
	game.World.Objects.Update()
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
	rl.EndDrawing()
}
