package geom

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
)

var (
	black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	white = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	blue  = color.RGBA{0x4C, 0x94, 0xFF, 0xFF}
)

func TestScale(t *testing.T) {
	tests := map[string]struct {
		tileWidth  float64
		tileHeight float64
		envelope   Envelope
		wkt        string
		expected   [][]float64
	}{
		"600 x 600": {
			float64(600),
			float64(600),
			Envelope{Min: []float64{387221.19853198, 410715.07842109}, Max: []float64{392221.19853198, 415715.07842109}},
			"LINESTRING(390380 413999.9999997685,390350 413979.99999976787,390226 413917.9999997675,390200.99999999994 413904.99999976775,390136.99999999994 413868.999999768,390104.99999999994 413850.9999997678,390280.99999999994 413942.9999997688,390259.00000000006 413932.9999997691,390350 413979.99999976787,390334 413969.9999997699,390200.99999999994 413904.99999976775,390191.00000000006 413899.9999997673,390334 413969.9999997699,390280.99999999994 413942.9999997688,390191.00000000006 413899.9999997673,390136.99999999994 413868.999999768,390104.99999999994 413850.9999997678,390062 413826.9999997678,390259.00000000006 413932.9999997691,390226 413917.9999997675,390540 414119.99999976903,390451.00000000006 414048.9999997688,390380 413999.9999997685,390688 414508.999999768,390644 414366.9999997688,390636 414336.99999976816,390576.99999999994 414167.9999997685,390576.99999999994 414167.9999997685,390564 414143.9999997677,390712.99999999994 414600.99999976816,390688 414508.999999768,390644 414366.9999997688,390636 414336.99999976816,390564 414143.9999997677,390540 414119.99999976903)",
			[][]float64{
				{379.05617616240164, 205.8094105585758},
				{375.4561761624017, 208.20941055865256},
				{360.5761761624017, 215.6494105586945},
				{357.5761761623947, 217.2094105586666},
				{349.8961761623947, 221.52941055863857},
				{346.0561761623947, 223.68941055865963},
				{367.1761761623947, 212.64941055854086},
				{364.53617616240865, 213.84941055850595},
				{375.4561761624017, 208.20941055865256},
				{373.53617616240166, 209.40941055840813},
				{357.5761761623947, 217.2094105586666},
				{356.3761761624087, 217.80941055872245},
				{373.53617616240166, 209.40941055840813},
				{367.1761761623947, 212.64941055854086},
				{356.3761761624087, 217.80941055872245},
				{349.8961761623947, 221.52941055863857},
				{346.0561761623947, 223.68941055865963},
				{340.8961761624017, 226.56941055865957},
				{364.53617616240865, 213.84941055850595},
				{360.5761761624017, 215.6494105586945},
				{398.2561761624017, 191.4094105585129},
				{387.5761761624087, 199.92941055854084},
				{379.05617616240164, 205.8094105585758},
				{416.0161761624017, 144.72941055863868},
				{410.7361761624017, 161.76941055854087},
				{409.77617616240167, 165.3694105586177},
				{402.69617616239475, 185.6494105585757},
				{402.69617616239475, 185.6494105585757},
				{401.1361761624017, 188.52941055867353},
				{419.01617616239474, 133.68941055861768},
				{416.0161761624017, 144.72941055863868},
				{410.7361761624017, 161.76941055854087},
				{409.77617616240167, 165.3694105586177},
				{401.1361761624017, 188.52941055867353},
				{398.2561761624017, 191.4094105585129},
			},
		},
	}

	for tname, tt := range tests {
		actual := [][]float64{}

		var tileWidth = tt.tileWidth
		var tileHeight = tt.tileHeight
		var envelope = tt.envelope

		scale := func(x, y float64) (float64, float64) {
			x = envelope.Px(x) * tileWidth
			y = tileHeight - (envelope.Py(y) * tileHeight)
			return x, y
		}

		g, err := gctx.NewGeomFromWKT(tt.wkt)
		if err != nil {
			t.Fatal(err)
		}

		cs := GetPoints(g)

		csd := *cs
		n := len(csd)
		x, y := csd[0][0], csd[0][1]
		x, y = scale(x, y)
		actual = append(actual, []float64{
			x,
			y,
		})

		for i := 1; i < n; i++ {
			x, y = scale(csd[i][0], csd[i][1])
			actual = append(actual, []float64{
				x,
				y,
			})
		}

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("%v: Expected [%+v]\nGot [%+v]", tname, tt.expected, actual)
		}
	}
}

