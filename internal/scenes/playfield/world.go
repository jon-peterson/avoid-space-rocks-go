package playfield

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

// The World object represents the state of the game within the playfield
type World struct {
	width     float32 // Width of the playfield in worldspace
	height    float32 // Height of the playfield in worldspace
	Spaceship Spaceship
}

// Constants for gameplay feel
const (
	spaceshipRotateSpeed float32 = math.Pi * 3 // 1.5 rotations per second
)

func MakeWorld(width float32, height float32) World {
	w := World{
		Spaceship: MakeSpaceship(),
		width:     width,
		height:    height,
	}
	// Spaceship starts in the middle of the playfield
	w.Spaceship.Position = rl.Vector2{
		X: width / 2,
		Y: height / 2,
	}
	return w
}
