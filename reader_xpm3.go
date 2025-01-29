package xpm

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	"github.com/vgarvardt/x11colors-go"
)

func (d *decoder) parseXPM3Metadata() error {
	// Hackily scan for the variable name. Not necessary, but might as well expose it.
	header := XPM3Header{}
	for d.scanner.Scan() {
		token := d.scanner.Text()
		lastSpace := -1
		for i, r := range token {
			if r == ' ' {
				lastSpace = i
			}
			if r == '[' {
				header.name = token[lastSpace+1 : i]
				break
			}
		}
		if header.name != "" {
			break
		}
	}
	// Read our values line.
	{
		var err error
		d.scanner.Scan()
		token := d.scanner.Text()
		token = strings.TrimPrefix(token, "\"")
		token = strings.TrimSuffix(token, "\",")
		values := strings.Split(token, " ")
		if len(values) < 4 {
			return FormatError("too few values in XPM3 metadata")
		}
		header.width, err = strconv.Atoi(values[0])
		if err != nil {
			return FormatError("invalid width in XPM3 metadata")
		}
		header.height, err = strconv.Atoi(values[1])
		if err != nil {
			return FormatError("invalid height in XPM3 metadata")
		}
		header.nColors, err = strconv.Atoi(values[2])
		if err != nil {
			return FormatError("invalid colors in XPM3 metadata")
		}
		header.cPP, err = strconv.Atoi(values[3])
		if err != nil {
			return FormatError("invalid characters per pixel in XPM3 metadata")
		}
		// Read hotspots as well.
		if len(values) == 6 {
			header.hotspotX, err = strconv.Atoi(values[4])
			if err != nil {
				return FormatError("invalid hotspot x in XPM3 metadata")
			}
			header.hotspotY, err = strconv.Atoi(values[5])
			if err != nil {
				return FormatError("invalid hotspot y in XPM3 metadata")
			}
		}
	}
	// Read our color entries.
	for i := 0; i < header.nColors; i++ {
		xclr := XPM3Color{}
		d.scanner.Scan()
		token := d.scanner.Text()
		token = strings.TrimPrefix(token, "\"")
		token = strings.TrimSuffix(token, "\",")

		swatch := ""
		var kind rune
		clr := ""
		foundSwatch := false
		foundKind := false
		for i, r := range token {
			if !foundSwatch {
				if (r == ' ' && i != 0) || r == '\t' {
					foundSwatch = true
					continue
				}
				swatch += string(r)
			} else if !foundKind {
				if r == ' ' || r == '\t' {
					foundKind = true
					continue
				}
				kind = r
			} else {
				if r == ' ' || r == '\t' || r == '\n' {
					break
				}
				clr += string(r)
			}
		}

		if !foundSwatch || !foundKind {
			return FormatError("invalid XPM3 color entry")
		}
		if clr == "" {
			return FormatError("invalid XPM3 color entry")
		}

		if len(swatch) != header.cPP {
			return FormatError("invalid characters per pixel in XPM3 color entry")
		}
		xclr.Chars = swatch
		switch kind {
		case 'c':
			xclr.Kind = XPM3ColorKindColor
			if clr[0] == '#' {
				// Read RGB hex.
				if len(clr) != 7 {
					return FormatError("invalid hex color in XPM3 color entry")
				}
				r := clr[1:3]
				g := clr[3:5]
				b := clr[5:7]
				red, err := strconv.ParseUint(r, 16, 8)
				if err != nil {
					return FormatError("invalid red value in XPM3 color entry")
				}
				green, err := strconv.ParseUint(g, 16, 8)
				if err != nil {
					return FormatError("invalid green value in XPM3 color entry")
				}
				blue, err := strconv.ParseUint(b, 16, 8)
				if err != nil {
					return FormatError("invalid blue value in XPM3 color entry")
				}
				xclr.Color = color.NRGBA{uint8(red), uint8(green), uint8(blue), 0xff}
			} else if clr == "None" {
				xclr.Color = color.NRGBA{0, 0, 0, 0}
			} else {
				x11color, exists := x11colors.GetByName(clr)
				if !exists {
					return FormatError("invalid X11 color name in XPM3 color entry")
				}
				xclr.Color = color.NRGBA{x11color.RGBA.R, x11color.RGBA.G, x11color.RGBA.B, 0xff}
			}
		case 'm':
			xclr.Kind = XPM3ColorKindMonochrome
			return FormatError("XPM3 monochrome color not implemented")
		case 'g':
			xclr.Kind = XPM3ColorKindGrayscale
			return FormatError("XPM3 grayscale color not implemented")
		case 's':
			xclr.Kind = XPM3ColorKindSymbolic
			return FormatError("XPM3 symbolic color not implemented")
		default:
			return FormatError("invalid color kind in XPM3 color entry")
		}
		header.colors = append(header.colors, xclr)
	}

	d.xpmHeader = header

	return nil
}

func (d *decoder) parseXPM3Pixels() error {
	header := d.xpmHeader.(XPM3Header)

	img := image.NewNRGBA(image.Rect(0, 0, header.width, header.height))

	y := 0
	for d.scanner.Scan() {
		token := d.scanner.Text()

		inDef := false
		x := 0
		for i := 0; i < len(token); i += header.cPP {
			r := token[i]
			if !inDef {
				if r == '"' {
					inDef = true
					continue
				}
			}
			if r == '"' && token[i-1] != '\\' {
				break
			}
			entry := token[i : i+header.cPP]
			color := header.XPM3Color(entry)
			if color == nil {
				return FormatError("invalid color entry in XPM3 pixels")
			}
			img.Set(x, y, color.Color)
			x++
		}
		y++
	}

	d.image = img

	return nil
}
