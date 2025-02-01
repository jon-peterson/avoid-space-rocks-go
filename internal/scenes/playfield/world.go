package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// The World object represents the state of the game within the playfield
type World struct {
	width     float32 // Width of the playfield in worldspace
	height    float32 // Height of the playfield in worldspace
	Spaceship Spaceship
	Objects   gameobjects.GameObjectCollection
}

func NewWorld(width float32, height float32) World {
	w := World{
		width:     width,
		height:    height,
		Objects:   gameobjects.NewGameObjectCollection(),
		Spaceship: NewSpaceship(),
	}
	// Spaceship starts in the middle of the playfield and is always first
	w.Spaceship.Position = rl.Vector2{
		X: width / 2,
		Y: height / 2,
	}
	return w
}

// Wraparound returns the position of the given position, wrapping around the edges of the playfield
func (w *World) Wraparound(p rl.Vector2) rl.Vector2 {
	if p.X < 0 {
		p.X = w.width
	} else if p.X > w.width {
		p.X = 0
	}
	if p.Y < 0 {
		p.Y = w.height
	} else if p.Y > w.height {
		p.Y = 0
	}
	return p
}
