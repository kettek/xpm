package xpm

import "image/color"

type XPMHeader interface {
	Width() int
	Height() int
	ColorModel() color.Model
}

type XPM1Header struct {
	format        int
	width, height int
	nColors       int
	cPP           int
}

func (x XPM1Header) Format() int {
	return x.format
}

func (x XPM1Header) Width() int {
	return x.width
}

func (x XPM1Header) Height() int {
	return x.height
}

func (x XPM1Header) ColorCount() int {
	return x.nColors
}

func (x XPM1Header) CharsPerPixel() int {
	return x.cPP
}

func (x XPM1Header) ColorModel() color.Model {
	return color.RGBAModel
}
