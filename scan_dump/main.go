// Command scan_dump reads user profiles as JSON from
// stardard input and inserts them into the database,
// fetching profile pictures as needed.
package main

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/unixpickle/bumble-dump"
)

const (
	NumPhotoWorkers  = 8
	MaxPhotosPerUser = 2
	MaxImageSize     = 512
)

func main() {
	db, err := bumble.OpenDatabase(bumble.GetConfig())
	if err != nil {
		log.Fatalln("scan_dump:", err)
	}

	photoChan := make(chan *bumble.Photo, 16)
	photoWg := sync.WaitGroup{}
	for i := 0; i < NumPhotoWorkers; i++ {
		photoWg.Add(1)
		go photoDownloader(db, photoChan, &photoWg)
	}
	defer photoWg.Wait()
	defer close(photoChan)

	dec := json.NewDecoder(os.Stdin)
	for {
		var user bumble.User
		if err := dec.Decode(&user); err != nil {
			log.Fatalln("scan_dump:", err)
		}
		db.AddUser(&user)
		photos := user.AllPhotos()
		if len(photos) > MaxPhotosPerUser {
			photos = photos[:MaxPhotosPerUser]
		}
		for _, photo := range photos {
			photoChan <- photo
		}
	}
}

func photoDownloader(db bumble.Database, ch <-chan *bumble.Photo, wg *sync.WaitGroup) {
	defer wg.Done()
	for photo := range ch {
		resp, err := http.Get("https:" + photo.LargeURL)
		if err != nil {
			log.Println("scan_dump:", err)
			continue
		}
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Println("scan_dump:", err)
			continue
		}
		data, err = shrinkPhoto(data)
		if err != nil {
			log.Println("scan_dump:", err)
			continue
		}
		if err := db.AddPhoto(photo, data); err != nil {
			log.Println("scan_dump:", err)
		}
	}
}

func shrinkPhoto(photoData []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(photoData))
	if err != nil {
		return nil, errors.Wrap(err, "shrink photo")
	}
	newImg := resize.Thumbnail(MaxImageSize, MaxImageSize, img, resize.Bilinear)
	var writer bytes.Buffer
	if err := jpeg.Encode(&writer, newImg, &jpeg.Options{Quality: 50}); err != nil {
		return nil, errors.Wrap(err, "shrink photo")
	}
	return writer.Bytes(), nil
}
