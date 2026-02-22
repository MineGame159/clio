package ui

import (
	"image"
	"image/draw"
	"unsafe"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/ploMP4/chafa-go"
)

var symbolMap *chafa.SymbolMap

func init() {
	symbolMap = chafa.SymbolMapNew()
	chafa.SymbolMapAddByTags(symbolMap, chafa.CHAFA_SYMBOL_TAG_ALL)
}

type Image struct {
	baseWidget

	source *image.RGBA

	lastWidth  uint
	lastHeight uint

	canvas *chafa.Canvas
}

func (w *Image) SetSource(source image.Image) {
	if source == nil {
		w.source = nil
		return
	}

	switch source := source.(type) {
	case *image.RGBA:
		w.source = source

	default:
		src := image.NewRGBA(source.Bounds())
		draw.Draw(src, src.Bounds(), source, src.Bounds().Min, draw.Src)

		w.source = src
	}
}

func (w *Image) CalcRequiredSize() (int, int) {
	if w.source != nil {
		ratio := float64(w.source.Rect.Dx()) / float64(w.source.Rect.Dy())

		w.requiredWidth = int(5 * ratio)
		w.requiredHeight = 5
	} else {
		w.requiredWidth = 0
		w.requiredHeight = 0
	}

	return w.requiredWidth, w.requiredHeight
}

func (w *Image) HandleEvent(_ any) {
}

func (w *Image) Draw(screen tcell.Screen, rect Rect) {
	if w.source == nil {
		return
	}

	// Calculate image size
	targetRatio := float64(w.source.Rect.Dx()) / float64(w.source.Rect.Dy())

	width := rect.Width
	height := rect.Height

	ratio := float64(width) * 0.45 / float64(height)

	if targetRatio > ratio {
		height = int(float64(width) * 0.45 / targetRatio)
	} else {
		width = int(float64(height) * targetRatio / 0.45)
	}

	// "Characterize" image
	config := chafa.CanvasConfigNew()
	defer chafa.CanvasConfigUnref(config)

	chafa.CanvasConfigSetGeometry(config, int32(width), int32(height))
	chafa.CanvasConfigSetSymbolMap(config, symbolMap)

	canvas := chafa.CanvasNew(config)
	defer chafa.CanvasUnRef(canvas)

	chafa.CanvasDrawAllPixels(
		canvas,
		chafa.CHAFA_PIXEL_RGBA8_UNASSOCIATED,
		w.source.Pix,
		int32(w.source.Rect.Dx()),
		int32(w.source.Rect.Dy()),
		int32(w.source.Stride),
	)

	// Copy canvas into screen
	//goland:noinspection GoRedundantConversion
	cells := unsafe.Slice((*canvasCell)(unsafe.Pointer(canvas.Cells)), width*height)

	for x := range width {
		for y := range height {
			cell := cells[y*width+x]

			fg := color.NewHexColor(int32(cell.FgColor))
			bg := color.NewHexColor(int32(cell.BgColor))
			style := tcell.StyleDefault.Foreground(fg).Background(bg)

			screen.SetContent(rect.X+x, rect.Y+y, cell.C, nil, style)
		}
	}
}

type canvasCell struct {
	C       rune
	FgColor uint32
	BgColor uint32
}
