package xpm

import (
	"strconv"
	"strings"
)

const (
	seenFormat = iota
	seenWidth
	seenHeight
	seenNColors
	seenCharsPerPixel
)

func (d *decoder) parseXPM1Metadata() error {
	xpmHeader := XPM1Header{}
	lastSeen := seenFormat
	// Step thru format token since it is stored from the earlier type check.
	words := strings.Split(d.lastLine, " ")
	if len(words) < 3 {
		return FormatError("too few words in #define")
	}
	formatValue, err := strconv.Atoi(words[2])
	if err != nil {
		return err
	}
	xpmHeader.format = uint8(formatValue)
	// Now get the rest of our defines.
	for d.scanner.Scan() {
		token := d.scanner.Text()
		if strings.HasPrefix(token, "#define") && lastSeen < seenCharsPerPixel {
			words := strings.Split(d.lastLine, " ")
			if len(words) < 3 {
				return FormatError("too few words in #define")
			}
			value, err := strconv.Atoi(words[2])
			if err != nil {
				return err
			}

			if lastSeen == seenFormat {
				xpmHeader.width = uint(value)
				lastSeen = seenWidth
			} else if lastSeen == seenWidth {
				xpmHeader.height = uint(value)
				lastSeen = seenHeight
			} else if lastSeen == seenHeight {
				xpmHeader.nColors = uint(value)
				lastSeen = seenNColors
			} else if lastSeen == seenNColors {
				xpmHeader.cPP = uint8(value)
				d.xpmHeader = xpmHeader
				lastSeen = seenCharsPerPixel
				// End of header, move on!
				break
			}
		} else {
			return FormatError("invalid XPM data")
		}
	}
	if err := d.scanner.Err(); err != nil {
		return err
	}
	return nil
}
