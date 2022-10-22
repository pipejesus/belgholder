package main

import (
	"image"
	"image/gif"
	"image/png"
	"net/http"

	"github.com/andybons/gogif"
	"github.com/gorilla/mux"
)

var avatar = Avatar{}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/czo", ComicHandler).Queries("gituser", "{gituser}", "msg", "{msg}")
	http.Handle("/", r)
	http.ListenAndServe(":30472", nil)
}

func ComicHandler(w http.ResponseWriter, r *http.Request) {
	gituser := mux.Vars(r)["gituser"]
	// msg := mux.Vars(r)["msg"]

	avatar_image, err := avatar.GetOne(gituser)

	if err != nil {
		return
	}

	out_width := 320
	out_height := 200

	animator := NewAnimator(out_width, out_height, avatar_image)
	// images := animator.RenderFrames()
	img_debug := animator.Debug()
	png.Encode(w, img_debug)
	// gif.EncodeAll(w, QuantizeImagesAndAddToGif(images))

	return
}

func QuantizeImagesAndAddToGif(frames []image.Image) *gif.GIF {
	out_gif := &gif.GIF{
		LoopCount: 100,
	}

	for _, simage := range frames {
		bounds := simage.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: 64}
		quantizer.Quantize(palettedImage, bounds, simage, image.ZP)

		// Add new frame to animated GIF
		out_gif.Image = append(out_gif.Image, palettedImage)
		out_gif.Delay = append(out_gif.Delay, 100)
	}

	return out_gif
}
