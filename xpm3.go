package xpm

import "image/color"

// XPM3ColorKind represents the kind of a color entry.
type XPM3ColorKind int

// Our c, m, g, and s values.
const (
	XPM3ColorKindColor XPM3ColorKind = iota
	XPM3ColorKindMonochrome
	XPM3ColorKindGrayscale
	XPM3ColorKindSymbolic
)

// XPM3Color is used to do something...
type XPM3Color struct {
	Chars string
	Kind  XPM3ColorKind
	Color color.Color // Resolved during parsing.
}

type XPM3Header struct {
	name          string
	width, height int
	nColors       int
	colors        []XPM3Color
	cPP           int
	hotspotX      int
	hotspotY      int
}

func (x *XPM3Header) XPM3Color(key string) *XPM3Color {
	for i := range x.colors {
		if x.colors[i].Chars == key {
			return &x.colors[i]
		}
	}
	return nil
}

func (x XPM3Header) Name() string {
	return x.name
}

func (x XPM3Header) Width() int {
	return x.width
}

func (x XPM3Header) Height() int {
	return x.height
}

func (x XPM3Header) ColorCount() int {
	return x.nColors
}

func (x XPM3Header) CharsPerPixel() int {
	return x.cPP
}

func (x XPM3Header) HotspotX() int {
	return x.hotspotX
}

func (x XPM3Header) HotspotY() int {
	return x.hotspotY
}

func (x XPM3Header) ColorModel() color.Model {
	return color.NRGBAModel
}
