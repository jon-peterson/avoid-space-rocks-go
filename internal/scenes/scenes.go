package scenes

type Scene interface {
	Init()
	Loop() SceneCode
	Close()
}

type SceneCode int

const (
	AttractMode SceneCode = iota
	Gameplay
	GameOver
)
