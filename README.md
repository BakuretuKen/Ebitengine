# Golang Game Engine Ebitengine

## 当たり判定スプライトライブラリ

当たり判定、クリック判定を実装したスプライトライブラリ

### [sprite](https://github.com/BakuretuKen/Ebitengine/tree/main/sprite)

[![](./img/sprite01.png)](https://github.com/BakuretuKen/Ebitengine/tree/main/sprite)

スプライト作成
```
spriteSheet, err = LoadEbitenImage("assets/icon.png")
sprite := NewSprite(spriteSheet, 2) // 2は画像のフレーム数（画像は横並び）
```

クリック判定
> func (s *Sprite) Contains(x, y int) bool

他のスプライトとの当たり判定
> func (s *Sprite) IsHitWith(other *Sprite) bool

## Windows データ保存ライブラリ

Windows のユーザーデータ保存（AppData）を使ったデータ保存ライブラリ

### [appdata](https://github.com/BakuretuKen/Ebitengine/tree/main/appdata)

データ保存

```go
SetAppDirName("BakuretuKenGame")
err := SaveGameData("save.json", []byte("Hello Sofmap World"))
if err != nil {
	log.Fatal(err)
}
```

データ読み込み

```go
SetAppDirName("BakuretuKenGame")
data, err := LoadGameData("save.json")
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(data))
```

## その他

Ebitengine 公式サイト<br />
https://ebitengine.org/ja/

