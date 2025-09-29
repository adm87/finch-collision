package collision

import "github.com/adm87/finch-core/geom"

type ColliderType string

const (
	ColliderDynamic ColliderType = "dynamic"
	ColliderStatic  ColliderType = "static"
)

type CollisionDetectionType int

const (
	CollisionDetectionDiscrete   CollisionDetectionType = iota // Standard collision detection
	CollisionDetectionContinuous                               // Swept collision detection for fast-moving objects
)

type Collider interface {
	AABB() geom.Rect64 // Axis-Aligned Bounding Box: the minimal rectangle that fully contains the collider

	AddToLayer(layer CollisionLayer)
	RemoveFromLayer(layer CollisionLayer)

	LayerMask() CollisionLayer

	Type() ColliderType
	SetType(colliderType ColliderType)

	CollisionDetection() CollisionDetectionType
	SetCollisionDetection(collisionDetectionType CollisionDetectionType)
}