func TestGetGeometryCenter(t *testing.T) {
	tests := map[string]struct {
		wkt      string
		expected []float64
	}{
		"Invalid": {
			wkt: "MULTILINESTRING ((380626.0000000000000000 413800.9999997689155862, 380627.0000000000582077 413806.9999997675768100, 380632.0000000000000000 413810.9999997689155862, 380642.0000000000000000 413807.9999997679842636, 380643.0000000000582077 413801.9999997690320015, 380638.0000000000000000 413795.9999997689737938, 380628.0000000000000000 413796.9999997678096406, 380626.0000000000000000 413800.9999997689155862))",
			expected: []float64{
				380634.5,
				413803.4999997689,
			},
		},
		"Spotland Road": {
			wkt: "LINESTRING (388603.3400000000256114 413816.2299999999813735, 388662.0000000000000000 413776.0000000000000000, 388682.0000000000000000 413767.0000000000000000, 388697.9799999999813735 413767.5700000000069849)",
			expected: []float64{
				388650.66000000003,
				413791.615,
			},
		},
		"SD": {
			wkt: "POLYGON ((300000.0000000000000000 400000.0000000000000000, 400000.0000000000000000 400000.0000000000000000, 400000.0000000000000000 500000.0000000000000000, 300000.0000000000000000 500000.0000000000000000, 300000.0000000000000000 400000.0000000000000000))",
			expected: []float64{
				350000.0,
				450000.0,
			},
		},
	}

	for tname, tt := range tests {
		g, err := gctx.NewGeomFromWKT(tt.wkt)
		if err != nil {
			t.Fatal(err)
		}

		actual := CenterFromGeometry(g)

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("%v: Expected [%#v]\nGot [%#v]", tname, tt.expected, actual)
		}
	}
}

func TestBoundary(t *testing.T) {
	tests := map[string]struct {
		wkt      string
		expected string
	}{
		"Invalid": {
			wkt:      "MULTILINESTRING ((380626.0000000000000000 413800.9999997689155862, 380627.0000000000582077 413806.9999997675768100, 380632.0000000000000000 413810.9999997689155862, 380642.0000000000000000 413807.9999997679842636, 380643.0000000000582077 413801.9999997690320015, 380638.0000000000000000 413795.9999997689737938, 380628.0000000000000000 413796.9999997678096406, 380626.0000000000000000 413800.9999997689155862))",
			expected: "[380626.000000 413796.000000 380643.000000 413811.000000]",
		},
	}

	for tname, tt := range tests {
		g, err := gctx.NewGeomFromWKT(tt.wkt)
		if err != nil {
			t.Fatal(err)
		}

		actual := g.Bounds()

		if tt.expected != actual.String() {
			t.Errorf("%v: Expected [%+v]\nGot [%+v]", tname, tt.expected, actual.String())
		}
	}
}

