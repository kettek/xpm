package xpm

import "image/color"

type Color struct {
	Type  ColorType
	Value string
	RGBA  color.RGBA /* Populated on decode or converted to Value on encode */
}

type ColorType uint8

const (
	SymbolicEntry ColorType = iota
	ColorEntry
	MonochromeEntry
	GrayscaleEntry
)
