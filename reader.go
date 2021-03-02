package xpm

import (
	"bufio"
	"image"
	"io"
	"strings"
)

// FormatError represents an invalid format error for the XPM format.
type FormatError string

// Error returns the string of the error.
func (e FormatError) Error() string { return "xpm: invalid format: " + string(e) }

// State machine constants
const (
	seenNothing = iota
	seenXPM
	seenMeta
	seenColors
	seenPixels
)

type XPMType uint8

const (
	NotXPM XPMType = iota
	XPM1
	XPM2
	XPM3
)

type decoder struct {
	scanner   *bufio.Scanner
	image     image.Image
	colors    map[string][]Color
	xpmType   XPMType
	xpmHeader XPMHeader
	lastLine  string
}

// parseMetadata parses the information from /* XPM */ up to /* colors */.
func (d *decoder) parseType() error {
	for d.scanner.Scan() {
		token := d.scanner.Text()
		if strings.HasPrefix(token, "#define") { // XPM1 (or XBM...)
			// Store token for later since XPM1 doesn't have a clean type header.
			d.lastLine = token
			d.xpmType = XPM1
			break
		} else if strings.HasPrefix(token, "! XPM2") { // XPM2
			d.xpmType = XPM2
			break
		} else if strings.HasPrefix(token, "/* XPM */") { // XPM3
			d.xpmType = XPM3
			break
		}
	}
	if err := d.scanner.Err(); err != nil {
		return err
	}
	if d.xpmType == NotXPM {
		return FormatError("could not find XPM type declaration")
	}
	return nil
}

func (d *decoder) parseMetadata() error {
	if d.xpmType == XPM1 {
		return d.parseXPM1Metadata()
	} else if d.xpmType == XPM2 {
		return d.parseXPM2Metadata()
	} else if d.xpmType == XPM3 {
		return d.parseXPM3Metadata()
	}
	return FormatError("metadata cannot be parsed before type")
}

func (d *decoder) parsePixels() error {
	return nil
}

func Decode(r io.Reader) (image.Image, error) {
	d := &decoder{
		scanner: bufio.NewScanner(r),
	}
	if err := d.parseType(); err != nil {
		return nil, err
	}
	if err := d.parseMetadata(); err != nil {
		return nil, err
	}
	if err := d.parsePixels(); err != nil {
		return nil, err
	}
	return d.image, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	d := &decoder{
		scanner: bufio.NewScanner(r),
	}
	if err := d.parseType(); err != nil {
		return image.Config{}, err
	}
	if err := d.parseMetadata(); err != nil {
		return image.Config{}, err
	}
	return image.Config{
		// TODO
	}, nil
}

func init() {
	image.RegisterFormat("xpm", "#define", Decode, DecodeConfig)   // XPM1
	image.RegisterFormat("xpm", "! XPM2", Decode, DecodeConfig)    // XPM2
	image.RegisterFormat("xpm", "/* XPM */", Decode, DecodeConfig) // XPM3
}
