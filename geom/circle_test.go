package geom

import (
	"fmt"
	"image"
	"image/draw"
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	godraw "github.com/rockwell-uk/go-draw/draw"
)

func TestCircle(t *testing.T) {
	m := image.NewRGBA(image.Rect(0, 0, 600.0, 600.0))
	draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.Point{0, 0}, draw.Src)
	gc := draw2dimg.NewGraphicContext(m)

	gc.SetDPI(72)

	origin := []float64{
		300.00,
		300.00,
	}
	radius := 280.0
	numPoints := 100

	points, err := godraw.Circle(
		origin,
		radius,
		numPoints,
	)
	if err != nil {
		t.Fatal(err)
	}

	scale := func(x, y float64) (float64, float64) {
		return x, y
	}

	fillColor := black
	strokeWidth := 0.0
	strokeColor := black
	pointRadius := 1.0

	for _, p := range points {
		g, err := gctx.NewGeomFromWKT(fmt.Sprintf("POINT(%v %v)", p[0], p[1]))
		if err != nil {
			t.Fatal(err)
		}
		err = DrawPoint(gc, g, pointRadius, fillColor, strokeWidth, strokeColor, scale)
		if err != nil {
			t.Fatal(err)
		}
	}

	// draw the image
	err = savePNG("test-output/circle.png", m)
	if err != nil {
		t.Fatal(err)
	}
}
