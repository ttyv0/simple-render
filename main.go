package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type (
	Img struct {
		f *os.File
		i *image.RGBA
	}
)

func main() {
	img := NewImg("output.png", 800, 800)
	red := color.NRGBA{255, 0, 0, 255}
	black := color.NRGBA{0, 0, 0, 255}
	img.Fill(black)
	//img.line(13, 20, 80, 40, red)
	//img.line(20, 13, 40, 80, red)
	model := NewModel("african_head.obj")
	img.DrawModel(model, red)
	img.Save()
}

func NewImg(filename string, w, h int) *Img {
	img := new(Img)
	img.f, _ = os.Create(filename)
	img.i = image.NewRGBA(image.Rect(0, 0, w, h))
	return img
}

func (img *Img) Save() error {
	encoder := png.Encoder{png.NoCompression}
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
	for _, face := range m.faces {
		//fmt.Println(face)
		for i := 0; i < 3; i++ {
			l := len(m.verts)
			if l < face[i] || l < face[(i+1)%3] {
				fmt.Println(face)
				continue
			}
			var v0, v1 Vec3f
			func() {
				defer func() {
					if err := recover(); err != nil {
						fmt.Println("Panic:", i, (i+1)%3, face, face[i], face[(i+1)%3], len(m.verts), m.verts[face[i]])
					}
				}()
				v0 = m.verts[face[i]]
				v1 = m.verts[face[(i+1)%3]]
			}()
			x0 := int(v0.raw[0]+1.) * w / 2
			y0 := int(v0.raw[1]+1.) * h / 2
			x1 := int(v1.raw[0]+1.) * w / 2
			y1 := int(v1.raw[1]+1.) * h / 2
			//fmt.Println(x0, y0, x1, y1)
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
	//13, 20, 80, 40
	for x := x0; x <= x1; x++ {
		t := float64(x-x0) / float64(x1-x0)
		y := int(float64(y0)*(1.-t) + float64(y1)*t)
		if steep {
			//fmt.Println(y, x)
			img.i.Set(y, x, c)
		} else {
			img.i.Set(x, y, c)
		}
	}
}
