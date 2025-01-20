package playfield

import rl "github.com/gen2brain/raylib-go/raylib"

func GameLoop() {
	rl.InitWindow(screenWidth, screenHeight, "Avoid the Space Rocks")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Avoid the Space Rocks", 190, 200, 20, rl.Black)

		rl.EndDrawing()
	}
}
