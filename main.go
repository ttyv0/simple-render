package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

type (
	Img struct {
		f *os.File
		i *image.RGBA
	}
)

var (
	modelfile = flag.String("modelfile", "african_head.obj", "obj file")
	output    = flag.String("output", "output.png", "name output file")
	width     = flag.Int("width", 800, "output width")
	height    = flag.Int("heigh", 800, "output height")
)

func main() {
	flag.Parse()
	img := NewImg(*output, *height, *width)
	white := color.NRGBA{255, 255, 255, 255}
	black := color.NRGBA{0, 0, 0, 255}
	//red := color.NRGBA{255, 0, 0, 255}
	img.Fill(black)
	model := NewModel(*modelfile)
	img.DrawModel(model, white)
	img.flip_vertically()
	img.Save()
}

func NewImg(filename string, w, h int) *Img {
	img := new(Img)
	img.f, _ = os.Create(filename)
	img.i = image.NewRGBA(image.Rect(0, 0, w, h))
	return img
}

func (img *Img) Save() error {
	encoder := png.Encoder{png.DefaultCompression}
	err := encoder.Encode(img.f, img.i)
	img.f.Close()
	return err
}

func (img *Img) Fill(c color.Color) {
	rect := img.i.Bounds()
	x := rect.Dx()
	y := rect.Dy()
	for a := 0; a <= x; a++ {
		for b := 0; b <= y; b++ {
			img.i.Set(a, b, c)
		}
	}
}

func (img *Img) DrawModel(m *Model, c color.Color) {
	rect := img.i.Bounds()
	w, h := rect.Dx(), rect.Dy()
	l := len(m.verts)
	for _, face := range m.faces {
		for i := 0; i < 3; i++ {
			if l < face[i] || l < face[(i+1)%3] {
				continue
			}
			var v0, v1 Vec3f
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println("Panic:", i, (i+1)%3, face, face[i], face[(i+1)%3], len(m.verts), m.verts[face[i]])
					}
				}()
				v0 = m.verts[face[i]]
				v1 = m.verts[face[(i+1)%3]]
			}()
			//log.Println(v0, v1)
			x0 := int((v0.X() + 1.) * float64(w/2))
			y0 := int((v0.Y() + 1.) * float64(h/2))
			x1 := int((v1.X() + 1.) * float64(w/2))
			y1 := int((v1.Y() + 1.) * float64(h/2))
			//log.Println(x0, y0, x1, y1)
			img.line(x0, y0, x1, y1, c)
		}
	}
}

func (img *Img) line(x0, y0, x1, y1 int, c color.Color) {
	var steep bool
	if math.Abs(float64(x0-x1)) < math.Abs(float64(y0-y1)) {
		x0, y0, x1, y1 = y0, x0, y1, x1
		steep = true
	}
	if x0 > x1 {
		x0, x1, y0, y1 = x1, x0, y1, y0
	}
	for x := x0; x <= x1; x++ {
		t := float64(x-x0) / float64(x1-x0)
		y := int(float64(y0)*(1.-t) + float64(y1)*t)
		if steep {
			img.i.Set(y, x, c)
		} else {
			img.i.Set(x, y, c)
		}
	}
}

func (img *Img) flip_vertically() {
	w, h := img.i.Rect.Dx(), img.i.Rect.Dy()
	flip_img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			flip_img.Set(x, y, img.i.At(w-x, h-y))
		}
	}
	img.i = flip_img
}
