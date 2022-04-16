package txtplacer_test

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"

    "github.com/keithroger/txtplacer"
)

var (
	imgdir  = filepath.Join("testdata", "images")
	fontdir = filepath.Join("testdata", "font")
)

func TestNewTextPlacer(t *testing.T) {
	tt := []struct {
		dst      draw.Image
		fontpath string
		size     float64
		want     txtplacer.Placer
	}{
		{
			image.NewGray(image.Rect(0, 0, 200, 200)),
			"fakepath.ttf",
			24.0,
			txtplacer.Placer{},
		},
	}

	for _, tc := range tt {
		ans, _ := txtplacer.NewPlacer(tc.dst, tc.fontpath, tc.size)
		if ans != tc.want {
			t.Errorf("%v != %v", ans, tc.want)
		}
	}
}

func TestWriteAt(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))

	placer, err := txtplacer.NewPlacer(img, filepath.Join(fontdir, "Lato.ttf"), 12.0)
	if err != nil {
		t.Error(err)
	}

	tt := []struct {
		pt        image.Point
		text      string
		wrapwidth int
		align     string
	}{
        {image.Point{200, 200}, "WriteAt\npos: 200, 200\nalign: left\nwrap: 200\nA multiline\nexample", 200, "left"},
        {image.Point{10, 200}, "WriteAt\npos: 10, 200\nalign: center\nwrap: 200\nA long word wraping example. A long long word wrapping example.", 200, "center"},
        {image.Point{10, 200}, "Write at\npos: 10, 200\nalign: right\nA\nright\naligned\nmultiline\nexample", 200, "right"},
	}

	for i, tc := range tt {
		draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)
		placer.WriteAt(tc.pt, tc.text, tc.wrapwidth, tc.align)

		outfile := filepath.Join(imgdir, "WriteAt"+strconv.Itoa(i)+".png")
		outputTestImg(outfile, img, t)
	}
}

func TestCenterAt(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 400, 400))

	placer, err := txtplacer.NewPlacer(img, filepath.Join(fontdir, "Lato.ttf"), 48.0)
	if err != nil {
		t.Error(err)
	}

	center := image.Point{200, 200}
	placer.CenterAt(center, "func: CenterAt\npt: 200, 200\nwrapwidth: 200\nalign: \"left\"", 200, "left")

	outfile := filepath.Join(imgdir, "CenterAt.png")
	outputTestImg(outfile, img, t)
}

func TestWriteAtCenter(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))

	placer, err := txtplacer.NewPlacer(img, filepath.Join(fontdir, "Lato.ttf"), 48.0)
	if err != nil {
		t.Error(err)
	}

	placer.WriteAtCenter("WriteAtCenter\nwrapwidth: 100\nalign: center", 100, "center")

	outfile := filepath.Join(imgdir, "WriteAtCenter.png")
	outputTestImg(outfile, img, t)
}

func outputTestImg(filepath string, img draw.Image, t *testing.T) {
	outfile, err := os.Create(filepath)
	if err != nil {
		t.Error(err)
	}

	err = png.Encode(outfile, img)
	if err != nil {
		t.Error(err)
	}
}
