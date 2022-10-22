package main

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v47/github"
)

type Avatar struct {
}

func (a Avatar) GetMultiple(users []string) []image.Image {
	avatar_images := make([]image.Image, 0)

	for _, gitUserName := range users {
		img, err := a.GetOne(gitUserName)
		if err == nil {
			avatar_images = append(avatar_images, img)
		}
	}

	return avatar_images
}

func (a Avatar) GetOne(githubUserName string) (image.Image, error) {
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
