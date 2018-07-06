package cache

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
)

// Получаем значение из кэш
func Load(width, height int, id string) (buf []byte, err error) {
	return ioutil.ReadFile(getDirectoryPath() + getSubPath(width, height) + id + ".png")
}

// Загружаем значение в кэш
func Save(width, height int, id string, img image.Image) (buf []byte, err error) {

	// Преобразуем фото в байт массив
	r := new(bytes.Buffer)
	err = png.Encode(r, img)
	if err != nil {
		return
	}

	buf = r.Bytes()

	// Если размер больше лимита, не сохраняем в кэш
	if int64(len(buf)) > getDirectoryLimit() {
		err = errors.New("File size is large than folder limit size.")
		return
	}

	// Если размер увеличиться после сохранение,
	// очищаем папку кэша
	if diskUsage(getDirectoryPath()) + int64(len(buf)) > getDirectoryLimit() {
		err = Clear()
		if err != nil {
			return
		}
	}

	// Создаем нужные директории
	err = os.MkdirAll(getDirectoryPath() + getSubPath(width, height), os.ModePerm)
	if err != nil {
		return
	}

	// Полный пут к сохраняемый файл
	filePath := getDirectoryPath() + getSubPath(width, height) + id + ".png"

	// Сохраняем данные в tmp файл,
	// потом переименуем в нужным нам файл
	err = ioutil.WriteFile(filePath + ".tmp", buf, os.ModePerm)
	if err != nil {
		os.Remove(filePath + ".tmp")
	} else {
		err = os.Rename(filePath + ".tmp", filePath)
	}

	return

}

// Очищает папку кэша
func Clear() error {
	return os.RemoveAll(getDirectoryPath())
}

// Возвращает под директорую из размеров
func getSubPath(width, height int) string {
	return strconv.Itoa(width) + "x" + strconv.Itoa(height) + "/"
}