func TestDrawInvalidBoundary(t *testing.T) {
	tests := map[string]struct {
		dim       int
		wkt       string
		boundsWkt string
		expected  string
	}{
		"invalid": {
			dim:       600,
			wkt:       "MULTILINESTRING ((380626.0000000000000000 413800.9999997689155862, 380627.0000000000582077 413806.9999997675768100, 380632.0000000000000000 413810.9999997689155862, 380642.0000000000000000 413807.9999997679842636, 380643.0000000000582077 413801.9999997690320015, 380638.0000000000000000 413795.9999997689737938, 380628.0000000000000000 413796.9999997678096406, 380626.0000000000000000 413800.9999997689155862))",
			boundsWkt: "POLYGON ((380596.0000000000000000 413770.0000000000000000, 380656.0000000000000000 413770.0000000000000000, 380656.0000000000000000 413830.0000000000000000, 380596.0000000000000000 413830.0000000000000000, 380596.0000000000000000 413770.0000000000000000))",
			expected:  "[380626.000000 413796.000000 380643.000000 413811.000000]",
		},
	}

	for tname, tt := range tests {
		b, err := gctx.NewGeomFromWKT(tt.boundsWkt)
		if err != nil {
			t.Fatal(err)
		}

		envelope, err := ToEnvelope(b)
		if err != nil {
			t.Fatal(err)
		}

		tileHeight := float64(tt.dim)
		tileWidth := float64(tt.dim)

		scale := func(x, y float64) (float64, float64) {
			x = envelope.Px(x) * tileWidth
			y = tileHeight - (envelope.Py(y) * tileHeight)
			return x, y
		}

		m := image.NewRGBA(image.Rect(0, 0, tt.dim, tt.dim))
		draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.Point{0, 0}, draw.Src)
		gc := draw2dimg.NewGraphicContext(m)

		gc.SetDPI(72)

		g, err := gctx.NewGeomFromWKT(tt.wkt)
		if err != nil {
			t.Fatal(err)
		}

		actual := g.Bounds()

		strokeColour := black
		fillColour := blue
		err = DrawLine(gc, g, 1, fillColour, 1, strokeColour, scale)
		if err != nil {
			t.Fatal(err)
		}

		if tt.expected != actual.String() {
			t.Errorf("%v: Expected [%+v]\nGot [%+v]", tname, tt.expected, actual)
		}

		err = savePNG(fmt.Sprintf("test-output/boundary_test_%v.png", tname), m)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestUnaryTopo(t *testing.T) {
	black := color.RGBA{0x00, 0x00, 0x00, 0xFF}
	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}

	m := image.NewRGBA(image.Rect(0, 0, 600.0, 600.0))
	draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.Point{0, 0}, draw.Src)
	gc := draw2dimg.NewGraphicContext(m)

	gc.SetDPI(72)

	scale := func(x, y float64) (float64, float64) {
		return 10 * x, 10 * y
	}

	strokeWidth := 0.0
	fillColor := black
	strokeColor := black

	g, err := gctx.NewGeomFromWKT("POLYGON ((10 5, 10 0, 0 0, 0 10, 5 10, 5 15, 15 15, 15 5, 10 5))")
	if err != nil {
		t.Fatal(err)
	}
	err = DrawPolygon(gc, g, fillColor, strokeColor, strokeWidth, scale)
	if err != nil {
		t.Fatal(err)
	}

	// draw the image
	err = savePNG("test-output/unary-topo.png", m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestToLinestring(t *testing.T) {
	tests := map[string]struct {
		wkt      string
		expected string
	}{
		"Invalid": {
			wkt:      "LINESTRING (388603.3400000000256114 413816.2299999999813735, 388662.0000000000000000 413776.0000000000000000, 388682.0000000000000000 413767.0000000000000000, 388697.9799999999813735 413767.5700000000069849)",
			expected: "LINESTRING (388603.3400000000256114 413816.2299999999813735, 388662.0000000000000000 413776.0000000000000000, 388682.0000000000000000 413767.0000000000000000, 388697.9799999999813735 413767.5700000000069849)",
		},
	}

	for tname, tt := range tests {
		g, err := gctx.NewGeomFromWKT(tt.wkt)
		if err != nil {
			t.Fatal(err)
		}

		actual, err := ToLineString(g)
		if err != nil {
			t.Fatal(err)
		}

		if tt.expected != actual.String() {
			t.Errorf("%v: Expected [%+v]\nGot [%+v]", tname, tt.expected, actual.String())
		}
	}
}

func TestToEnvelope(t *testing.T) {
	xmin := 0.0
	xmax := 600.0
	ymin := 0.0
	ymax := 600.0

	bounds, err := BoundsGeom(xmin, xmax, ymin, ymax)
	if err != nil {
		t.Fatal(err)
	}

	envelope, err := ToEnvelope(bounds)
	if err != nil {
		t.Fatal(err)
	}

	expected := Envelope{
		Min: []float64{
			0,
			0,
		},
		Max: []float64{
			600,
			600,
		},
	}

	if !reflect.DeepEqual(expected, envelope) {
		t.Fatalf("expected %#v, got %#v", expected, envelope)
	}
}

func TestScaleLine(t *testing.T) {
	scale := func(x, y float64) (float64, float64) {
		return 10 * x, 10 * y
	}

	tests := map[string]struct {
		wkt      string
		expected string
	}{
		"Invalid": {
			wkt:      "MULTILINESTRING ((388888.9999999999417923 413242.9999997675186023, 388874.0000000000000000 413258.9999997682753019))",
			expected: "MULTILINESTRING ((3888889.9999999995343387 4132429.9999976754188538, 3888740.0000000000000000 4132589.9999976828694344))",
		},
	}

	for tname, tt := range tests {
		g, err := gctx.NewGeomFromWKT(tt.wkt)
		if err != nil {
			t.Fatal(err)
		}

		scaled, err := ScaleLine(g, scale)
		if err != nil {
			t.Fatal(err)
		}

		if tt.expected != scaled.String() {
			t.Errorf("%v: ScaleLine - Expected [%+v]\nGot [%+v]", tname, tt.expected, scaled.String())
		}
	}
}

func savePNG(fname string, m image.Image) error {
	dir, _ := path.Split(fname)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	err = draw2dimg.SaveToPngFile(fname, m)
	if err != nil {
		return err
	}

	return nil
}

func BenchmarkScaleLine(b *testing.B) {
	wkt := "MULTILINESTRING ((388874 413258.9999997683,388844.99999999994 413290.99999976775,388740.99999999994 413427.9999997701,388659.00000000006 413499.99999976833,388648 413512.9999997685,388583.00000000006 413820.9999997686))"
	tileHeight := float64(600)
	tileWidth := float64(600)

	boundsWkt := "LINESTRING (386374.0000000000000000 410959.0000000000000000, 391374.0000000000000000 410959.0000000000000000, 391374.0000000000000000 415959.0000000000000000, 386374.0000000000000000 415959.0000000000000000, 386374.0000000000000000 410959.0000000000000000)"
	bounds, err := gctx.NewGeomFromWKT(boundsWkt)
	if err != nil {
		b.Fatal(err)
	}

	envelope, err := ToEnvelope(bounds)
	if err != nil {
		b.Fatal(err)
	}

	g, err := gctx.NewGeomFromWKT(wkt)
	if err != nil {
		b.Fatal(err)
	}

	scale := func(x, y float64) (float64, float64) {
		x = envelope.Px(x) * tileWidth
		y = tileHeight - (envelope.Py(y) * tileHeight)
		return x, y
	}

	for i := 0; i < b.N; i++ {
		_, err = ScaleLine(g, scale)
		if err != nil {
			b.Fatal(err)
		}
	}
}
