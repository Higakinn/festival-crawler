package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// import (
// 	"bytes"
// 	"fmt"
// 	"image"
// 	"image/png"
// 	"log"
// 	"os"

// 	"github.com/golang/freetype/truetype"
// 	"golang.org/x/image/font"
// 	"golang.org/x/image/math/fixed"
// )

// func main() {
// 	err := run("ありがとうございます")
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}
// }

// func run(text string) error {
// 	// TODO: 日本語対応
// 	// TODO: テキストのながさを元に動的にサイズ決定ができたら嬉しい
// 	fontSize := 20.0
// 	imageWidth, imageHeight, textTopMargin := 640, 400, 220
// 	// if len(text) > 12 {
// 	// 	return xerrors.New("text length should less 12")
// 	// }
// 	ttf, err := os.ReadFile("./Koruri-Light.ttf")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ft, err := truetype.Parse(ttf)
// 	// ft, err := truetype.Parse(goregular.TTF)
// 	if err != nil {
// 		return err
// 	}

// 	opt := truetype.Options{
// 		Size:              fontSize,
// 		DPI:               0,
// 		Hinting:           0,
// 		GlyphCacheEntries: 0,
// 		SubPixelsX:        0,
// 		SubPixelsY:        0,
// 	}
// 	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
// 	face := truetype.NewFace(ft, &opt)
// 	dr := &font.Drawer{
// 		Dst:  img,
// 		Src:  image.White,
// 		Face: face,
// 		Dot:  fixed.Point26_6{},
// 	}
// 	dr.Dot.X = (fixed.I(imageWidth) - dr.MeasureString(text)) / 2
// 	dr.Dot.Y = fixed.I(textTopMargin)
// 	dr.DrawString(text)

// 	buf := &bytes.Buffer{}
// 	err = png.Encode(buf, img)
// 	if err != nil {
// 		return err
// 	}

// 	f, err := os.Create("./hoge.png")
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	f.Write(buf.Bytes())
// 	return nil
// }

func main() {
	// 256x128の白い画像を作る。
	img := image.NewRGBA(image.Rect(0, 0, 640, 640))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 描画する範囲を決めておく。
	area := image.Rect(30, 20, 640-30, 640-20)
	draw.Draw(img, area, &image.Uniform{color.Gray16{0xdddf}}, image.Point{}, draw.Src)

	// フォントを読み込んで、image/font.faceを作る。
	ttf, err := os.ReadFile("Koruri-Bold.ttf")
	if err != nil {
		log.Fatal(err)
	}
	font_, err := truetype.Parse(ttf)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(font_, &truetype.Options{
		Size: 22,
	})

	// 描画用の構造体を準備する。
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
	}

	// フォントフェイスから1行の高さを取得する。
	lineHeight := face.Metrics().Height.Ceil()
	a := image.Rect(30, 20, 640-30, (lineHeight+22)*2)
	draw.Draw(img, a, &image.Uniform{color.RGBA{R: 255, G: 165, B: 0, A: 255}}, image.Point{}, draw.Src)

	// 描画する文字列。
	// text := "Hello, World! こんにちは、世界！\nThis is a test."
	text := `

【愛媛県】新居浜太鼓祭り 10月16日 ~ 10月18日
【愛媛県】西条祭り 10月16日 ~ 10月18日
【大阪府】岸和田だんじり祭り
【千葉県】佐原の大祭秋祭り`

	y := lineHeight + 22
	d.Dot = fixed.Point26_6{X: fixed.I(area.Min.X + 100), Y: fixed.I(y + 22)}
	d.DrawString("＜2024年10月開催の祭り一覧＞")
	// 折り返しを考慮しながら1行ずつに分割する。
	runes := []rune(text)
	var lines []string
	start := 0
	for i := 0; i < len(runes); i++ {
		// 改行文字を見つけたら改行する。
		if runes[i] == '\n' {
			lines = append(lines, string(runes[start:i]))
			start = i + 1
			continue
		}

		// ここまでの文字列の横幅を計算する。
		width := d.MeasureString(string(runes[start:i]))

		// 横幅が描画範囲を越えていたら改行する。
		if width > fixed.I(area.Dx()) {
			i--
			lines = append(lines, string(runes[start:i]))
			start = i
		}
	}
	// 最後の1行をlinesに加えておく。
	if start < len(runes) {
		lines = append(lines, string(runes[start:]))
	}

	// 1行ずつ描画する。
	for lineOffset, line := range lines {
		y := area.Min.Y + (lineOffset+1)*lineHeight + 22
		d.Dot = fixed.Point26_6{X: fixed.I(area.Min.X), Y: fixed.I(y)}
		d.DrawString(line)
	}

	// 画像をoutput.jpgとして保存する。
	out, err := os.Create("output.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	if err := jpeg.Encode(out, img, nil); err != nil {
		log.Fatal(err)
	}
}
