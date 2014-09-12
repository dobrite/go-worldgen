package main

import (
	"image"
	"image/jpeg"
	"math"
	"os"

	"github.com/afolmert/libtcod-go/tcod"
	"github.com/larspensjo/Go-simplex-noise/simplexnoise"
)

const (
	dx         = 1024
	dy         = 1024
	scale      = 10.0
	octaves    = 128
	lacunarity = 2.0
)

var noise *tcod.Noise
var rnd *tcod.Random

func main() {
	rnd = tcod.NewRandomFromSeed(0xDEADBEEF)
	noise = tcod.NewNoise(2, rnd)
	pixels := Pic(dx, dy)
	img := Create(pixels)
	WriteImage("img.jpg", img)
}

func noise1(x float64) float64 {
	return simplexnoise.Noise1(x)
}

func noise2(x, y float64) float64 {
	return simplexnoise.Noise2(x, y)
}

func getTCODFBM(x, y float64, octaves int) float64 {
	arr := []float32{float32(x), float32(y)}
	n := noise.GetFbm(tcod.FloatArray(arr), float32(octaves))
	return float64(n)
}

func noise3(x, y, z float64) float64 {
	return simplexnoise.Noise3(x, y, z)
}

func getSimplexFBM(x, y float64, octaves int, lacunarity float64) float64 {
	result := 0.0
	f := 1.0
	var exponents [129]float64

	for i := 0; i <= octaves; i++ {
		exponents[i] = 1.0 / f
		f *= lacunarity
		result += noise2(x, y) * exponents[i]
		x *= lacunarity
		y *= lacunarity
	}

	// clamp -1.0 to 1.0
	ret := math.Max(-1, math.Min(1, result))
	return ret
}

func Pic(dx, dy int) [][]uint8 {
	pixels := make([][]uint8, dy)
	for y := 0; y < dy; y++ {
		pixels[y] = make([]uint8, dx)
		for x := 0; x < dx; x++ {
			val := getSimplexFBM(float64(x)/scale, float64(y)/scale, octaves, lacunarity)
			//val := getTCODFBM(float64(x)/scale, float64(y)/scale, octaves)
			val = (1 + val) / 2
			pixels[y][x] = uint8(val * 255)
		}
	}
	return pixels
}

func Create(data [][]uint8) image.Image {
	m := image.NewRGBA(image.Rect(0, 0, dx, dy))
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			v := data[y][x]
			i := y*m.Stride + x*4
			m.Pix[i] = v
			m.Pix[i+1] = v
			m.Pix[i+2] = v
			m.Pix[i+3] = 255
		}
	}
	return m
}

func WriteImage(n string, m image.Image) {
	f, err := os.OpenFile(n, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(f, m, nil)
	if err != nil {
		panic(err)
	}
}
