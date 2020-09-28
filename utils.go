package main

import (
	"image/color"
	"io"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
)

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func imageUpsideDown(filepath string) error {
	img, err := imaging.Open(filepath)
	if err != nil {
		return errors.Wrap(err, "failed to load image")
	}

	// img, degrees, set a color to the background
	img = imaging.Rotate(img, 180, color.RGBA{0, 0, 0, 1})
	err = imaging.Save(img, filepath)
	if err != nil {
		return errors.Wrap(err, "failed to save image")
	}
	return err
}
