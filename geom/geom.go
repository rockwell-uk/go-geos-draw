package geom

import (
	"fmt"

	geos "github.com/twpayne/go-geos"
)

var gctx = geos.NewContext()

var lineStringEmpty, _ = gctx.NewGeomFromWKT("LINESTRING EMPTY")
var multiLineStringEmpty, _ = gctx.NewGeomFromWKT("MULTILINESTRING EMPTY")
var polygonEmpty, _ = gctx.NewGeomFromWKT("POLYGON EMPTY")

func SimplifyGeom(g *geos.Geom, lvl float64) *geos.Geom {
	return g.Simplify(lvl)
}

//nolint:exhaustive
func GetGeometryCenter(g *geos.Geom, scale func(x, y float64) (float64, float64)) ([]float64, error) {
	var c []float64
	var err error

	switch g.TypeID() {
	case geos.TypeIDPoint:
		x, y := scale(g.X(), g.Y())

		return []float64{
			x,
			y,
		}, nil

	case geos.TypeIDMultiLineString, geos.TypeIDLineString:
		g, err = ScaleLine(g, scale)
		if err != nil {
			return []float64{}, err
		}

		c = CenterFromGeometry(g)

	default:
		return c, fmt.Errorf("geom type not supported %v", g.TypeID())
	}

	return c, nil
}

func CenterFromGeometry(g *geos.Geom) []float64 {
	var Xmin, Ymin, Xmax, Ymax float64

	b := g.Bounds()

	Xmin = b.MinX
	Ymin = b.MinY
	Xmax = b.MaxX
	Ymax = b.MaxY

	x := Xmin + ((Xmax - Xmin) / 2)
	y := Ymin + ((Ymax - Ymin) / 2)

	return []float64{
		x,
		y,
	}
}

func noscale(x, y float64) (float64, float64) {
	return x, y
}

func ScaleLine(g *geos.Geom, scale func(x, y float64) (float64, float64)) (*geos.Geom, error) {
	r := "MULTILINESTRING"

	t, err := transform(r, g, true, scale)
	if err != nil {
		return multiLineStringEmpty, err
	}

	return t, nil
}

func ToLineString(g *geos.Geom) (*geos.Geom, error) {
	r := "LINESTRING"

	t, err := transform(r, g, false, noscale)
	if err != nil {
		return lineStringEmpty, err
	}

	return t, nil
}

func ToPolygon(g *geos.Geom) (*geos.Geom, error) {
	r := "POLYGON"

	t, err := transform(r, g, false, noscale)
	if err != nil {
		return polygonEmpty, err
	}

	return t, nil
}

func transform(gType string, g *geos.Geom, multi bool, scale func(x, y float64) (float64, float64)) (*geos.Geom, error) {
	s := transformToPoints(gType, g, multi, scale)
	if s == "" {
		return lineStringEmpty, nil
	}

	p, err := gctx.NewGeomFromWKT(s)
	if err != nil {
		return lineStringEmpty, fmt.Errorf("transform %v: %v", err.Error(), s)
	}

	return p, nil
}

func transformToPoints(r string, g *geos.Geom, multi bool, scale func(x, y float64) (float64, float64)) string {
	var startSep, endSep string

	if multi {
		startSep = "(("
		endSep = "))"
	} else {
		startSep = "("
		endSep = ")"
	}
	r = fmt.Sprintf("%v %v", r, startSep)
	points := GetPoints(g)
	l := len(*points)
	for i, p := range *points {
		x, y := scale(p[0], p[1])
		r = fmt.Sprintf("%v%v %v", r, x, y)
		if i < l-1 {
			r = fmt.Sprintf("%v,", r)
		}
	}
	r = fmt.Sprintf("%v%v", r, endSep)

	return r
}

func GetOrd(cs *geos.CoordSeq, fn func(*geos.CoordSeq, int) float64) func(int) float64 {
	return func(idx int) float64 {
		ord := fn(cs, idx)
		return ord
	}
}

func BoundsWKT(xmin, xmax, ymin, ymax float64) string {
	tl := []float64{
		xmin,
		ymin,
	}
	tr := []float64{
		xmax,
		ymin,
	}
	bl := []float64{
		xmin,
		ymax,
	}
	br := []float64{
		xmax,
		ymax,
	}

	return fmt.Sprintf(
		"POLYGON ((%v %v, %v %v, %v %v, %v %v, %v %v))",
		tl[0], tl[1], tr[0], tr[1], br[0], br[1], bl[0], br[1], tl[0], tl[1],
	)
}

func BoundsGeom(xmin, xmax, ymin, ymax float64) (*geos.Geom, error) {
	wkt := BoundsWKT(xmin, xmax, ymin, ymax)

	g, err := gctx.NewGeomFromWKT(wkt)
	if err != nil {
		return polygonEmpty, err
	}

	return g, nil
}

func GetPoints(list ...*geos.Geom) *[][]float64 {
	var res [][]float64

	for _, g := range list {
		_type := g.TypeID()

		switch _type {
		case geos.TypeIDPoint, geos.TypeIDLineString, geos.TypeIDLinearRing, geos.TypeIDMultiPoint, geos.TypeIDMultiLineString:
			for i := 0; i < g.NumGeometries(); i++ {
				res = append(res, g.Geometry(i).CoordSeq().ToCoords()...)
			}
		case geos.TypeIDPolygon:
			l := g.ExteriorRing()
			for i := 0; i < l.NumGeometries(); i++ {
				res = append(res, l.Geometry(i).CoordSeq().ToCoords()...)
			}
		case geos.TypeIDMultiPolygon:
			for i := 0; i < g.NumGeometries(); i++ {
				t := g.Geometry(i)
				l := t.ExteriorRing()
				for i := 0; i < l.NumGeometries(); i++ {
					res = append(res, l.Geometry(i).CoordSeq().ToCoords()...)
				}
			}
		case geos.TypeIDGeometryCollection:
			var subgeoms []*geos.Geom
			for i := 0; i < g.NumGeometries(); i++ {
				subgeoms = append(subgeoms, g.Geometry(i))
			}
			s := GetPoints(subgeoms...)
			res = append(res, *s...)
		}
	}

	return &res
}
