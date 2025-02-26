package attractmode

import (
	"avoid_the_space_rocks/internal/scenes"
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
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
	for !rl.WindowShouldClose() {

		key := rl.GetKeyPressed()
		if key == rl.KeyEscape {
			return scenes.Quit
		} else if key != rl.KeyNull {
			return scenes.GameplayScene
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		utils.CenterText("Avoid", rl.Vector2{X: am.width / 2, Y: am.height/3 - 55}, 80)
		utils.CenterText("the", rl.Vector2{X: am.width / 2, Y: am.height / 3}, 40)
		utils.CenterText("Space Rocks", rl.Vector2{X: am.width / 2, Y: am.height/3 + 50}, 80)

		utils.CenterText("Press any key to start", rl.Vector2{X: am.width / 2, Y: am.height/3 + 180}, 20)
		rl.EndDrawing()
	}

	return scenes.Quit
}
