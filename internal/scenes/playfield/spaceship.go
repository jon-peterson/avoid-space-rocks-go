package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Spaceship struct {
	gameobjects.Transform
	SpriteSheet rl.Texture2D
	frameWidth  int32
	frameHeight int32
}

func MakeSpaceship() Spaceship {
	ship := Spaceship{
		SpriteSheet: rl.LoadTexture("assets/sprites/spaceship.png"),
	}
	ship.frameWidth = ship.SpriteSheet.Width / 3
	ship.frameHeight = ship.SpriteSheet.Height
	return ship
}
