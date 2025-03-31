package gameover

import (
	"avoid_the_space_rocks/internal/core"
	"avoid_the_space_rocks/internal/scenes"
	"avoid_the_space_rocks/internal/utils"
	"github.com/dustin/go-humanize"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type GameOverMode struct {
	width  float32
	height float32
}

var _ scenes.Scene = (*GameOverMode)(nil)

func (am *GameOverMode) Init(width, height float32) {
	am.width = width
	am.height = height
}

func (am *GameOverMode) Close() {
}

func (am *GameOverMode) Loop() scenes.SceneCode {
	next := scenes.GameOverScene

	go func() {
		time.Sleep(4 * time.Second)
		next = scenes.AttractModeScene
	}()

	for !rl.WindowShouldClose() && next == scenes.GameOverScene {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		utils.CenterText("Game Over", rl.Vector2{X: am.width / 2, Y: am.height / 3}, 80)

		score := humanize.Comma(int64(core.GetGame().Score))
		utils.CenterText("Your Score", rl.Vector2{X: am.width / 2, Y: am.height / 2}, 30)
		utils.CenterText(score, rl.Vector2{X: am.width / 2, Y: am.height/2 + 75}, 60)

		rl.EndDrawing()
	}

	if rl.WindowShouldClose() {
		return scenes.Quit
	}
	return next
}
