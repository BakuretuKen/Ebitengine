// Ebitengine RGB Draw Image (refactored)

package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	stageLineLimit = 84 // visibleLines がここに到達したら次の stage へ
	copyStep       = 84 // ピクセル作画間隔
	blockSize      = 2  // 2x2 ブロックで作画
)

type rgbStage int

const (
	stageB rgbStage = iota   // B only
	stageGB                  // G+B
	stageRGB                 // R+G+B（元画像）
	stageDone
)

type RgbImage struct {
	src          *ebiten.Image
	width, height int

	// 描画先（蓄積バッファ）
	buffer *ebiten.Image

	// チャンネル抽出済み画像（B / GB は事前生成、RGB は src を使う）
	channels [2]*ebiten.Image

	// 状態
	frameCount   int
	visibleLines int
	stage        rgbStage
	started      bool
	done         bool
	isReverse    bool
	waitFrame    int
}

func NewRgbImage(src *ebiten.Image, reverse bool) *RgbImage {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	r := &RgbImage{
		src:    src,
		width:  w,
		height: h,

		buffer: ebiten.NewImage(w, h),

		stage:   stageB,
		started: false,
		done:    false,
		isReverse: reverse,
	}

	// 事前生成：実行中に ColorM で毎回作り直す必要はない
	r.channels[0] = createChannelImage(src, stageB)
	r.channels[1] = createChannelImage(src, stageGB)

	if r.isReverse {
		r.buffer = createChannelImage(src, stageRGB)
	}

	return r
}

func (r *RgbImage) StartDraw(waitFrame int) {
	r.waitFrame = waitFrame
	r.started = true
}

func (r *RgbImage) Reset() {
	r.frameCount = 0
	r.visibleLines = 0
	r.stage = stageB
	r.done = false
	r.started = false
	r.buffer.Clear()
}

func (r *RgbImage) Update() {
	if !r.started || r.done {
		return
	}
	r.frameCount++
	// 作画速度設定
	if (r.frameCount % r.waitFrame) == 0 {
		r.visibleLines++
	}

	// stage 遷移（Draw ではなく Update で行う）
	if r.visibleLines * blockSize >= stageLineLimit {
		r.visibleLines = 0
		r.frameCount = 1 // 次フレームで確実に visibleLines++ が実行されるように奇数にする
		r.stage++
		if r.stage >= stageDone {
			r.done = true
		}
	}
}

func (r *RgbImage) Draw(screen *ebiten.Image) {
	if !r.started {
		return
	}

	// 完了後は buffer をそのまま出すだけ（状態を固定）
	if r.done {
		if !r.isReverse {
			screen.DrawImage(r.buffer, nil)
		}
		return
	}

	var src *ebiten.Image
	if r.isReverse {
		src = r.stageImageReverse(r.stage)
	} else {
		src = r.stageImage(r.stage)
	}
	if src == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Go 元コードの式を保つ（意図が「同じ処理」ならここは変えない）
	limit := r.width * (r.height / 2)

	for i := r.visibleLines * 2; i < limit; i += copyStep {
		x := i % r.width
		y := (i / r.width) * 2

		// 2x2 がはみ出すと SubImage が落ちるのでガード
		if x+blockSize > r.width || y+blockSize > r.height {
			continue
		}

		rect := image.Rect(x, y, x+blockSize, y+blockSize)
		tile := src.SubImage(rect).(*ebiten.Image)

		op.GeoM.Reset()
		op.GeoM.Translate(float64(x), float64(y))
		r.buffer.DrawImage(tile, op)
	}

	// 画面描画はフレーム最後に 1 回で十分
	screen.DrawImage(r.buffer, nil)
}

func (r *RgbImage) stageImage(s rgbStage) *ebiten.Image {
	switch s {
	case stageB:
		return r.channels[0]
	case stageGB:
		return r.channels[1]
	case stageRGB:
		return r.src
	default:
		return nil
	}
}

func (r *RgbImage) stageImageReverse(s rgbStage) *ebiten.Image {
	switch s {
	case stageB:
		return r.channels[1]
	case stageGB:
		return r.channels[0]
	case stageRGB:
		return nil
	default:
		return nil
	}
}


// stageB / stageGB のみ新規生成する（stageRGB は src をそのまま使う）
func createChannelImage(src *ebiten.Image, s rgbStage) *ebiten.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	dst := ebiten.NewImage(w, h)

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Reset()

	switch s {
	case stageB:
		op.ColorM.Scale(0, 0, 1, 1) // B のみ
	case stageGB:
		op.ColorM.Scale(0, 1, 1, 1) // G+B
	case stageRGB:
		op.ColorM.Scale(1, 1, 1, 1) // 変化なし（通常ここは使わない）
	default:
		op.ColorM.Scale(1, 1, 1, 1)
	}

	dst.DrawImage(src, op)
	return dst
}
