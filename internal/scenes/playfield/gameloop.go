package playfield

import (
	"avoid_the_space_rocks/internal/core"
	"avoid_the_space_rocks/internal/utils"
	"github.com/dustin/go-humanize"
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
	spaceship := &core.GetGame().World.Spaceship
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
		core.GetGame().Paused = !core.GetGame().Paused
	}
	spaceship.FuelBurning = rl.IsKeyDown(rl.KeyUp)
}

// Update all game state since last time through game loop
func update() {
	game := core.GetGame()
	if !game.Paused {
		if err := game.World.Spaceship.Update(); err != nil {
			rl.TraceLog(rl.LogError, "error updating spaceship: %v", err)
		}
		game.World.Objects.Update()
	}
}

// Draw all game state
func render() {
	game := core.GetGame()
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	drawHud()

	if err := game.World.Spaceship.Draw(); err != nil {
		rl.TraceLog(rl.LogError, "error drawing spaceship: %v", err)
	}
	game.World.Objects.Draw()

	if game.Paused {
		utils.CenterText("PAUSED", rl.Vector2{X: game.World.Width / 2, Y: game.World.Height / 3}, 40)
	}

	rl.EndDrawing()
}

// drawHud displays the score and the number of lives remaining
func drawHud() {
	game := core.GetGame()

	score := humanize.Comma(int64(game.Score))
	utils.WriteText(score, rl.Vector2{X: 10, Y: 10}, 30)

	size := game.World.Spaceship.SpriteSheet.GetSize()
	for i := range game.Lives {
		pos := rl.Vector2{X: game.World.Width - 20 - (float32(i) * size.X * 0.6), Y: 20 + (size.Y / 2)}
		if err := game.World.Spaceship.SpriteSheet.Draw(0, 0, pos, rl.Vector2{X: 0, Y: -1}); err != nil {
			rl.TraceLog(rl.LogError, "error drawing spaceship for lives: %v", err)
		}
	}
}
