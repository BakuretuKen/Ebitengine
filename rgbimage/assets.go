package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assetsFS embed.FS

// LoadImage は埋め込みアセットから生の画像を読み込む
func LoadImage(path string) (image.Image, error) {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// LoadEbitenImage は埋め込みアセットから画像を読み込む
func LoadEbitenImage(path string) (*ebiten.Image, error) {
	img, err := LoadImage(path)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}
