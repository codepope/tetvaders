package main

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// World represents the game state.
type World struct {
	area     []bool
	width    int
	height   int
	shipxpos int
	shipload int
	dropping []Shape
}

type Shape struct {
	xpos    int
	ypos    int
	tettype int
}

func (s *Shape) draw(area []bool, width int, heigh int) {
	area[s.ypos*width+s.xpos] = true
}

type Particle struct {
	xpos int
	ypos int
}

// NewWorld clears a world
func NewWorld(width, height int) *World {
	w := &World{
		area:     make([]bool, width*height),
		width:    width,
		height:   height,
		shipxpos: width / 2,
	}

	return w
}

func (w *World) init() {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			w.area[y*w.width+x] = false
		}
	}
	w.dropping = make([]Shape, 20)
}

// Update game state by one tick.
func (w *World) Update() {
	width := w.width
	height := w.height

	// This is the screen we will draw on
	next := make([]bool, width*height)

	// Add a shape
	w.dropping = append(w.dropping, Shape{ypos: 0, xpos: rand.Intn(width)})

	// Move all the shapes

	for i, d := range w.dropping {
		if d.ypos < height-1 {
			w.dropping[i].ypos = d.ypos + 1
		} else {
			// On the floor
		}
	}

	// Draw all the shapes

	for _, d := range w.dropping {
		d.draw(next, width, height)
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if w.shipxpos > 10 {
			w.shipxpos = w.shipxpos - 1
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if w.shipxpos < (w.width - 10) {
			w.shipxpos = w.shipxpos + 1
		}
	}
	next[(height-1)*width+w.shipxpos] = true
	next[(height-1)*width+(w.shipxpos-1)] = true
	next[(height-1)*width+(w.shipxpos+1)] = true
	next[(height-2)*width+w.shipxpos] = true

	w.area = next
}

func pixelDrop(w *World, x, y int) bool {
	if w.area[y*w.width+x] {
		return true
	}
	return false
}

func (w *World) Draw(pix []byte) {
	for i, v := range w.area {
		if v {
			pix[4*i] = 0xff
			pix[4*i+1] = 0xff
			pix[4*i+2] = 0xff
			pix[4*i+3] = 0xff
		} else {
			pix[4*i] = 0
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0
		}
	}
}

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	world  = NewWorld(screenWidth, screenHeight)
	pixels = make([]byte, screenWidth*screenHeight*4)
)

func update(screen *ebiten.Image) error {
	world.Update()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	world.Draw(pixels)
	screen.ReplacePixels(pixels)
	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Faller Experiment"); err != nil {
		log.Fatal(err)
	}
}
