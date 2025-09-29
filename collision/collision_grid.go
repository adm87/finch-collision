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
	cells    map[GridKey]hashset.Set[*geom.Rect64]

	cellsByCollider map[*geom.Rect64]hashset.Set[GridKey]
}

func NewGrid(cellSize float64) *Grid {
	return &Grid{
		cellSize:        cellSize,
		cells:           make(map[GridKey]hashset.Set[*geom.Rect64]),
		cellsByCollider: make(map[*geom.Rect64]hashset.Set[GridKey]),
	}
}

func (g *Grid) CellSize() float64 {
	return g.cellSize
}

func (g *Grid) Insert(rect *geom.Rect64) {
	if rect.Width <= 0 || rect.Height <= 0 {
		panic("cannot insert rect with non-positive width or height")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	keys := g.getGridKeys(rect)
	for _, key := range keys {
		if _, ok := g.cells[key]; !ok {
			g.cells[key] = hashset.New[*geom.Rect64]()
		}

		g.cells[key].AddDistinct(rect)

		if _, ok := g.cellsByCollider[rect]; !ok {
			g.cellsByCollider[rect] = hashset.New[GridKey]()
		}

		g.cellsByCollider[rect].AddDistinct(key)
	}
}

func (g *Grid) Reinsert(rect *geom.Rect64) {
	g.Remove(rect)
	g.Insert(rect)
}

func (g *Grid) Remove(rect *geom.Rect64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if keys, ok := g.cellsByCollider[rect]; ok {
		for key := range keys {
			if cell, ok := g.cells[key]; ok {
				cell.Remove(rect)
				if len(cell) == 0 {
					delete(g.cells, key)
				}
			}
		}
		delete(g.cellsByCollider, rect)
	}
}

func (g *Grid) GetCellsInArea(area *geom.Rect64) hashset.Set[GridKey] {
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

func (g *Grid) getGridKeys(rect *geom.Rect64) []GridKey {
	keys := make([]GridKey, 0)

	minx, miny := rect.Min()
	maxx, maxy := rect.Max()

	offset := g.cellSize / 2.0

	startCell := geom.Point64{
		X: math.Floor((minx - offset) / g.cellSize),
		Y: math.Floor((miny - offset) / g.cellSize),
	}
	endCell := geom.Point64{
		X: math.Floor((maxx + offset) / g.cellSize),
		Y: math.Floor((maxy + offset) / g.cellSize),
	}

	for x := startCell.X; x <= endCell.X; x++ {
		for y := startCell.Y; y <= endCell.Y; y++ {
			keys = append(keys, GridKey{X: x, Y: y})
		}
	}

	return keys
}
