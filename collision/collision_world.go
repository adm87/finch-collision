package collision

import (
	"github.com/adm87/finch-core/fsys"
	"github.com/adm87/finch-core/geom"
	"github.com/adm87/finch-core/hashset"
)

const (
	MinMovementThreshold     = 0.001
	SweptCollisionStepFactor = 0.25
	MinSweptCollisionSteps   = 3
	MaxSweptCollisionSteps   = 50
)

type ContactInfo struct {
	ColliderA Collider
	ColliderB Collider
	Normal    geom.Point64
	Depth     float64
}

type CollisionWorld struct {
	profiles         CollisionProfile
	grid             *Grid
	dynamicTracking  map[Collider]geom.Point64
	dynamicColliders hashset.Set[Collider]
	colliders        hashset.Set[Collider]
}

func NewCollisionWorld(cellSize float64) *CollisionWorld {
	return &CollisionWorld{
		profiles:         make(CollisionProfile),
		grid:             NewGrid(cellSize),
		dynamicTracking:  make(map[Collider]geom.Point64),
		dynamicColliders: hashset.New[Collider](),
		colliders:        hashset.New[Collider](),
	}
}

func (c *CollisionWorld) SetProfiles(profiles CollisionProfile) {
	c.profiles = profiles
}

func (c *CollisionWorld) AddCollider(collider Collider) {
	if c.colliders.Contains(collider) {
		return
	}

	c.colliders.AddDistinct(collider)
	c.grid.Insert(collider)

	if collider.Type() == ColliderDynamic {
		c.dynamicColliders.AddDistinct(collider)
		c.dynamicTracking[collider] = geom.NewPoint64(collider.AABB().X, collider.AABB().Y)
	}
}

func (c *CollisionWorld) RemoveCollider(collider Collider) {
	if !c.colliders.Contains(collider) {
		return
	}

	delete(c.dynamicTracking, collider)
	delete(c.dynamicColliders, collider)
	delete(c.colliders, collider)

	c.grid.Remove(collider)
}

func (c *CollisionWorld) UpdateCollider(collider Collider) {
	if !c.colliders.Contains(collider) {
		// If the collider is not in the system, make sure it isn't being tracked
		c.grid.Remove(collider)
		if collider.Type() == ColliderDynamic {
			delete(c.dynamicTracking, collider)
			delete(c.dynamicColliders, collider)
		}
		return
	}

	// Handle changes in collider type
	if collider.Type() == ColliderDynamic && !c.dynamicColliders.Contains(collider) {
		c.dynamicColliders.AddDistinct(collider)
		c.dynamicTracking[collider] = geom.NewPoint64(collider.AABB().X, collider.AABB().Y)
	} else if collider.Type() == ColliderStatic && c.dynamicColliders.Contains(collider) {
		delete(c.dynamicColliders, collider)
		delete(c.dynamicTracking, collider)
	}

	// Reinsert the collider into the spatial grid
	c.grid.Reinsert(collider)
}

func (c *CollisionWorld) QueryArea(area geom.Rect64) hashset.Set[Collider] {
	cells := c.grid.GetCellsInArea(area)

	results := hashset.New[Collider]()
	for cell := range cells {
		for collider := range c.grid.cells[cell] {
			results.AddDistinct(collider)
		}
	}

	return results
}

func (c *CollisionWorld) Clear() {
	c.colliders.Clear()
	c.grid = NewGrid(c.grid.CellSize())
}

func (c *CollisionWorld) Grid() *Grid {
	return c.grid
}

