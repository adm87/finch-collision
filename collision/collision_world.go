package collision

import (
	"github.com/adm87/finch-core/geom"
	"github.com/adm87/finch-core/hashset"
)

type CollisionWorld struct {
	grid *Grid
}

func NewCollisionWorld(cellSize float64) *CollisionWorld {
	return &CollisionWorld{
		grid: NewGrid(cellSize),
	}
}

func (c *CollisionWorld) AddCollider(rect *geom.Rect64) {
	c.grid.Insert(rect)
}

func (c *CollisionWorld) RemoveCollider(rect *geom.Rect64) {
	c.grid.Remove(rect)
}

func (c *CollisionWorld) UpdateCollider(rect *geom.Rect64) {
	c.grid.Reinsert(rect)
}

func (c *CollisionWorld) QueryArea(area *geom.Rect64) hashset.Set[*geom.Rect64] {
	cells := c.grid.GetCellsInArea(area)

	results := hashset.New[*geom.Rect64]()
	for cell := range cells {
		for collider := range c.grid.cells[cell] {
			results.AddDistinct(collider)
		}
	}

	return results
}

func (c *CollisionWorld) CheckCollision() {

}

func (c *CollisionWorld) Clear() {
	c.grid = NewGrid(c.grid.CellSize())
}

func (c *CollisionWorld) Grid() *Grid {
	return c.grid
}
