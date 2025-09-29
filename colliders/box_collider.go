package colliders

import (
	"github.com/adm87/finch-collision/collision"
	"github.com/adm87/finch-core/geom"
)

type BoxCollider struct {
	geom.Rect64

	layerMask          collision.CollisionLayer
	colliderType       collision.ColliderType
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
		layerMask:          collision.NoCollisionLayers,
		colliderType:       collision.ColliderStatic,
		collisionDetection: collision.CollisionDetectionDiscrete,
	}
}

func (b *BoxCollider) AABB() geom.Rect64 {
	return b.Rect64
}

func (b *BoxCollider) Shape() any {
	return b.Rect64
}

func (b *BoxCollider) AddToLayer(layer collision.CollisionLayer) {
	b.layerMask.AddLayer(layer)
}

func (b *BoxCollider) RemoveFromLayer(layer collision.CollisionLayer) {
	b.layerMask.RemoveLayer(layer)
}

func (b *BoxCollider) LayerMask() collision.CollisionLayer {
	return b.layerMask
}

func (b *BoxCollider) Type() collision.ColliderType {
	return b.colliderType
}

func (b *BoxCollider) SetType(colliderType collision.ColliderType) {
	b.colliderType = colliderType
}

func (b *BoxCollider) CollisionDetection() collision.CollisionDetectionType {
	return b.collisionDetection
}

func (b *BoxCollider) SetCollisionDetection(collisionDetectionType collision.CollisionDetectionType) {
	b.collisionDetection = collisionDetectionType
}
