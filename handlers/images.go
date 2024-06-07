package handlers

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/nfnt/resize"
	"testTask/middleware"
	"testTask/models"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Center crop to square
	var size int
	if img.Bounds().Dx() > img.Bounds().Dy() {
		size = img.Bounds().Dy()
	} else {
		size = img.Bounds().Dx()
	}

	x0 := (img.Bounds().Dx() - size) / 2
	y0 := (img.Bounds().Dy() - size) / 2
	squareImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(x0, y0, x0+size, y0+size))

	userDir := filepath.Join("images", strconv.Itoa(int(userID)))
	err = os.MkdirAll(userDir, os.ModePerm)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for i := 30; i <= 80; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(size int) {
			defer wg.Done()
			defer func() { <-sem }()

			resizedImg := resize.Resize(uint(size), uint(size), squareImg, resize.Lanczos3)
			filePath := filepath.Join(userDir, strconv.Itoa(size)+"x"+strconv.Itoa(size)+".jpg")
			outFile, err := os.Create(filePath)
			if err != nil {
				return
			}

			defer func(outFile *os.File) {
				err = outFile.Close()
				if err != nil {

				}
			}(outFile)
			err = jpeg.Encode(outFile, resizedImg, nil)
			if err != nil {
				return
			}

			db.Create(&models.Image{UserID: userID, URL: filePath})
		}(i)
	}

	wg.Wait()
	w.WriteHeader(http.StatusOK)
}

func GetImages(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	var images []models.Image
	db.Where("user_id = ?", userID).Find(&images)

	urls := make([]string, len(images))
	for i, img := range images {
		urls[i] = img.URL
	}

	err := json.NewEncoder(w).Encode(urls)
	if err != nil {
		return
	}
}
