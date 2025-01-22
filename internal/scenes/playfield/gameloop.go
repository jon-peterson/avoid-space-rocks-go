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

}

// Update all game state since last time through game loop
func update() {

}

// Draw all game state
func render() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	rl.EndDrawing()
}
