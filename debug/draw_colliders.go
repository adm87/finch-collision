package debug

import (
	"image/color"

	"github.com/adm87/finch-collision/collision"
	"github.com/adm87/finch-core/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawColliders(world *collision.CollisionWorld, screen *ebiten.Image, viewport geom.Rect64, viewMatrix ebiten.GeoM) {
	for collider := range world.QueryArea(viewport) {
		aabb := collider.AABB()
		drawRect(screen, &aabb, viewMatrix)
	}
}

func drawRect(screen *ebiten.Image, rect *geom.Rect64, viewMatrix ebiten.GeoM) {
	rectImg.Fill(color.RGBA{R: 255, A: 255})

	path := vector.Path{}

	minx, miny := rect.Min()
	maxx, maxy := rect.Max()

	sminx, sminy := viewMatrix.Apply(minx, miny)
	smaxx, smaxy := viewMatrix.Apply(maxx, maxy)

	path.MoveTo(float32(sminx), float32(sminy))
	path.LineTo(float32(smaxx), float32(sminy))
	path.LineTo(float32(smaxx), float32(smaxy))
	path.LineTo(float32(sminx), float32(smaxy))
	path.Close()

	vs, ic := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: 1,
	})

	screen.DrawTriangles(vs, ic, rectImg, nil)

	vs, ic = path.AppendVerticesAndIndicesForFilling(nil, nil)

	for i := range vs {
		vs[i].ColorA = 0.25
	}

	screen.DrawTriangles(vs, ic, rectImg, nil)
}
