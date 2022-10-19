package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/google/go-github/v47/github"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/placebelg/{width:[0-9]+}/{height:[0-9]+}/{users:[a-z]+}", PlaceBelgHandler)
	r.HandleFunc("/hero-badge", HeroBadgeHandler).Queries("gitusers", "{gitusers}")

	http.Handle("/", r)
	http.ListenAndServe(":30472", nil)
}

func HeroBadgeHandler(w http.ResponseWriter, r *http.Request) {

	users := strings.Split(mux.Vars(r)["gitusers"], ",")

	if len(users) == 0 {
		return
	}

	avatar_images := getCachedAvatarImages(users)

	dst_width := 128
	dst_height := 128

	out_width := len(avatar_images) * dst_width
	out_height := dst_height

	dc := gg.NewContext(out_width, out_height)

	x := 0
	for _, avatar := range avatar_images {
		cropped_image := imaging.Fill(avatar, dst_width, dst_height, imaging.Center, imaging.Lanczos)
		dc.DrawImage(cropped_image, x, 0)
		x += cropped_image.Bounds().Dx()
	}

	dc.DrawCircle(float64(out_width/2), float64(out_height/2), float64(out_height/4))
	dc.SetRGBA(0, 0, 0, 0.6)
	dc.Fill()
	dc.EncodePNG(w)

	// buf := new(bytes.Buffer)
	// png.Encode(buf, imgAv)
	// return buf.Bytes()

	return
}

func getCachedAvatarImages(users []string) []image.Image {
	avatar_images := make([]image.Image, 0)

	for _, gitUserName := range users {
		img, err := downloadUserAvatar(gitUserName)
		if err == nil {
			avatar_images = append(avatar_images, img)
		}
	}

	return avatar_images
}

func downloadUserAvatar(githubUserName string) (image.Image, error) {
	out_file_name := "avatars/" + githubUserName + ".png"

	if _, err := os.Stat(out_file_name); err == nil {
		cached_avatar, _ := os.Open(out_file_name)
		return png.Decode(cached_avatar)
	}

	client := github.NewClient(nil)

	userinfo, _, err := client.Users.Get(context.Background(), githubUserName)

	if err != nil {
		return nil, errors.New("no")
	}

	fmt.Println(userinfo)

	avatarUrl := userinfo.AvatarURL
	response, _ := http.Get(*avatarUrl)
	defer response.Body.Close()

	image_type := response.Header.Get("Content-Type")

	out_file, err := os.Create(out_file_name)
	defer out_file.Close()

	switch image_type {
	case "image/png":
		_, _ = io.Copy(out_file, response.Body)
		decoded, _ := png.Decode(response.Body)
		return decoded, nil

	case "image/jpeg":
		img_jpg, _ := jpeg.Decode(response.Body)
		buf := new(bytes.Buffer)
		png.Encode(buf, img_jpg)
		out_file.Write(buf.Bytes())
		return img_jpg, nil

	case "image/gif":
		img_gif, _ := gif.Decode(response.Body)
		buf := new(bytes.Buffer)
		png.Encode(buf, img_gif)
		out_file.Write(buf.Bytes())
		return img_gif, nil
	default:
		return nil, errors.New("Unknown image type")
	}

}

func PlaceBelgHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	width, _ := strconv.Atoi(vars["width"])
	height, _ := strconv.Atoi(vars["height"])
	users := vars["users"]

	fileBytes := createImg(width, height, users)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}

func createImg(width int, height int, users string) []byte {
	client := github.NewClient(nil)

	czoko, _, _ := client.Users.Get(context.Background(), users)

	avatarUrl := czoko.AvatarURL
	response, _ := http.Get(*avatarUrl)

	defer response.Body.Close()

	// body, _ := ioutil.ReadAll(response.Body)

	imgAv, _ := jpeg.Decode(response.Body)

	// upLeft := image.Point{0, 0}
	// lowRight := image.Point{width, height}

	// img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// custom := color.RGBA{10, 12, 140, 0xff}

	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		img.Set(x, y, custom)
	// 	}
	// }

	buf := new(bytes.Buffer)
	png.Encode(buf, imgAv)
	return buf.Bytes()
}
