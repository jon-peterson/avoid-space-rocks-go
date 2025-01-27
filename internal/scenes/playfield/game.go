package playfield

import "sync"

var instance *Game
var once sync.Once

type Game struct {
	World World
	Lives uint8
	Level uint8
	Score uint64
}

func GetGame() *Game {
	if instance == nil {
		panic("Game not initialized. Call InitGame first")
	}
	return instance
}

func InitGame(w World) *Game {
	once.Do(func() {
		instance = &Game{
			World: w,
			Lives: 3,
			Level: 1,
			Score: 0,
		}
	})
	return instance
}
