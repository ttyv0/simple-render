package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type (
	Model struct {
		verts []Vec3f
		faces [][]int
	}

	Vec3f struct {
		raw []float64
	}
)

func NewModel(fn string) *Model {
	m := new(Model)
	m.verts = make([]Vec3f, 0)
	m.faces = make([][]int, 0)
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	bf := bufio.NewReader(f)
	var eof bool
	for !eof {
		line, err := bf.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		}
		if strings.HasPrefix(line, "v ") {
			var v Vec3f
			v.raw = make([]float64, 3)
			for i, s := range strings.Fields(line[2:]) {
				v.raw[i], _ = strconv.ParseFloat(s, 64)
			}
			m.verts = append(m.verts, v)
		} else if strings.HasPrefix(line, "f ") {
			var f []int
			for _, s := range strings.Fields(line[2:]) {
				for _, s2 := range strings.Split(s, "/") {
					i, _ := strconv.Atoi(s2)
					f = append(f, i-1)
				}
			}
			m.faces = append(m.faces, f)
		}
	}
	log.Println("Len:", len(m.verts))
	return m
}
