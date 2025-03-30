package attractmode

import (
	"avoid_the_space_rocks/internal/scenes"
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type AttractMode struct {
	width  float32
	height float32
}

var _ scenes.Scene = (*AttractMode)(nil)

func (am *AttractMode) Init(width, height float32) {
	am.width = width
	am.height = height
}

func (am *AttractMode) Close() {
}

func (am *AttractMode) Loop() scenes.SceneCode {
	screenDuration := time.Second * time.Duration(5)
	lastSwitchTime := rl.GetTime()
	currentScreen := 0

	for !rl.WindowShouldClose() {

		key := rl.GetKeyPressed()
		if key == rl.KeyEscape {
			return scenes.Quit
		} else if key != rl.KeyNull {
			return scenes.GameplayScene
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if rl.GetTime()-lastSwitchTime >= screenDuration.Seconds() {
			currentScreen = (currentScreen + 1) % 3
			lastSwitchTime = rl.GetTime()
		}

		switch currentScreen {
		case 0:
			am.titleScreen()
		case 1:
			am.howToPlayScreen()
		case 2:
			am.keymappingsScreen()
		}

		utils.CenterText("Press any key to start", rl.Vector2{X: am.width / 2, Y: am.height - 100}, 20)
		rl.EndDrawing()
	}

	return scenes.Quit
}

func (am *AttractMode) titleScreen() {
	utils.CenterText("Avoid", rl.Vector2{X: am.width / 2, Y: am.height/3 - 55}, 80)
	utils.CenterText("the", rl.Vector2{X: am.width / 2, Y: am.height / 3}, 40)
	utils.CenterText("Space Rocks", rl.Vector2{X: am.width / 2, Y: am.height/3 + 50}, 80)
}

func (am *AttractMode) howToPlayScreen() {
	utils.CenterText("How to Play", rl.Vector2{X: am.width / 2, Y: am.height/3 - 125}, 70)
	utils.CenterText("1. Avoid rocks", rl.Vector2{X: am.width / 2, Y: am.height / 3}, 50)
	utils.CenterText("2. Shoot aliens", rl.Vector2{X: am.width / 2, Y: am.height/3 + 50}, 50)
}

func (am *AttractMode) keymappingsScreen() {
	utils.CenterText("Keys", rl.Vector2{X: am.width / 2, Y: am.height/3 - 125}, 70)

	utils.CenterText("left", rl.Vector2{X: am.width/2 - 175, Y: am.height / 3}, 50)
	utils.CenterText("right", rl.Vector2{X: am.width/2 - 175, Y: am.height/3 + 50}, 50)
	utils.CenterText("up", rl.Vector2{X: am.width/2 - 175, Y: am.height/3 + 100}, 50)
	utils.CenterText("space", rl.Vector2{X: am.width/2 - 175, Y: am.height/3 + 150}, 50)
	utils.CenterText("enter", rl.Vector2{X: am.width/2 - 175, Y: am.height/3 + 200}, 50)

	utils.CenterText("Rotate left", rl.Vector2{X: am.width/2 + 175, Y: am.height / 3}, 50)
	utils.CenterText("Rotate right", rl.Vector2{X: am.width/2 + 175, Y: am.height/3 + 50}, 50)
	utils.CenterText("Thrust", rl.Vector2{X: am.width/2 + 175, Y: am.height/3 + 100}, 50)
	utils.CenterText("Fire", rl.Vector2{X: am.width/2 + 175, Y: am.height/3 + 150}, 50)
	utils.CenterText("Hyperspace", rl.Vector2{X: am.width/2 + 175, Y: am.height/3 + 200}, 50)
}
