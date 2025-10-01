package collision

import "github.com/adm87/finch-core/geom"

type Collider interface {
	AABB() geom.Rect64 // Axis-Aligned Bounding Box: the minimal rectangle that fully contains the collider

	Layer() CollisionLayer
	SetLayer(layer CollisionLayer)

	Type() ColliderType
	SetType(colliderType ColliderType)

	DetectionType() CollisionDetectionType
	SetDetectionType(detectionType CollisionDetectionType)
}
