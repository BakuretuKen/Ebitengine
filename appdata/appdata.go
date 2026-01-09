// Windows用 AppData 読み書きライブラリ
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var appDirName string

// SetAppDirName アプリケーションのディレクトリ名を設定します
// main関数の最初で一度だけ呼び出してください
func SetAppDirName(name string) {
	appDirName = name
}

func getUserSaveDir(appName string) (string, error) {
	// Windows の場合
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		return "", fmt.Errorf("LOCALAPPDATA not defined")
	}
	saveDir := filepath.Join(appData, appName)
	return saveDir, nil
}

func SaveGameData(filename string, data []byte) error {
	dir, err := getUserSaveDir(appDirName)
	if err != nil {
		return err
	}
	// ディレクトリが存在しない場合は作る
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, filename)
	return os.WriteFile(path, data, 0644)
}

func LoadGameData(filename string) ([]byte, error) {
	dir, err := getUserSaveDir(appDirName)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, filename)
	return os.ReadFile(path)
}
