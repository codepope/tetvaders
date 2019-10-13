package tetvaders

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

// Tetvaders is a handle to the game
type Tetvaders struct {
	world           *World
	smallArcadeFont font.Face
	smallFontSize   int
	screenWidth     int
	screenHeight    int
	dpi             int
}

var currgame *Tetvaders

//Init initialises a Tetvaders game and makes it current
func (t *Tetvaders) Init() {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}

	t.smallFontSize = 8
	t.screenWidth = 128
	t.screenHeight = 128
	t.dpi = 72

	t.smallArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(t.smallFontSize),
		DPI:     float64(t.dpi),
		Hinting: font.HintingFull,
	})

	t.world = NewWorld(t.screenWidth, t.screenHeight)
	currgame = t
}

// Run runs the currgame
func (t *Tetvaders) Run() {
	if err := ebiten.Run(update, t.screenWidth, t.screenHeight, 4, "Tetvaders"); err != nil {
		log.Fatal(err)
	}
}

// World represents the game state.
type World struct {
	width     int
	height    int
	base      Base
	dropbear  int
	dropstart int
	dropping  []Shape
	particles []*Particle
	bullets   []Bullet
}

// Shape represents a shape
type Shape struct {
	position PVector
	tettype  int
	bmap     [][]bool
	offxpos  float64
	offypos  float64
	height   float64
	width    float64
	destroy  bool // Remove ASAP
}

func (s *Shape) initial(w *World) {
	s.tettype = rand.Intn(2)
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
	xpos := rand.Float64()*(float64(w.width)-s.width) + s.offxpos
	ypos := float64(0) + s.height

	s.position = NewPVector2D(xpos, ypos)

}

func (s *Shape) move(w *World) {
	if int(s.position.Y-s.offypos+s.height) < w.height {
		s.position.Add(NewPVector2D(0, 1))
		return
	}

	for i := 0; i < int(s.height); i = i + 1 {
		for j := 0; j < int(s.width); j = j + 1 {
			if s.bmap[i][j] {
				w.particles = append(w.particles,
					&Particle{position: NewPVector2D(float64(s.position.X)-float64(s.offxpos)+float64(i), float64(s.position.Y)-float64(s.offypos)+float64(j)),
						gravity: true, direction: NewRandom2dPVector(), velocity: 2.0 * rand.Float64()})
			}
		}
	}

	s.destroy = true
}

var (
	smallArcadeFont font.Face
)

// Particle is a debris
type Particle struct {
	position  PVector
	gravity   bool // Under gravity or propelled?
	direction PVector
	velocity  float64
	deleted   bool
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
	w.dropping = make([]Shape, 0, 20)
	w.bullets = make([]Bullet, 0, 20)
	w.particles = make([]*Particle, 0, 20)
}

// Update game state by one tick.
func (w *World) logicupdate() {
	// Add a shape
	w.dropbear = w.dropbear - 1
	if w.dropbear < 0 {
		newshape := &Shape{}
		newshape.initial(w)
		w.dropping = append(w.dropping, *newshape)
		w.dropbear = w.dropstart
	}

	// Move all the shapes

	for i := range w.dropping {
		w.dropping[i].move(w)
	}

	newdropping := make([]Shape, 20)
	for _, s := range w.dropping {
		if !s.destroy {
			newdropping = append(newdropping, s)
		}
	}
	w.dropping = newdropping

	for i, b := range w.bullets {
		if b.ypos > 0 {
			w.bullets[i].ypos = b.ypos - 1
		} else {
			// Through the roof
		}
	}

	for i := range w.particles {
		w.particles[i].position.Add(w.particles[i].direction)

		if w.particles[i].position.X < 0 || w.particles[i].position.X >= float64(w.width) {
			w.particles[i].direction.Mult(-1)
			w.particles[i].position.Add(w.particles[i].direction)
		} else if w.particles[i].position.Y < 0 || w.particles[i].position.Y >= float64(w.height) {
			w.particles[i].deleted = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if w.base.xpos > 10 {
			w.base.xpos = w.base.xpos - 1.
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
	//fmt.Printf("%+v\n", *s)
	for iy, c := range s.bmap {
		for ix, p := range c {
			if p {
				screen.Set(int(s.position.X-s.offxpos+float64(ix)), int(s.position.Y-s.offypos+float64(iy)), color.White)
			}
		}
	}
}

func (b *Bullet) draw(screen *ebiten.Image) {
	screen.Set(b.xpos, b.ypos, color.White)
}

func (p *Particle) draw(screen *ebiten.Image) {
	if p.deleted {
		return
	}
	screen.Set(int(p.position.X), int(p.position.Y), color.White)
}

func (b *Base) draw(screen *ebiten.Image) {
	screen.Set(b.xpos-1, b.ypos, color.White)
	screen.Set(b.xpos, b.ypos, color.White)
	screen.Set(b.xpos+1, b.ypos, color.White)
	screen.Set(b.xpos, b.ypos-1, color.White)
}

// Draw the world
func (t *Tetvaders) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {

	for _, d := range t.world.dropping {
		d.draw(screen)
	}
	for _, b := range t.world.bullets {
		b.draw(screen)
	}

	for _, p := range t.world.particles {
		p.draw(screen)
	}

	t.world.base.draw(screen)
	load := fmt.Sprintf("Load:%d", t.world.base.load)
	text.Draw(screen, load, t.smallArcadeFont, t.world.width-(len(load)*t.smallFontSize), t.smallFontSize, color.White)
}

func update(screen *ebiten.Image) error {
	currgame.world.logicupdate()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-1, -1)

	currgame.Draw(screen, op)
	return nil
}
