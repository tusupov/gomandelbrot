package mandelbrot

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
)

type point struct {
	X			int		`json:"x"`
	Y			int		`json:"y"`
	priority	int
}

type sizes map[string]point

// Значение размеров по умолчанию
var defaulSizes = sizes{
	"small": {
		X: 64,
		Y: 64,
	},
}

var sizesList = loadSizes()

// Загружаем размеры из конфиг файла
func loadSizes() (s sizes) {

	defer func() {
		// Расчет приоритеты размеров
		s = calcSizesPriority(s)
	}()

	configFile, err := os.Open("config/sizes.json")
	if err != nil {
		s = defaulSizes
		return
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&s)
	if err != nil || len(s) == 0 {
		s = defaulSizes
		return
	}

	return

}

// Расчет приоритеты размеров
func calcSizesPriority(s sizes) sizes {

	type pointN struct {
		name	string
		x, y 	int
	}

	sizeList := make([]pointN, 0)

	for k, v := range s {
		sizeList = append(
			sizeList, pointN{
				name: 	k,
				x:		v.X,
				y:		v.Y,
			},
		)
	}

	sort.Slice(sizeList, func(i, j int) bool { return sizeList[i].x  < sizeList[j].x })

	for k, v := range sizeList {
		s[v.name] = point{
			X: 			v.x,
			Y: 			v.y,
			priority: 	k + 1,
		}
	}

	return s

}

// Получит значение размера
func GetSize(rec string) (int, int, int) {

	if size, ok := sizesList[strings.ToLower(rec)]; ok {
		return size.X, size.Y, size.priority
	}

	return GetSize("small")

}
