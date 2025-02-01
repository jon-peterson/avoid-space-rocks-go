package gameobjects

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameObject interface {
	Update() error
	Draw() error
}

type GameObjectCollection struct {
	objects []GameObject
}

func NewGameObjectCollection() GameObjectCollection {
	return GameObjectCollection{
		objects: make([]GameObject, 0, 100),
	}
}

func (c *GameObjectCollection) Add(obj GameObject) {
	c.objects = append(c.objects, obj)
}

func (c *GameObjectCollection) Update() {
	for idx, obj := range c.objects {
		if err := obj.Update(); err != nil {
			rl.TraceLog(rl.LogError, "error updating object %d %v: %v", idx, obj, err)
		}
	}
}

func (c *GameObjectCollection) Draw() {
	for idx, obj := range c.objects {
		if err := obj.Draw(); err != nil {
			rl.TraceLog(rl.LogError, "error drawing object %d %v: %v", idx, obj, err)
		}
	}
}
