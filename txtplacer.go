package txtplacer

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func NewPlacer(dst draw.Image, fontPath string, size float64) (Placer, error) {
	b, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return Placer{}, fmt.Errorf("Failed to load font: %v", err)
	}

	f, err := opentype.Parse(b)
	if err != nil {
		return Placer{}, fmt.Errorf("Failed to parse font: %v", err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return Placer{}, fmt.Errorf("Failed to create font face: %v", err)
	}

	drawer := font.Drawer{
		Dst:  dst,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(0, 0),
	}

	return Placer{
		Face:   face,
		Drawer: drawer,
		Dst:    dst,
		Size:   size,
	}, nil
}

type Placer struct {
	Face   font.Face
	Drawer font.Drawer
	Dst    draw.Image
	Size   float64
}

// Writes string starting from pt. If there are multiple lines text
// will be aligned within the boundary.
func (p *Placer) WriteAt(pt image.Point, text string, wrapwidth int, align string) {
	lines := p.wordwrap(text, wrapwidth)
	longest := p.maxLineLen(lines)
	textHeight := p.Face.Metrics().Ascent.Ceil() + p.Face.Metrics().Descent.Ceil()

	switch align {
	case "left":
		for i, line := range lines {
			p.Drawer.Dot = fixed.P(pt.X, pt.Y+i*textHeight)
			p.Drawer.DrawString(line)
		}
	case "center":
		for i, line := range lines {
			p.Drawer.Dot = fixed.P(pt.X+(longest-p.pixWidth(line))/2, pt.Y+i*textHeight)
			p.Drawer.DrawString(line)
		}
	case "right":
		for i, line := range lines {
			p.Drawer.Dot = fixed.P(pt.X+(longest-p.pixWidth(line)), pt.Y+i*textHeight)
			p.Drawer.DrawString(line)
		}
	default:
	}
}

func (p *Placer) WriteAtCenter(text string, wrapwidth int, align string) {
	center := image.Point{p.Dst.Bounds().Dx() / 2, p.Dst.Bounds().Dy() / 2}
	p.CenterAt(center, text, wrapwidth, align)
}

// Write string where pt will be in the center after drawn to image.
func (p *Placer) CenterAt(pt image.Point, text string, wrapwidth int, align string) {
	width, height := p.Bounds(text, wrapwidth)
	descent := p.Face.Metrics().Descent.Ceil()
	height -= descent

	p.WriteAt(image.Point{pt.X - width/2, pt.Y - height/2}, text, wrapwidth, align)
}

func (p *Placer) SetFont(fontPath string) error {
	b, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return fmt.Errorf("Failed to load font: %v", err)
	}

	f, err := opentype.Parse(b)
	if err != nil {
		return fmt.Errorf("Failed to parse font: %v", err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    p.Size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return fmt.Errorf("Failed to create font face: %v", err)
	}

	drawer := font.Drawer{
		Dst:  p.Dst,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(0, 0),
	}

	p.Face = face
	p.Drawer = drawer

	return nil
}

// Sets text to be a new color otherwise defaults to image.Black.
func (p *Placer) SetColor(c color.Color) {
	p.Drawer.Src = &image.Uniform{c}
}

// Returns the size the string will take up on the image.
// Font may be drawn outside of the box depending on the font.
func (p *Placer) Bounds(text string, wrapwidth int) (width int, height int) {
	lines := p.wordwrap(text, wrapwidth)
	textHeight := p.Face.Metrics().Ascent.Ceil() + p.Face.Metrics().Descent.Ceil()

	width = p.maxLineLen(lines)
	height = len(lines) * textHeight

	return width, height
}

func (p *Placer) pixWidth(text string) int {
	return font.MeasureString(p.Face, text).Ceil()
}

func (p *Placer) maxLineLen(lines []string) int {
	max := 0
	for _, l := range lines {
		if p.pixWidth(l) > max {
			max = p.pixWidth(l)
		}
	}

	return max
}

func (p *Placer) wordwrap(text string, rowWidth int) []string {
	iter := spaceIter{s: &text, currIdx: 0}
	for iter.nextSpace() != -1 || iter.currIdx != -1 {

		if iter.nextSpace() != -1 && p.pixWidth(text[iter.prevNewline():iter.nextSpace()]) > rowWidth {
			text = text[:iter.currIdx] + strings.Replace(text[iter.currIdx:], " ", "\n", 1)
		} else {
			iter.currIdx = iter.nextSpace()
		}

		iter.currIdx = iter.nextSpace()
	}

	return strings.Split(text, "\n")
}

type spaceIter struct {
	s       *string
	currIdx int
}

func (it *spaceIter) nextSpace() int {
	return it.nextIndex(" ")
}

func (it *spaceIter) nextIndex(sep string) int {
	if it.currIdx == -1 {
		return -1
	}

	if it.currIdx == len(*it.s) {
		return -1
	}

	idx := strings.Index((*it.s)[it.currIdx+1:], sep)
	if idx > -1 {
		idx += it.currIdx + 1
	}

	return idx
}

func (it *spaceIter) prevNewline() int {
	for i := it.currIdx; i > 0; i-- {
		if (*it.s)[i] == '\n' {
			return i
		}
	}

	return 0
}
