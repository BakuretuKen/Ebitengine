// 薬飲み忘れタイマー
// Windowsスタートアップフォルダにプログラムを登録
package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// タイマー設定値 20:00
var myTimerHour = 20
var myTimerMinute = 0

const (
	screenWidth  = 300
	screenHeight = 260
)

var imageReimu *ebiten.Image

func init() {
	var err error
	imageReimu, err = LoadEbitenImage("assets/reimu.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	inited        bool
	prevMouseLeft bool
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(imageReimu, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func isRunAppDate() bool {
	// 現在の 月日 を文字列で取得
	now := time.Now()
	month := now.Month()
	day := now.Day()
	dateString := fmt.Sprintf("%02d%02d", month, day)

	// time.txt を読み込む
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDir := filepath.Dir(execPath)
	// fmt.Printf("DEBUG: execPath = %s\n", execPath)
	// fmt.Printf("DEBUG: execDir = %s\n", execDir)
	timeText, err := os.ReadFile(filepath.Join(execDir, "last_date.txt"))
	if err != nil {
		timeText = []byte("")
	}
	timeTextString := string(timeText)
	// 最終実行日と今日の日付が一致
	if timeTextString == dateString {
		fmt.Println("Already executed today")
		return false
	}

	return true
}

func writeLastDate() {
	// 現在の 月日 を文字列で取得
	now := time.Now()
	month := now.Month()
	day := now.Day()
	dateString := fmt.Sprintf("%02d%02d", month, day)

	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDir := filepath.Dir(execPath)

	// dateString をプログラムのフォルダに time.txt として保存
	os.WriteFile(filepath.Join(execDir, "last_date.txt"), []byte(dateString), 0644)
}

func isWaitAppTime() bool {
	// 現在の HH:MM を数値で取得
	now := time.Now()
	timeInt := now.Hour()*100 + now.Minute()

	// 設定時刻以降なら即実行
	if timeInt >= myTimerHour*100+myTimerMinute {
		fmt.Println("Execute after " + fmt.Sprintf("%d - %d", timeInt, (myTimerHour*100+myTimerMinute)))
		return false
	}

	return true
}

func main() {
	time.Sleep(10 * time.Second)

	isRun := isRunAppDate()
	if !isRun {
		return
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Reimu Drag Timer")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowFloating(true)

	isWait := isWaitAppTime()
	if !isWait {
		// 画像を表示
		if err := ebiten.RunGame(&Game{}); err != nil {
			log.Fatal(err)
		}
	} else {
		// myTimerHourMinuteになるまで待機
		now := time.Now()
		lastTimeInt := now.Hour()*3600 + now.Minute()*60
		targetTimeInt := myTimerHour*3600 + myTimerMinute*60
		durationSec := targetTimeInt - lastTimeInt
		log.Printf("Waiting until %d (%v remaining)", durationSec, durationSec)
		time.Sleep(time.Duration(durationSec) * time.Second)

		// 最終実行日を保存
		writeLastDate()

		// 画像を表示
		if err := ebiten.RunGame(&Game{}); err != nil {
			log.Fatal(err)
		}
	}
}
