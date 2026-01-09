// Ebitengine Sprite 呼び出しテスト
// @see https://ebitengine.org/ja/examples/sprites.html

package main

import (
	_ "image/png"
	"log"
	"math/rand/v2"

	"github.com/ebitengine/debugui"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	maxAngle     = 256
)

var spriteSheet *ebiten.Image

func init() {
	var err error
	spriteSheet, err = LoadEbitenImage("assets/icon.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Sprites struct {
	sprites []*Sprite
	num     int
}

func (s *Sprites) Update() {
	for i := 0; i < s.num; i++ {
		s.sprites[i].Update()
	}
}

const (
	MinSprites = 0
	MaxSprites = 500
)

type Game struct {
	debugui debugui.DebugUI

	sprites       Sprites
	inited        bool
	prevMouseLeft bool
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	g.sprites.sprites = make([]*Sprite, MaxSprites)
	g.sprites.num = 40
	fw, fh := FrameSize()
	if fw == 0 {
		// 初回呼び出しでフレーム分割
		NewSprite(spriteSheet, 3)
		fw, fh = FrameSize()
	}
	for i := range g.sprites.sprites {
		s := NewSprite(spriteSheet, 3)
		s.x = rand.IntN(screenWidth - fw)
		s.y = rand.IntN(screenHeight - fh)
		s.vx = 2*rand.IntN(2) - 1
		s.vy = 2*rand.IntN(2) - 1
		s.angle = rand.IntN(maxAngle)
		s.vAngle = rand.IntN(5) - 2
		s.hitArea = 80
		g.sprites.sprites[i] = s
	}

	// 最初のスプライトだけキャラを2にする
	g.sprites.sprites[0].SetFrame(2)
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}
	// スプライト同士の当たり判定
	for j := 1; j < g.sprites.num; j++ {
		if g.sprites.sprites[0].IsHitWith(g.sprites.sprites[j]) {
			log.Printf("No.0 Sprite hit No.%d Sprite", j)
			g.sprites.sprites[j].SetFrame(1)
		}
	}

	// マウスクリック検出
	currentMouseLeft := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if currentMouseLeft && !g.prevMouseLeft {
		// クリックされた瞬間
		x, y := ebiten.CursorPosition()
		for i := g.sprites.num; i >= 0; i-- {
			// クリック済みのスプライトはスキップ
			if g.sprites.sprites[i].CurrentFrame() == 1 {
				continue
			}
			if g.sprites.sprites[i].Contains(x, y) {
				log.Printf("Clicked sprite %d at position (%d, %d)", i, x, y)
				// フレームを0→1に切り替え
				g.sprites.sprites[i].SetFrame(1)
			}
		}
	}
	g.prevMouseLeft = currentMouseLeft

	g.sprites.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < g.sprites.num; i++ {
		g.sprites.sprites[i].Draw(screen)
	}

	g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sprite Contains Test(Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
