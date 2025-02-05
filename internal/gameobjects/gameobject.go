package gameobjects

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameObject interface {
	Update() error
	Draw() error
	IsAlive() bool
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
	// Remove all dead objects from the collection
	for i := len(c.objects) - 1; i >= 0; i-- {
		if !c.objects[i].IsAlive() {
			// Replace the current element with the one at the end
			c.objects[i] = c.objects[len(c.objects)-1]
			c.objects = c.objects[:len(c.objects)-1]
		}
	}

	// Update all the remaining
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
