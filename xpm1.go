package xpm

type XPMHeader interface {
}

type XPM1Header struct {
	format        uint8
	width, height uint
	nColors       uint
	cPP           uint8
}

func (x *XPM1Header) Format() uint8 {
	return x.format
}

func (x *XPM1Header) Width() uint {
	return x.width
}

func (x *XPM1Header) Height() uint {
	return x.height
}

func (x *XPM1Header) ColorCount() uint {
	return x.nColors
}

func (x *XPM1Header) CharsPerPixel() uint8 {
	return x.cPP
}
