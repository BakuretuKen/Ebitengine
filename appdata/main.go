package main

import (
	"fmt"
	"log"
)

// main 読み書きテスト
func main() {
	SetAppDirName("BakuretuKenGame")

	err := SaveGameData("save.json", []byte("Hello Sofmap World"))
	if err != nil {
		log.Fatal(err)
	}

	data, err := LoadGameData("save.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
