package x11

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"image"
	"image/color"

	"github.com/patrislav/wmapi/domain"
)

func (c *Client) GetWindowIcon(id domain.WindowID) (image.Image, error) {
	icons, _ := ewmh.WmIconGet(c.xu, xproto.Window(id))
	if len(icons) <= 0 {
		return nil, fmt.Errorf("no icon")
	}
	icon := icons[0]
	img := image.NewRGBA(image.Rect(0, 0, int(icon.Width), int(icon.Height)))
	for i, argb := range icon.Data {
		x := i % int(icon.Width)
		y := i / int(icon.Height)
		col := color.RGBA{
			A: uint8((0xff000000 & argb) >> 24),
			R: uint8((0x00ff0000 & argb) >> 16),
			G: uint8((0x0000ff00 & argb) >> 8),
			B: uint8(0x000000ff & argb),
		}
		img.Set(x, y, col)
	}
	return img, nil
}
