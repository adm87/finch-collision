package collision

// CollisionResponse defines the function signature for responding to a collision event.
type CollisionResponse func(contact *ContactInfo)

// CollisionRules defines how a specific layer responds to collisions with other layers.
type CollisionRules map[CollisionLayer]CollisionResponse

func (rule CollisionRules) HasRule(layer CollisionLayer) bool {
	_, exists := rule[layer]
	return exists
}

// CollisionProfile defines the complete set of collision rules for all layers.
type CollisionProfile map[CollisionLayer]CollisionRules

func (profile CollisionProfile) HasProfile(layer CollisionLayer) bool {
	_, exists := profile[layer]
	return exists
}
