package colliders

import (
	"github.com/adm87/finch-collision/collision"
	"github.com/adm87/finch-core/geom"
)

type BoxCollider struct {
	geom.Rect64

	colliderType       collision.ColliderType
	collisionLayer     collision.CollisionLayer
	collisionDetection collision.CollisionDetectionType
}

func NewBoxCollider(x, y, width, height float64) *BoxCollider {
	return &BoxCollider{
		Rect64: geom.Rect64{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
		colliderType:       collision.ColliderStatic,
		collisionDetection: collision.CollisionDetectionDiscrete,
		collisionLayer:     0,
	}
}

func (b *BoxCollider) AABB() geom.Rect64 {
	return b.Rect64
}

func (b *BoxCollider) Shape() any {
	return b.Rect64
}

func (b *BoxCollider) Layer() collision.CollisionLayer {
	return b.collisionLayer
}

func (b *BoxCollider) SetLayer(layer collision.CollisionLayer) {
	b.collisionLayer = layer
}

func (b *BoxCollider) Type() collision.ColliderType {
	return b.colliderType
}

func (b *BoxCollider) SetType(colliderType collision.ColliderType) {
	b.colliderType = colliderType
}

func (b *BoxCollider) DetectionType() collision.CollisionDetectionType {
	return b.collisionDetection
}

func (b *BoxCollider) SetDetectionType(detectionType collision.CollisionDetectionType) {
	b.collisionDetection = detectionType
}
