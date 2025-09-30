package collision

import (
	"sync"

	"github.com/adm87/finch-core/enum"
)

// ======================================================
// Collider Type
// ======================================================

type ColliderType int

const (
	ColliderDynamic ColliderType = iota
	ColliderStatic
)

func (ct ColliderType) String() string {
	switch ct {
	case ColliderDynamic:
		return "Dynamic"
	case ColliderStatic:
		return "Static"
	default:
		return "Unknown"
	}
}

func (ct ColliderType) IsValid() bool {
	return ct >= ColliderDynamic && ct <= ColliderStatic
}

func (ct ColliderType) MarshalJSON() ([]byte, error) {
	return enum.MarshalEnum(ct)
}

func (ct *ColliderType) UnmarshalJSON(data []byte) error {
	val, err := enum.UnmarshalEnum[ColliderType](data)
	if err != nil {
		return err
	}
	*ct = val
	return nil
}

// ======================================================
// Collision Detection Type
// ======================================================

type CollisionDetectionType int

const (
	CollisionDetectionDiscrete   CollisionDetectionType = iota // Standard collision detection
	CollisionDetectionContinuous                               // Swept collision detection for fast-moving objects
)

func (cdt CollisionDetectionType) String() string {
	switch cdt {
	case CollisionDetectionDiscrete:
		return "Discrete"
	case CollisionDetectionContinuous:
		return "Continuous"
	default:
		return "Unknown"
	}
}

func (cdt CollisionDetectionType) IsValid() bool {
	return cdt >= CollisionDetectionDiscrete && cdt <= CollisionDetectionContinuous
}

func (cdt CollisionDetectionType) MarshalJSON() ([]byte, error) {
	return enum.MarshalEnum(cdt)
}

func (cdt *CollisionDetectionType) UnmarshalJSON(data []byte) error {
	val, err := enum.UnmarshalEnum[CollisionDetectionType](data)
	if err != nil {
		return err
	}
	*cdt = val
	return nil
}

// ======================================================
// Collision Layers - User Defined
// ======================================================

var (
	layersByName  = make(map[string]CollisionLayer)
	layersByValue = make(map[CollisionLayer]string)
	layerMu       = sync.RWMutex{}
)

type CollisionLayer int

// NewCollisionLayer creates a new collision layer with the given name.
// It panics if a layer with the same name already exists.
//
// Example usage:
//
//	playerLayer := collision.NewCollisionLayer("Player")
//	enemyLayer := collision.NewCollisionLayer("Enemy")
//	platformLayer := collision.NewCollisionLayer("Platform")
//
// The created layers can then be used to set up collision profiles for colliders.
//
// Important: When serializing and deserializing collision layers, it is recommended to use their names
// rather than their integer values. This ensures that the correct layers are referenced even if the
// order of layer creation changes in the code.
func NewCollisionLayer(name string) CollisionLayer {
	layerMu.Lock()
	defer layerMu.Unlock()

	length := len(layersByName)

	if _, exists := layersByName[name]; exists {
		panic("collision layer with this name already exists: " + name)
	}

	layer := CollisionLayer(length)
	layersByName[name] = layer
	layersByValue[layer] = name

	return layer
}

func (cl CollisionLayer) String() string {
	layerMu.RLock()
	defer layerMu.RUnlock()

	if name, exists := layersByValue[cl]; exists {
		return name
	}
	return "Unknown"
}

func (cl CollisionLayer) IsValid() bool {
	layerMu.RLock()
	defer layerMu.RUnlock()

	_, exists := layersByValue[cl]
	return exists
}

func (cl CollisionLayer) MarshalJSON() ([]byte, error) {
	return enum.MarshalEnum(cl)
}

func (cl *CollisionLayer) UnmarshalJSON(data []byte) error {
	val, err := enum.UnmarshalEnum[CollisionLayer](data)
	if err != nil {
		return err
	}
	*cl = val
	return nil
}
