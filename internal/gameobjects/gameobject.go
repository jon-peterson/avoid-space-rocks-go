package gameobjects

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"sync"
)

type GameObject interface {
	Update() error
	Draw() error
	IsAlive() bool
}

type Collidable interface {
	OnCollision(other Collidable) error
	GetHitbox() rl.Rectangle
}

type GameObjectCollection struct {
	objects     []GameObject
	newObjects  []GameObject
	objectsLock sync.RWMutex
}

func NewGameObjectCollection() GameObjectCollection {
	return GameObjectCollection{
		objects:    make([]GameObject, 0, 100),
		newObjects: make([]GameObject, 0, 100),
	}
}

// Add adds a game object to the collection.
func (c *GameObjectCollection) Add(obj GameObject) {
	c.newObjects = append(c.newObjects, obj)
}

// Update all the objects in the collection. Removes dead objects, updates the rest.
// Checks for collisions between objects.
func (c *GameObjectCollection) Update() {
	// Remove all dead objects from the collection
	c.removeDead()
	c.birthNew()

	// Update all the remaining
	c.objectsLock.RLock()
	defer c.objectsLock.RUnlock()
	for idx, obj := range c.objects {
		if err := obj.Update(); err != nil {
			rl.TraceLog(rl.LogError, "error updating object %d %v: %v", idx, obj, err)
		}
	}

	// Check for collisions on all the collidable objects
	c.collisionCheck()
}

// removeDead removes all dead objects from the collection, after which the collection will be a full
// slice of alive objects.
func (c *GameObjectCollection) removeDead() {
	c.objectsLock.Lock()
	defer c.objectsLock.Unlock()

	for i := len(c.objects) - 1; i >= 0; i-- {
		if !c.objects[i].IsAlive() {
			// Replace the current element with the one at the end
			c.objects[i] = c.objects[len(c.objects)-1]
			c.objects = c.objects[:len(c.objects)-1]
		}
	}
}

// birthNew adds all the new objects to the collection.
func (c *GameObjectCollection) birthNew() {
	c.objectsLock.Lock()
	defer c.objectsLock.Unlock()

	c.objects = append(c.objects, c.newObjects...)
	c.newObjects = c.newObjects[:0]
}

// CollisionCheck checks for collisions between all the Collidable objects in the
// collection. When two objects collide they have their OnCollision methods called.
// TODO: Use spatial partitioning so it's not O(n^2)
func (c *GameObjectCollection) collisionCheck() {
	for i := len(c.objects) - 1; i >= 0; i-- {
		hammer := c.getCollidable(i)
		if hammer == nil {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			anvil := c.getCollidable(j)
			if anvil == nil {
				continue
			}
			if rl.CheckCollisionRecs(hammer.GetHitbox(), anvil.GetHitbox()) {
				if err := hammer.OnCollision(anvil); err != nil {
					rl.TraceLog(rl.LogError, "error handling collision between %d %v and %d %v: %v", i, hammer, j, anvil, err)
				}
				if err := anvil.OnCollision(hammer); err != nil {
					rl.TraceLog(rl.LogError, "error handling collision between %d %v and %d %v: %v", j, anvil, i, hammer, err)
				}
			}
		}
	}
}

// getCollidable returns an interface on a game object that is both collidable and alive, or nil.
func (c *GameObjectCollection) getCollidable(idx int) Collidable {
	obj := c.objects[idx]
	cast, ok := obj.(Collidable)
	if ok && obj.IsAlive() {
		return cast
	} else {
		return nil
	}
}

func (c *GameObjectCollection) Draw() {
	for idx, obj := range c.objects {
		if err := obj.Draw(); err != nil {
			rl.TraceLog(rl.LogError, "error drawing object %d %v: %v", idx, obj, err)
		}
	}
}
