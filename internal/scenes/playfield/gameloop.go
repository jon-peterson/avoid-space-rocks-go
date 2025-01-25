package playfield

import rl "github.com/gen2brain/raylib-go/raylib"

func GameLoop(world *World) {

	for !rl.WindowShouldClose() {
		handleInput(world)
		update(world)
		render(world)
	}
}

// Handle player input
func handleInput(world *World) {
	delta := rl.GetFrameTime()
	if rl.IsKeyDown(rl.KeyLeft) {
		world.Spaceship.Rotation = rl.Vector2Rotate(world.Spaceship.Rotation, -spaceshipRotateSpeed*delta)
	}
	if rl.IsKeyDown(rl.KeyRight) {
		world.Spaceship.Rotation = rl.Vector2Rotate(world.Spaceship.Rotation, spaceshipRotateSpeed*delta)
	}
}

// Update all game state since last time through game loop
func update(world *World) {

}

// Draw all game state
func render(world *World) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	world.Spaceship.Draw()

	rl.EndDrawing()
}
