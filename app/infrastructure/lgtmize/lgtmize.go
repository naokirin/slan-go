package lgtmize

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	domain "github.com/naokirin/slan-go/app/domain/lgtmize"
)

const maskPath = "img/mask.png"

var _ domain.LGTMize = (*LGTMize)(nil)

// LGTMize for lgtm creator
type LGTMize struct{}

// CreateLGTM create lgtm image
func (l *LGTMize) CreateLGTM(url string, size int, color string) (string, error) {
	path, err := download(url)
	if err != nil {
		return "", err
	}
	path, err = lgtm(path, size, color)
	return path, err
}

func download(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(body)
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("content-type %s is not supported", contentType)
	}
	ext := contentType[6:len(contentType)]
	ext = strings.Replace(ext, "jpeg", "jpg", 1)
	_, filename := path.Split(url)
	r := strings.NewReplacer(
		"?", "-",
		"=", "-",
		":", "-",
	)
	filename = r.Replace(filename)
	filename = filepath.Join("lgtm", filename+"."+ext)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		return "", err
	}
	file.Write(body)
	return filename, nil
}

func lgtm(path string, size int, color string) (string, error) {
	src, err := imaging.Open(path)
	if err != nil {
		return path, err
	}
	resized := resize(src, size)
	lgtmized, err := drawLGTM(resized, size, color)
	if err != nil {
		return path, err
	}
	result, err := save(lgtmized, path)
	if err != nil {
		return path, err
	}
	return result, nil
}

func calculateRect(img image.Image, size int) image.Rectangle {
	s := img.Bounds().Size()

	var min, max image.Point
	switch {
	case s.X == s.Y:
		min, max = image.ZP, image.Pt(size, size)
	case s.X > s.Y:
		min = image.Pt((s.X-size)/2, 0)
		max = image.Pt(min.X+size, min.Y+size)
	case s.X < s.Y:
		min = image.Pt(0, (s.Y-size)/2)
		max = image.Pt(min.X+size, min.Y+size)
	}

	return image.Rect(min.X, min.Y, max.X, max.Y)
}

func resize(img image.Image, size int) image.Image {
	s := img.Bounds().Size()

	var x, y int
	switch {
	case s.X == s.Y:
		x, y = size, size
	case s.X > s.Y:
		ratio := float32(s.Y) / float32(s.X)
		x, y = size, int(float32(size)*ratio)
	case s.X < s.Y:
		ratio := float32(s.X) / float32(s.Y)
		x, y = int(float32(size)*ratio), size
	}

	return imaging.Resize(img, x, y, imaging.Box)
}

func drawLGTM(img image.Image, size int, lgtmColor string) (image.Image, error) {
	rect := img.Bounds()
	s := rect.Size()
	result := imaging.New(s.X, s.Y, color.RGBA{0, 0, 0, 0})
	draw.Draw(result, rect, img, rect.Min, draw.Src)

	lgtm := imaging.New(size, size, colorFromString(lgtmColor))
	mask, err := imaging.Open(maskPath)
	if err != nil {
		return nil, err
	}
	maskSize := rect.Size().X
	maskMin := mask.Bounds().Min
	maskPt := image.Point{maskMin.X, maskMin.Y - (size - maskSize)}
	if rect.Size().X > rect.Size().Y {
		maskSize = rect.Size().Y
		maskPt = image.Point{maskMin.X - int(float32(size-maskSize)/2), maskMin.Y}
	}
	mask = resize(mask, maskSize)
	draw.DrawMask(result, calculateRect(img, size), lgtm, lgtm.Bounds().Min, mask, maskPt, draw.Over)
	return result, nil
}

func save(img image.Image, srcPath string) (string, error) {
	srcExt := filepath.Ext(srcPath)
	resultPath := filepath.Join(filepath.Dir(srcPath), "lgtm"+time.Now().Format(time.RFC3339Nano)+srcExt)
	return resultPath, imaging.Save(img, resultPath)
}

func colorFromString(str string) color.RGBA {
	if strings.ToLower(str) == "white" {
		return color.RGBA{255, 255, 255, 255}
	}
	if strings.ToLower(str) == "black" {
		return color.RGBA{0, 0, 0, 255}
	}
	if strings.ToLower(str) == "blue" {
		return color.RGBA{0, 0, 255, 255}
	}
	if strings.ToLower(str) == "red" {
		return color.RGBA{255, 0, 0, 255}
	}
	if strings.ToLower(str) == "green" {
		return color.RGBA{0, 255, 0, 255}
	}
	if strings.ToLower(str) == "magenta" {
		return color.RGBA{215, 21, 126, 255}
	}
	if strings.ToLower(str) == "cyan" {
		return color.RGBA{0, 163, 219, 255}
	}
	if strings.ToLower(str) == "yellow" {
		return color.RGBA{252, 212, 27, 255}
	}
	return color.RGBA{255, 255, 255, 255}
}