func (c *CollisionWorld) CheckForCollisions(dt float64) {
	for colliderA := range c.dynamicColliders {
		layerA := colliderA.Layer()

		// Skip if there is no collision profile for the layer this collider is on
		if !c.profiles.HasProfile(layerA) {
			continue
		}

		queryArea := colliderA.AABB()
		if colliderA.DetectionType() == CollisionDetectionContinuous {
			if prevPos, exists := c.dynamicTracking[colliderA]; exists {
				currentAABB := colliderA.AABB()
				minX := min(prevPos.X, currentAABB.X)
				minY := min(prevPos.Y, currentAABB.Y)
				maxX := max(prevPos.X+currentAABB.Width, currentAABB.X+currentAABB.Width)
				maxY := max(prevPos.Y+currentAABB.Height, currentAABB.Y+currentAABB.Height)
				queryArea = geom.Rect64{X: minX, Y: minY, Width: maxX - minX, Height: maxY - minY}
			}
		}

		others := c.QueryArea(queryArea)
		if len(others) == 0 || (len(others) == 1 && others.Contains(colliderA)) {
			continue
		}
		others.Remove(colliderA)

		for other := range others {
			layerB := other.Layer()

			// Skip if there is no collision rule for the initiating collider's layer against the other collider's layer
			if !c.profiles[layerA].HasRule(layerB) {
				continue
			}

			var contact *ContactInfo
			var collided bool

			if colliderA.DetectionType() == CollisionDetectionContinuous {
				contact, collided = c.detectSweptCollision(colliderA, other)
			} else {
				contact, collided = c.detectCollision(colliderA.AABB(), other.AABB())
			}

			if !collided {
				continue
			}

			contact.ColliderA = colliderA
			contact.ColliderB = other

			c.profiles[layerA][layerB](contact)
		}

		if colliderA.Type() == ColliderDynamic {
			c.dynamicTracking[colliderA] = geom.NewPoint64(colliderA.AABB().X, colliderA.AABB().Y)
		}
	}
}

func (c *CollisionWorld) detectCollision(a, b geom.Rect64) (*ContactInfo, bool) {
	if !a.Intersects(b) {
		return nil, false
	}

	minxA, minyA := a.Min()
	maxxA, maxyA := a.Max()
	minxB, minyB := b.Min()
	maxxB, maxyB := b.Max()

	overlapX := min(maxxA, maxxB) - max(minxA, minxB)
	overlapY := min(maxyA, maxyB) - max(minyA, minyB)

	if overlapX <= 0 || overlapY <= 0 {
		return nil, false
	}

	var normal geom.Point64
	var depth float64

	if overlapX < overlapY {
		if (minxA + maxxA) < (minxB + maxxB) {
			normal = geom.Point64{X: -1, Y: 0}
		} else {
			normal = geom.Point64{X: 1, Y: 0}
		}
		depth = overlapX
	} else {
		if (minyA + maxyA) < (minyB + maxyB) {
			normal = geom.Point64{X: 0, Y: -1}
		} else {
			normal = geom.Point64{X: 0, Y: 1}
		}
		depth = overlapY
	}

	return &ContactInfo{
		Normal: normal,
		Depth:  depth,
	}, true
}

func (c *CollisionWorld) detectSweptCollision(collider, other Collider) (*ContactInfo, bool) {
	prevPos, exists := c.dynamicTracking[collider]
	if !exists {
		return c.detectCollision(collider.AABB(), other.AABB())
	}

	currentAABB := collider.AABB()
	otherAABB := other.AABB()

	currentPosition := geom.NewPoint64(currentAABB.X, currentAABB.Y)
	movement := currentPosition.Sub(prevPos)
	distance := movement.Length()

	// If no movement, just do regular collision detection
	if distance < MinMovementThreshold {
		return c.detectCollision(currentAABB, otherAABB)
	}

	stepSize := min(currentAABB.Width, currentAABB.Height) * SweptCollisionStepFactor
	steps := fsys.Clamp(int(distance/stepSize)+1, MinSweptCollisionSteps, MaxSweptCollisionSteps)
	stepVector := movement.Div(float64(steps))

	for i := 0; i <= steps; i++ {
		testAABB := geom.Rect64{
			X:      prevPos.X + stepVector.X*float64(i),
			Y:      prevPos.Y + stepVector.Y*float64(i),
			Width:  currentAABB.Width,
			Height: currentAABB.Height,
		}
		if contact, collided := c.detectCollision(testAABB, otherAABB); collided {
			return contact, true
		}
	}

	return nil, false
}
