package main

import (
	"avoid_the_space_rocks/internal/scenes"
	"avoid_the_space_rocks/internal/scenes/attractmode"
	"avoid_the_space_rocks/internal/scenes/playfield"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1024.0
	screenHeight = 768.0
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Avoid the Space Rocks")
	defer rl.CloseWindow()
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)

	sceneCode := scenes.AttractModeScene
	for sceneCode != scenes.Quit {
		rl.TraceLog(rl.LogInfo, "Starting scene code %v", sceneCode)
		scene := initScene(sceneCode)
		sceneCode = scene.Loop()
		scene.Close()
	}
}

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
		return nil
	} else {
		return nil
	}
}
