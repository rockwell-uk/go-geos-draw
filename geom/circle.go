package geom

import (
	"fmt"

	"github.com/rockwell-uk/go-draw/draw"
	geos "github.com/twpayne/go-geos"
)

func CircleWKT(origin []float64, radius float64, numPoints int) (string, error) {
	points, err := draw.Circle(origin, radius, numPoints)
	if err != nil {
		return "LINESTRING EMPTY", err
	}

	n := len(points)
	s := "LINESTRING("
	for i, c := range points {
		s = fmt.Sprintf("%v%v %v", s, c[0], c[1])
		if i < n-1 {
			s = fmt.Sprintf("%v,", s)
		}
	}

	s = fmt.Sprintf("%v)", s)

	return s, nil
}

func CircleGeom(origin []float64, radius float64, numPoints int) (*geos.Geom, error) {
	e := geos.Geom{}

	r, err := CircleWKT(origin, radius, numPoints)
	if err != nil {
		return &e, err
	}

	g, err := gctx.NewGeomFromWKT(r)
	if err != nil {
		return &e, err
	}

	return g, nil
}
