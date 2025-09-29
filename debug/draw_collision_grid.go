package debug

import (
	"image/color"

	"github.com/adm87/finch-collision/collision"
	"github.com/adm87/finch-core/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawCollisionGrid(world *collision.CollisionWorld, screen *ebiten.Image, viewport geom.Rect64, viewMatrix ebiten.GeoM) {
	rectImg.Fill(color.RGBA{G: 255, A: 10})

	path := vector.Path{}
	grid := world.Grid()

	cells := grid.GetCellsInArea(&viewport)
	cellSize := grid.CellSize()

	for cell := range cells {
		minx := float64(cell.X) * cellSize
		miny := float64(cell.Y) * cellSize
		maxx := minx + cellSize
		maxy := miny + cellSize

		sminx, sminy := viewMatrix.Apply(minx, miny)
		smaxx, smaxy := viewMatrix.Apply(maxx, maxy)

		path.MoveTo(float32(sminx), float32(sminy))
		path.LineTo(float32(smaxx), float32(sminy))
		path.LineTo(float32(smaxx), float32(smaxy))
		path.LineTo(float32(sminx), float32(smaxy))
		path.Close()
	}

	vs, ic := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: 1,
	})

	screen.DrawTriangles(vs, ic, rectImg, nil)

	vs, ic = path.AppendVerticesAndIndicesForFilling(nil, nil)

	for i := range vs {
		vs[i].ColorA = 0.15
	}

	screen.DrawTriangles(vs, ic, rectImg, nil)
}
