package cache

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"strconv"
	"testing"
)

var _ = os.Setenv("CACHE_LIMIT", "1M")
var testJpg []byte

// empty black jpg
func getTestJpg() (width int, heigth int, id string, buf []byte, err error) {

	width = 2048
	heigth = 2048
	id = "test"

	if len(testJpg) == 0 {

		m := image.NewRGBA(image.Rect(0, 0, width, heigth))

		imgReader := new(bytes.Buffer)
		err = png.Encode(imgReader, m)
		if err != nil {
			return
		}

		testJpg = imgReader.Bytes()

	}

	return width, heigth, id, testJpg, nil

}

func TestCacheSaveLoad(t *testing.T) {

	width, heigth, id, imgBuf, err := getTestJpg()
	if err != nil {
		t.Error(err)
		return
	}

	img, err := png.Decode(bytes.NewReader(imgBuf))
	if err != nil {
		t.Error(err)
		return
	}

	saveImgBuf, err := Save(width, heigth, id, img)
	if err != nil {
		t.Error(err)
		return
	}

	loadImgBuf, err := Load(width, heigth, id)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(saveImgBuf, loadImgBuf) {
		t.Error("Save and load file data not equal")
	}

}

// cache folder limit
func TestCacheLimit(t *testing.T) {

	width, heigth, id, imgBuf, err := getTestJpg()
	if err != nil {
		t.Error(err)
		return
	}

	img, err := png.Decode(bytes.NewReader(imgBuf))
	if err != nil {
		t.Error(err)
		return
	}

	err = Clear()
	if err != nil {
		t.Error(err)
		return
	}

	limit := getDirectoryLimit()
	cnt := 0

	for limit >= 0 {

		cnt++

		_, err = Save(width, heigth, id+strconv.Itoa(cnt), img)
		if err != nil {
			t.Error(err)
			return
		}

		limit -= int64(len(imgBuf))

	}

	if diskUsage(getDirectoryPath()) > getDirectoryLimit() {
		t.Error("Failed to clean the cache folder")
	}

}

// cache folder clear
func TestCacheClear(t *testing.T) {

	err := Clear()
	if err != nil {
		t.Error(err)
	}

}
