package collision

import (
	"math"
	"sync"

	"github.com/adm87/finch-core/geom"
	"github.com/adm87/finch-core/hashset"
)

type GridKey geom.Point64

type Grid struct {
	mu sync.RWMutex

	cellSize float64
	cells    map[GridKey]hashset.Set[Collider]

	cellsByCollider map[Collider]hashset.Set[GridKey]
}

func NewGrid(cellSize float64) *Grid {
	return &Grid{
		cellSize:        cellSize,
		cells:           make(map[GridKey]hashset.Set[Collider]),
		cellsByCollider: make(map[Collider]hashset.Set[GridKey]),
	}
}

func (g *Grid) CellSize() float64 {
	return g.cellSize
}

func (g *Grid) Insert(collider Collider) {
	aabb := collider.AABB()

	if aabb.Width <= 0 || aabb.Height <= 0 {
		panic("cannot insert collider with non-positive width or height bounding box")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	keys := g.getGridKeys(aabb)
	for _, key := range keys {
		if _, ok := g.cells[key]; !ok {
			g.cells[key] = hashset.New[Collider]()
		}

		g.cells[key].AddDistinct(collider)

		if _, ok := g.cellsByCollider[collider]; !ok {
			g.cellsByCollider[collider] = hashset.New[GridKey]()
		}

		g.cellsByCollider[collider].AddDistinct(key)
	}
}

func (g *Grid) Reinsert(collider Collider) {
	g.Remove(collider)
	g.Insert(collider)
}

func (g *Grid) Remove(collider Collider) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if keys, ok := g.cellsByCollider[collider]; ok {
		for key := range keys {
			if cell, ok := g.cells[key]; ok {
				delete(cell, collider)
				if len(cell) == 0 {
					delete(g.cells, key)
				}
			}
		}
		delete(g.cellsByCollider, collider)
	}
}

func (g *Grid) GetCellsInArea(area geom.Rect64) hashset.Set[GridKey] {
	g.mu.RLock()
	defer g.mu.RUnlock()

	keys := g.getGridKeys(area)
	cells := hashset.New[GridKey]()

	for _, key := range keys {
		if _, ok := g.cells[key]; ok {
			cells.AddDistinct(key)
		}
	}

	return cells
}

func (g *Grid) getGridKeys(rect geom.Rect64) []GridKey {
	minx, miny := rect.Min()
	maxx, maxy := rect.Max()

	offset := g.cellSize / 2.0

	startX := math.Floor((minx - offset) / g.cellSize)
	startY := math.Floor((miny - offset) / g.cellSize)
	endX := math.Floor((maxx + offset) / g.cellSize)
	endY := math.Floor((maxy + offset) / g.cellSize)

	keys := make([]GridKey, 0, int((endX-startX+1)*(endY-startY+1)))
	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {
			keys = append(keys, GridKey{X: x, Y: y})
		}
	}
	return keys
}
