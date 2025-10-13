package main

import (
	"avoid_the_space_rocks/internal/scenes"
	"avoid_the_space_rocks/internal/scenes/attractmode"
	"avoid_the_space_rocks/internal/scenes/gameover"
	"avoid_the_space_rocks/internal/scenes/playfield"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

// main is the entry point for the application. Spawns a window and runs the game inside.
func main() {
	rl.InitWindow(screenWidth, screenHeight, "Avoid the Space Rocks")
	defer rl.CloseWindow()
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)

	if os.Getenv("DEBUG") != "" {
		rl.SetTraceLogLevel(rl.LogDebug)
	}

	sceneCode := scenes.AttractModeScene
	for sceneCode != scenes.Quit {
		rl.TraceLog(rl.LogInfo, "Starting scene code %v", sceneCode)
		scene := initScene(sceneCode)
		sceneCode = scene.Loop()
		scene.Close()
	}
}

// initScene initializes and returns the scene corresponding to the given scene code.
func initScene(code scenes.SceneCode) scenes.Scene {
	if code == scenes.AttractModeScene {
		am := &attractmode.AttractMode{}
		am.Init(screenWidth, screenHeight)
		return am
	} else if code == scenes.GameplayScene {
		gm := &playfield.Gameloop{}
		gm.Init(screenWidth, screenHeight)
		return gm
	} else if code == scenes.GameOverScene {
		gom := &gameover.GameOverMode{}
		gom.Init(screenWidth, screenHeight)
		return gom
	} else {
		rl.TraceLog(rl.LogError, "Unknown scene code %v", code)
		return nil
	}
}
