package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/placebelg/{width:[0-9]+}/{height:[0-9]+}", PlaceBelgHandler)

	http.Handle("/", r)
	http.ListenAndServe(":30472", nil)
}

func PlaceBelgHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	width, _ := strconv.Atoi(vars["width"])
	height, _ := strconv.Atoi(vars["height"])

	fileBytes := createImg(width, height)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}

func createImg(width int, height int) []byte {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	custom := color.RGBA{10, 12, 140, 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, custom)
		}
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}
