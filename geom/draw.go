package geom

import (
	"errors"
	"image/color"
	"math"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	geos "github.com/twpayne/go-geos"
	"golang.org/x/image/font"
)

func DrawString(gc *draw2dimg.GraphicContext, pos []float64, rotation float64, text string) error {
	radians := rotation * (math.Pi / 180)
	rm := draw2d.NewRotationMatrix(radians)

	trx, try := rm.InverseTransformPoint(pos[0], pos[1])

	gc.Save()
	gc.Rotate(radians)
	gc.StrokeStringAt(text, trx, try)
	gc.FillStringAt(text, trx, try)
	gc.Restore()

	return nil
}

func DrawRune(gc *draw2dimg.GraphicContext, pos []float64, f font.Face, rotation float64, char rune) error {
	radians := rotation * (math.Pi / 180)
	rm := draw2d.NewRotationMatrix(radians)

	trlx, trly := rm.InverseTransformPoint(pos[0], pos[1])

	gc.Save()
	gc.Rotate(radians)
	gc.StrokeStringAt(string(char), trlx, trly)
	gc.FillStringAt(string(char), trlx, trly)
	gc.Restore()

	return nil
}

func DrawPoint(gc draw2d.GraphicContext, g *geos.Geom, radius float64, fillColor color.Color, strokeWidth float64, strokeColor color.Color, scale func(x, y float64) (float64, float64)) error {
	gc.SetFillColor(fillColor)
	gc.SetStrokeColor(strokeColor)
	gc.SetLineWidth(strokeWidth)

	x := g.X()
	y := g.Y()

	x, y = scale(x, y)

	draw2dkit.Circle(gc, x, y, radius)
	gc.FillStroke()

	return nil
}

func DrawLine(gc draw2d.GraphicContext, g *geos.Geom, lineWidth float64, fillColor color.Color, strokeWidth float64, strokeColor color.Color, scale func(x, y float64) (float64, float64)) error {
	cs := GetPoints(g)

	if lineWidth == 0.0 {
		return errors.New("line width cannot be zero")
	}

	// first line is for the stroke (beneath)
	gc.SetStrokeColor(strokeColor)
	gc.SetLineWidth(lineWidth + strokeWidth)

	err := lineCoordSeq(gc, cs, scale)
	if err != nil {
		return (err)
	}
	gc.Stroke()

	// the actual line
	gc.SetStrokeColor(fillColor)
	gc.SetLineWidth(lineWidth)

	err = lineCoordSeq(gc, cs, scale)
	if err != nil {
		return (err)
	}
	gc.Stroke()

	return nil
}

func DrawCoordLine(gc draw2d.GraphicContext, lineCoords [][]float64, lineWidth float64, fillColor color.Color, strokeWidth float64, strokeColor color.Color, scale func(x, y float64) (float64, float64)) error {
	if lineWidth == 0.0 {
		return errors.New("line width cannot be zero")
	}

	// first line is for the stroke (beneath)
	gc.SetStrokeColor(strokeColor)
	gc.SetLineWidth(lineWidth + strokeWidth)

	err := lineCoordSeq(gc, &lineCoords, scale)
	if err != nil {
		return (err)
	}
	gc.Stroke()

	// the actual line
	gc.SetStrokeColor(fillColor)
	gc.SetLineWidth(lineWidth)

	err = lineCoordSeq(gc, &lineCoords, scale)
	if err != nil {
		return (err)
	}
	gc.Stroke()

	return nil
}

func DrawPolygon(gc draw2d.GraphicContext, g *geos.Geom, fillColor color.Color, strokeColor color.Color, strokeWidth float64, scale func(x, y float64) (float64, float64)) error {
	gc.SetFillColor(fillColor)
	gc.SetStrokeColor(strokeColor)
	gc.SetLineWidth(strokeWidth)

	// exterior ring
	cs := GetPoints(g.ExteriorRing())

	err := lineCoordSeq(gc, cs, scale)
	if err != nil {
		return (err)
	}
	gc.FillStroke()
	// interior rings...

	n := g.NumInteriorRings()
	for i := 0; i < n; i++ {
		err = lineCoordSeq(gc, GetPoints(g.InteriorRing(i)), scale)
		if err != nil {
			return (err)
		}
		gc.SetFillColor(color.White)
		gc.FillStroke()
	}

	return nil
}

func DrawDot(gc *draw2dimg.GraphicContext, radius, x, y float64) error {
	gc.MoveTo(x, y)
	gc.ArcTo(x, y, radius, radius, 0, 2*math.Pi)
	gc.Fill()

	return nil
}

func lineCoordSeq(gc draw2d.GraphicContext, cs *[][]float64, scale func(x, y float64) (float64, float64)) error {
	if cs == nil {
		return errors.New("coord seq cannot be nil")
	}

	csd := *cs

	gc.MoveTo(scale(csd[0][0], csd[0][1]))

	for i := 1; i < len(csd); i++ {
		x, y := scale(csd[i][0], csd[i][1])
		gc.LineTo(x, y)
	}

	return nil
}
