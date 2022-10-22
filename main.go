package main

import (
	"image"
	"image/color"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
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
	users := []string{mux.Vars(r)["gituser"]}
	// msg := mux.Vars(r)["msg"]

	avatar_images := avatar.GetMultiple(users)

	out_width := 320
	out_height := 200
	avatar_width := 80
	avatar_height := 80

	dc := gg.NewContext(out_width, out_height)
	cropped_image := imaging.Fill(avatar_images[0], avatar_width, avatar_height, imaging.Center, imaging.Lanczos)
	upper_jaw := imaging.Crop(cropped_image, image.Rectangle{image.Point{0, 0}, image.Point{79, 39}})
	lower_jaw := imaging.Crop(cropped_image, image.Rectangle{image.Point{0, 40}, image.Point{79, 79}})
	transparent := color.RGBA{0, 0, 0, 0}

	upper_jaw = imaging.Rotate(upper_jaw, 45, transparent)
	lower_jaw = imaging.Rotate(lower_jaw, -45, transparent)
	dc.Push()
	dc.Pop()
	// pat := gg.NewSurfacePattern(cropped_image, gg.RepeatNone)

	dc.SetColor(color.Black)
	dc.DrawRectangle(0, 0, 319, 199)
	dc.Fill()
	// dc.SetFillStyle(pat)
	// dc.RotateAbout(math.Pi/3, 40, 20)
	// dc.Fill()
	dc.DrawImageAnchored(upper_jaw, 0, 100, 0.0, 1.0)
	dc.DrawImageAnchored(lower_jaw, 0, 100, 0.0, 0.0)
	dc.EncodePNG(w)
}
