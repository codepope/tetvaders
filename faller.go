package main

// TODO: Sweep for ammo
// TODO: Explode shapes on bottom
// TODO: Better shape templats
// TODO: Fast bullet/shape/particle collision tests

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

// World represents the game state.
type World struct {
	width     int
	height    int
	base      Base
	dropbear  int
	dropstart int
	dropping  []Shape
	bullets   []Bullet
}

// Shape represents a shape
type Shape struct {
	xpos    int
	ypos    int
	tettype int
	bmap    [][]bool
	offxpos int
	offypos int
	height  int
	width   int
}

func (s *Shape) initial(ypos int, xpos int) {
	s.tettype = rand.Intn(2)
	s.xpos = xpos
	s.ypos = ypos
	switch s.tettype {
	case 0:
		s.bmap = make([][]bool, 3)
		s.bmap[0] = []bool{false, true, false, true, false}
		s.bmap[1] = []bool{false, true, true, true, false}
		s.bmap[2] = []bool{false, false, true, false, false}
		s.offxpos = 2
		s.offypos = 1
		s.height = 3
		s.width = 5
	case 1:
		s.bmap = make([][]bool, 3)
		s.bmap[0] = []bool{true, true, true, true, true}
		s.bmap[1] = []bool{true, false, false, false, false}
		s.bmap[2] = []bool{true, false, false, false, false}
		s.offxpos = 2
		s.offypos = 1
		s.height = 3
		s.width = 5
	}
}

func (s *Shape) move(w *World) {
	if (s.ypos - s.offypos + s.height) < w.height {
		s.ypos = s.ypos + 1
	}
}

var (
	smallArcadeFont font.Face
)

// Particle is a debris
type Particle struct {
	xpos int
	ypos int
}

// Bullet is bullet :)
type Bullet struct {
	xpos int
	ypos int
}

// Base represents our base
type Base struct {
	xpos    int
	ypos    int
	load    int
	maxload int
}

const smallFontSize = 8

// NewWorld clears a world
func NewWorld(width, height int) *World {
	w := &World{
		width:     width,
		height:    height,
		dropbear:  100,
		dropstart: 100,
		base:      Base{xpos: width / 2, ypos: height - 1, load: 20 / 2, maxload: 20},
	}

	return w
}

func (w *World) init() {
	w.dropping = make([]Shape, 20)
	w.bullets = make([]Bullet, 20)
}

// Update game state by one tick.
func (w *World) logicupdate() {
	width := w.width

	// Add a shape
	w.dropbear = w.dropbear - 1
	if w.dropbear < 0 {
		newshape := &Shape{}
		newshape.initial(0, rand.Intn(width))
		w.dropping = append(w.dropping, *newshape)
		w.dropbear = w.dropstart
	}

	// Move all the shapes

	for i := range w.dropping {
		w.dropping[i].move(w)
	}

	for i, b := range w.bullets {
		if b.ypos > 0 {
			w.bullets[i].ypos = b.ypos - 1
		} else {
			// Through the roof
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if w.base.xpos > 10 {
			w.base.xpos = w.base.xpos - 1
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if w.base.xpos < (w.width - 10) {
			w.base.xpos = w.base.xpos + 1
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyZ) {
		if w.base.load > 0 {
			w.bullets = append(w.bullets, Bullet{ypos: w.base.ypos - 1, xpos: w.base.xpos})
			w.base.load = w.base.load - 1
		}
	}

}

func (s *Shape) draw(screen *ebiten.Image) {
	for iy, c := range s.bmap {
		for ix, p := range c {
			if p {
				screen.Set(s.xpos-s.offxpos+ix, s.ypos-s.offypos+iy, color.White)
			}
		}
	}
}

func (b *Bullet) draw(screen *ebiten.Image) {
	screen.Set(b.xpos, b.ypos, color.White)
}

func (b *Base) draw(screen *ebiten.Image) {
	screen.Set(b.xpos-1, b.ypos, color.White)
	screen.Set(b.xpos, b.ypos, color.White)
	screen.Set(b.xpos+1, b.ypos, color.White)
	screen.Set(b.xpos, b.ypos-1, color.White)
}

// Draw the world
func (w *World) Draw(screen *ebiten.Image) {
	for _, d := range w.dropping {
		d.draw(screen)
	}
	for _, b := range w.bullets {
		b.draw(screen)
	}
	w.base.draw(screen)
	load := fmt.Sprintf("Load:%d", w.base.load)
	text.Draw(screen, load, smallArcadeFont, w.width-(len(load)*smallFontSize), smallFontSize, color.White)
}

const (
	screenWidth  = 240
	screenHeight = 240
)

var (
	world = NewWorld(screenWidth, screenHeight)
)

func update(screen *ebiten.Image) error {
	world.logicupdate()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	world.Draw(screen)
	return nil
}

func main() {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72

	smallArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Faller Experiment"); err != nil {
		log.Fatal(err)
	}
}
