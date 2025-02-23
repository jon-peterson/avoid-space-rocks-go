package scenes

type Scene interface {
	Init(width, height float32)
	Loop() SceneCode
	Close()
}

type SceneCode int

const (
	AttractModeScene SceneCode = iota
	GameplayScene
	GameOverScene
	Quit
)
