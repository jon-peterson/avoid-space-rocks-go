package playfield

import (
	"avoid_the_space_rocks/internal/core"
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/scenes"
	"avoid_the_space_rocks/internal/utils"
	"github.com/dustin/go-humanize"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Init() {
	game := core.GetGame()
	game.Observers = append(game.Observers, NewAudioManager(), NewScoreKeeper(), NewGameWarden())

	for _, obs := range game.Observers {
		if err := obs.Register(game); err != nil {
			rl.TraceLog(rl.LogError, "error registering observer: %v", err)
		}
	}
}

func Close() {
	game := core.GetGame()
	for _, obs := range game.Observers {
		if err := obs.Deregister(game); err != nil {
			rl.TraceLog(rl.LogError, "error deregistering observer: %v", err)
		}
	}
}

func Loop() scenes.SceneCode {
	for !rl.WindowShouldClose() {
		handleInput()
		update()
		render()
	}
	return scenes.AttractMode
}

// Handle player input
func handleInput() {
	game := core.GetGame()
	if game.DebugMode {
		handleDebugInput()
	}
	// Player input
	spaceship := &game.World.Spaceship
	if spaceship.IsAlive() {
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
	// Game state input
	if rl.IsKeyPressed(rl.KeyEscape) {
		core.GetGame().Paused = !core.GetGame().Paused
	}
}

func handleDebugInput() {
	game := core.GetGame()
	if rl.IsKeyPressed(rl.KeyF1) {
		game.Lives += 1
	}
	if rl.IsKeyPressed(rl.KeyF2) {
		game.World.Objects.ForEach(func(obj gameobjects.GameObject) {
			if rock, ok := obj.(*core.Rock); ok {
				_ = rock.OnDestruction(rl.Vector2{})
			}
		})
	}
}

// Update all game state since last time through game loop
func update() {
	game := core.GetGame()
	if game.Paused {
		return
	}
	game.World.Objects.Update()
}

// Draw all game state
func render() {
	game := core.GetGame()
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	drawHud()

	game.World.Objects.Draw()

	if game.Paused {
		utils.CenterText("PAUSED", rl.Vector2{X: game.World.Width / 2, Y: game.World.Height / 3}, 40)
	} else if game.Overlay != nil {
		game.Overlay()
	}

	rl.EndDrawing()
}

// drawHud displays the score and the number of lives remaining
func drawHud() {
	game := core.GetGame()

	score := humanize.Comma(int64(game.Score))
	utils.WriteText(score, rl.Vector2{X: 15, Y: 12}, 36)

	size := game.World.Spaceship.Spritesheet.GetSize()
	for i := range game.Lives {
		pos := rl.Vector2{X: game.World.Width - 20 - (float32(i) * size.X * 0.6), Y: 20 + (size.Y / 2)}
		if err := game.World.Spaceship.Spritesheet.Draw(0, 0, pos, rl.Vector2{X: 0, Y: -1}); err != nil {
			rl.TraceLog(rl.LogError, "error drawing spaceship for lives: %v", err)
		}
	}
}
