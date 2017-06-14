package arn

import (
	"image"
	"os"

	"github.com/chai2010/webp"
)

var OriginalImageExtensions = []string{
	".jpg",
	".png",
	".gif",
}

const (
	AvatarSmallSize = 100
	AvatarMaxSize   = 560
)

// LoadImage loads an image from the given path.
func LoadImage(path string) (img image.Image, format string, err error) {
	f, openErr := os.Open(path)

	if openErr != nil {
		return nil, "", openErr
	}

	img, format, decodeErr := image.Decode(f)

	if decodeErr != nil {
		return nil, "", decodeErr
	}

	return img, format, nil
}

// SaveWebP saves an image as a file in WebP format.
func SaveWebP(img image.Image, out string, quality float32) error {
	file, writeErr := os.Create(out)

	if writeErr != nil {
		return writeErr
	}

	encodeErr := webp.Encode(file, img, &webp.Options{
		Quality: quality,
	})

	return encodeErr
}

// FindFileWithExtension tries to test different file extensions.
func FindFileWithExtension(baseName string, dir string, extensions []string) string {
	for _, ext := range extensions {
		if _, err := os.Stat(dir + baseName + ext); !os.IsNotExist(err) {
			return dir + baseName + ext
		}
	}

	return ""
}