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
	game := GetGame()
	delta := rl.GetFrameTime()
	if rl.IsKeyDown(rl.KeyLeft) {
		game.World.Spaceship.Rotation = rl.Vector2Rotate(game.World.Spaceship.Rotation, -shipRotateSpeed*delta)
	}
	if rl.IsKeyDown(rl.KeyRight) {
		game.World.Spaceship.Rotation = rl.Vector2Rotate(game.World.Spaceship.Rotation, shipRotateSpeed*delta)
	}
	game.World.Spaceship.FuelBurning = rl.IsKeyDown(rl.KeyUp)
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
