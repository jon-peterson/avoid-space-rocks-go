package gameobjects

import rl "github.com/gen2brain/raylib-go/raylib"

type GameObject struct {
	Body
	SpriteSheet rl.Texture2D
}

// TODO: add filecheck and panic
func MakeGameObject(spritesheet string) GameObject {
	return GameObject{
		SpriteSheet: rl.LoadTexture("assets/sprites/" + spritesheet),
	}
}
