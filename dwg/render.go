package dwg

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
)

// RenderOptions controls rasterization of a Document to a raster image.
type RenderOptions struct {
	Width, Height int
	Margin        int         // pixels of padding around the fitted drawing
	Background    color.Color // defaults to white
	Line          color.Color // defaults to black
}

func (o RenderOptions) withDefaults() RenderOptions {
	if o.Width <= 0 {
		o.Width = 1600
	}
	if o.Height <= 0 {
		o.Height = 1200
	}
	if o.Margin <= 0 {
		o.Margin = 20
	}
	if o.Background == nil {
		o.Background = color.White
	}
	if o.Line == nil {
		o.Line = color.Black
	}
	return o
}

// Render rasterizes the document's entities into an RGBA image, fitting the
// drawing's bounding box to the requested canvas (preserving aspect ratio,
// flipping Y since DWG space is Y-up and image space is Y-down).
func Render(doc *Document, opts RenderOptions) *image.RGBA {
	opts = opts.withDefaults()
	img := image.NewRGBA(image.Rect(0, 0, opts.Width, opts.Height))
	fillBackground(img, opts.Background)

	min, max, ok := doc.BoundingBox()
	if !ok {
		return img
	}
	drawW := float64(opts.Width - 2*opts.Margin)
	drawH := float64(opts.Height - 2*opts.Margin)
	spanX := max.X - min.X
	spanY := max.Y - min.Y
	if spanX <= 0 {
		spanX = 1
	}
	if spanY <= 0 {
		spanY = 1
	}
	scale := math.Min(drawW/spanX, drawH/spanY)
	if scale <= 0 || math.IsInf(scale, 0) || math.IsNaN(scale) {
		scale = 1
	}
	offsetX := float64(opts.Margin) + (drawW-spanX*scale)/2
	offsetY := float64(opts.Margin) + (drawH-spanY*scale)/2

	project := func(p Point3) (int, int) {
		px := offsetX + (p.X-min.X)*scale
		py := float64(opts.Height) - (offsetY + (p.Y-min.Y)*scale)
		return int(math.Round(px)), int(math.Round(py))
	}

	for _, e := range doc.Entities {
		for _, poly := range e.Polylines() {
			for i := 0; i+1 < len(poly); i++ {
				x0, y0 := project(poly[i])
				x1, y1 := project(poly[i+1])
				drawLine(img, x0, y0, x1, y1, opts.Line)
			}
		}
	}
	return img
}

// EncodePNG writes img as a PNG to w.
func EncodePNG(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}

func fillBackground(img *image.RGBA, c color.Color) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}

// drawLine rasterizes a line segment using Bresenham's algorithm.
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	err := dx + dy
	b := img.Bounds()
	for {
		if x0 >= b.Min.X && x0 < b.Max.X && y0 >= b.Min.Y && y0 < b.Max.Y {
			img.Set(x0, y0, c)
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
