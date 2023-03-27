package geom

import (
	geos "github.com/twpayne/go-geos"
)

type Envelope struct {
	Min, Max []float64
}

func (e Envelope) Dx() float64 {
	return e.Max[0] - e.Min[0]
}

func (e Envelope) Dy() float64 {
	return e.Max[1] - e.Min[1]
}

func (e Envelope) Px(x float64) float64 {
	return (x - e.Min[0]) / e.Dx()
}

func (e Envelope) Py(y float64) float64 {
	return (y - e.Min[1]) / e.Dy()
}

func ToEnvelope(g *geos.Geom) (Envelope, error) {
	l, err := ToLineString(g)
	if err != nil {
		return Envelope{}, err
	}

	bl := l.Point(4)
	tr := l.Point(2)

	minX := bl.X()
	maxX := tr.X()
	minY := bl.Y()
	maxY := tr.Y()

	return Envelope{
		[]float64{
			minX,
			minY,
		},
		[]float64{
			maxX,
			maxY,
		},
	}, nil
}
