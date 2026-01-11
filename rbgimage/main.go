// Ebitengine RGB Draw Image Test

package main

import (
	_ "image/png"
	"log"

	"github.com/ebitengine/debugui"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 400
)

var sourceImage *ebiten.Image

func init() {
	var err error
	sourceImage, err = LoadEbitenImage("assets/image1.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Images struct {
	rgbImage *RgbImage
}

func (s *Images) Update() {
	s.rgbImage.Update()
}

const (
	MinImages = 0
	MaxImages = 500
)

type Game struct {
	debugui debugui.DebugUI

	rgbImage        *RgbImage
	inited        bool
	prevMouseLeft bool
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	g.rgbImage = NewRgbImage(sourceImage, false)
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	// 画面クリックで描画開始
	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	// 	g.rgbImage.StartDraw(1)
	// }
	g.rgbImage.StartDraw(1)

	g.rgbImage.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.rgbImage.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("RGB Image Test(Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
