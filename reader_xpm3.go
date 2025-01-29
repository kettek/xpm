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
		clr := XPM3Color{}
		d.scanner.Scan()
		token := d.scanner.Text()
		token = strings.TrimPrefix(token, "\"")
		token = strings.TrimSuffix(token, "\",")
		values := strings.Split(token, " ")
		if len(values) < 3 {
			return FormatError("too few values in XPM3 color entry")
		}
		if len(values[0]) != header.cPP {
			return FormatError("invalid characters per pixel in XPM3 color entry")
		}
		clr.Chars = values[0]
		switch values[1] {
		case "c":
			clr.Kind = XPM3ColorKindColor
			if values[2][0] == '#' {
				// Read RGB hex.
				if len(values[2]) != 7 {
					return FormatError("invalid hex color in XPM3 color entry")
				}
				r := values[2][1:3]
				g := values[2][3:5]
				b := values[2][5:7]
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
				clr.Color = color.NRGBA{uint8(red), uint8(green), uint8(blue), 0xff}
			} else if values[2] == "None" {
				clr.Color = color.NRGBA{0, 0, 0, 0}
			} else {
				x11color, exists := x11colors.GetByName(values[2])
				if !exists {
					return FormatError("invalid X11 color name in XPM3 color entry")
				}
				clr.Color = color.NRGBA{x11color.RGBA.R, x11color.RGBA.G, x11color.RGBA.B, 0xff}
			}
		case "m":
			clr.Kind = XPM3ColorKindMonochrome
			return FormatError("XPM3 monochrome color not implemented")
		case "g":
			clr.Kind = XPM3ColorKindGrayscale
			return FormatError("XPM3 grayscale color not implemented")
		case "s":
			clr.Kind = XPM3ColorKindSymbolic
			return FormatError("XPM3 symbolic color not implemented")
		default:
			return FormatError("invalid color kind in XPM3 color entry")
		}
		header.colors = append(header.colors, clr)
	}

	d.xpmHeader = header

	return nil
}

func (d *decoder) parseXPM3Pixels() error {
	header := d.xpmHeader.(*XPM3Header)

	img := image.NewNRGBA(image.Rect(0, 0, header.width, header.height))

	y := 0
	for d.scanner.Scan() {
		token := d.scanner.Text()

		if strings.HasPrefix(token, "};") { // Bail once we reach closing block.
			break
		}
		// Tidy up quotations and commas.
		token = strings.TrimPrefix(token, "\"")
		token = strings.TrimSuffix(token, "\",") // This might break things...
		token = strings.TrimSuffix(token, "\"")

		// Step thru according to cPP
		x := 0
		for i := 0; i < len(token); i += header.cPP {
			key := token[i : i+header.cPP]
			color := header.XPM3Color(key)
			if color == nil {
				return FormatError("invalid color key in XPM3 pixel data")
			}
			img.Set(x, y, color.Color)
			x++
		}
		y++
	}

	return nil
}
