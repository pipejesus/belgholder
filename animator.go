package main

import (
	"image"
	"image/color"
	"math"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

type animator struct {
	out_w     int
	out_h     int
	upper_jaw image.Image
	lower_jaw image.Image
	av_w      int
	av_h      int
}

var msg_padding = 20

func NewAnimator(out_w int, out_h int, sprite image.Image) animator {
	avatar_width := 70
	avatar_height := 70

	upper_jaw, lower_jaw := createJaws(avatar_width, avatar_height, sprite)

	return animator{
		out_w,
		out_h,
		upper_jaw,
		lower_jaw,
		avatar_width,
		avatar_height,
	}
}

func createJaws(w int, h int, sprite image.Image) (image.Image, image.Image) {
	cropped_image := imaging.Fill(sprite, w, h, imaging.Center, imaging.Lanczos)
	upper_jaw := imaging.Crop(cropped_image, image.Rectangle{image.Point{0, 0}, image.Point{w, h/2 - 1}})
	lower_jaw := imaging.Crop(cropped_image, image.Rectangle{image.Point{0, h / 2}, image.Point{w, h - 1}})

	return upper_jaw, lower_jaw
}

func (a animator) Debug() image.Image {
	c := gg.NewContext(320, 200*10)
	imgs := a.RenderFrames("Test")
	for i, v := range imgs {
		c.DrawImage(v, 0, i*200)
	}
	return c.Image()
}

func (a animator) RenderFrames(msg string) []image.Image {
	transparent := color.RGBA{0, 0, 0, 0}

	images := []image.Image{}
	time_passed := 0.0
	angle := 0.0

	for i := 0; i < 10; i++ {
		ctx := gg.NewContext(a.out_w, a.out_h)

		upper_jaw := imaging.Rotate(a.upper_jaw, angle, transparent)
		lower_jaw := imaging.Rotate(a.lower_jaw, -1.0*angle, transparent)

		ctx.Clear()
		ctx.SetColor(color.RGBA{255, 255, 255, 255})
		ctx.DrawRectangle(0, 0, 320, 200)
		ctx.Fill()
		ctx.DrawImageAnchored(upper_jaw, 0, 100, 0.0, 1.0)
		ctx.DrawImageAnchored(lower_jaw, 0, 100, 0.0, 0.0)

		// ctx.LoadFontFace("fonts/Comic_CAT.ttf", 16)
		ctx.LoadFontFace("fonts/AlmaMono-Heavy.ttf", 16)

		bb := upper_jaw.Bounds()

		msg_center_x := float64(bb.Max.X + msg_padding)
		msg_center_y := float64(a.out_h / 2)
		msg_max_w := float64(a.out_w - (2 * msg_padding) - a.av_w)

		if len(msg) > 0 {

			wrapped_msg := ctx.WordWrap(msg, msg_max_w)
			m_h := float64(len(wrapped_msg)) * 1.2 * 16.0

			// h := float64(len(lines)) * dc.fontHeight * lineSpacing
			// h -= (lineSpacing - 1) * dc.fontHeight
			m_y := msg_center_y - 0.5*m_h

			ctx.SetRGBA(0.1, 0.7, 0.21, 1.0)
			ctx.DrawRoundedRectangle(msg_center_x-5.0, m_y, msg_max_w+10.0, m_h, 3.0)
			ctx.Fill()
			ctx.DrawRegularPolygon(3, msg_center_x-5.0-5.0, msg_center_y, 10.0, -math.Pi/2)
			ctx.Fill()
			ctx.SetRGBA(0.9, 0.89, 0.91, 1.0)
			ctx.DrawStringWrapped(msg, msg_center_x, msg_center_y, 0, 0.5, msg_max_w, 1.2, gg.AlignRight)
			// ctx.DrawStringAnchored(msg, float64(a.out_w-msg_padding), float64(msg_padding), 1.0, 1.0)
		}

		// ctx.EncodePNG(w)
		img := ctx.Image()

		images = append(images, img)
		// angle = math.Abs(math.Sin(time_passed) * math.Sin(time_passed) * 5.0)
		angle = math.Abs(math.Sin(time_passed) * 5.0)
		time_passed = time_passed + 0.8
	}

	return images

}
