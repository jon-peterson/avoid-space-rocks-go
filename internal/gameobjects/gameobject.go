package gameobjects

type GameObject interface {
	Update() error
	Draw() error
}
