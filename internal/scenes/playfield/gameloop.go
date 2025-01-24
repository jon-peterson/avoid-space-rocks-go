package playfield

import rl "github.com/gen2brain/raylib-go/raylib"

func GameLoop(world World) {

	for !rl.WindowShouldClose() {
		handleInput(world)
		update(world)
		render(world)
	}
}

// Handle player input
func handleInput(world World) {

}

// Update all game state since last time through game loop
func update(world World) {

}

// Draw all game state
func render(world World) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	rl.DrawTexture(world.Spaceship.SpriteSheet,
		int32(world.Spaceship.Transform.Position.X), int32(world.Spaceship.Transform.Position.Y),
		rl.Black)

	rl.EndDrawing()
}
