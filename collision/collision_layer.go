package collision

type CollisionLayer uint64

const NoCollisionLayers CollisionLayer = 0

func (cm CollisionLayer) HasLayer(layer CollisionLayer) bool {
	return (cm & layer) != 0
}

func (cm *CollisionLayer) AddLayer(layer CollisionLayer) CollisionLayer {
	(*cm) |= layer
	return *cm
}

func (cm *CollisionLayer) RemoveLayer(layer CollisionLayer) CollisionLayer {
	return *cm &^ layer
}

func (cm CollisionLayer) Intersects(other CollisionLayer) bool {
	return (cm & other) != 0
}
