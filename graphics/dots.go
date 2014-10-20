package graphics

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"image/color"
	"math"
)

// DotPlot plots a XY plot with dots in it
type Dots struct {
	X []float64
	Y []float64
	color.Color
}

func NewDots(x, y []float64) (*Dots, error) {
	return &Dots{x, y, color.RGBA{G: 100, B: 200, A: 255}}, nil
}

func (pt *Dots) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)

	da.SetColor(pt.Color)

	for i := range pt.Y {
		// Transform the data x, y coordinate of this bubble
		// to the corresponding drawing coordinate.
		x := trX(pt.X[i])
		y := trY(pt.Y[i])

		// Get the radius of this bubble.  The radius
		// is specified in drawing units (i.e., its size
		// is given as the final size at which it will
		// be drawn) so it does not need to be transformed.
		rad := vg.Length(2)

		// Fill a circle centered at x,y on the draw area.
		var p vg.Path
		p.Move(x+rad, y)
		p.Arc(x, y, rad, 0, 2*math.Pi)
		p.Close()
		da.Fill(p)
	}
}

func (pt *Dots) DataRange() (xmin, xmax, ymin, ymax float64) {

	xmin = math.MaxFloat64
	xmax = -math.MaxFloat64
	ymin = math.MaxFloat64
	ymax = -math.MaxFloat64

	for i := range pt.Y {
		if pt.Y[i] > ymax {
			ymax = pt.Y[i]
		}
		if pt.Y[i] < ymin {
			ymin = pt.Y[i]
		}
		if pt.X[i] > xmax {
			xmax = pt.X[i]
		}
		if pt.X[i] < xmin {
			xmin = pt.X[i]
		}
	}
	return 0, xmax+1, 0, ymax
}

func DotsPaint(x, y []float64, title, xlabel, ylabel, file string) {
	p, err := plot.New()
	if err != nil {
		return
	}

	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	bs, err := NewDots(x, y)

	// bs.Color = color.RGBA{R: 196, B: 128, A: 255}
	p.Add(bs)

	if err := p.Save(5, 5, file); err != nil {
		return
	}
}
