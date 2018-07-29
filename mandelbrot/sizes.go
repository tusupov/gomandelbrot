package mandelbrot

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"
)

var (
	errSizeNotFound = errors.New("Size not found")
)

var (
	sizeList   = make(map[string]Size)
	inQueueSum = make([]int32, 0)
)

type Size struct {
	X        int `json:"size"`
	priority int
}

func LoadSize(configPath string) (n int, err error) {

	size, err := loadFile(configPath)
	if err != nil {
		return
	}

	n = len(size)

	sizeList = calcSizesPriority(size)
	inQueueSum = make([]int32, n)

	return

}

// Load from file
func loadFile(configPath string) (s map[string]Size, err error) {

	configFile, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&s)
	if err != nil {
		return
	}
	if len(s) == 0 {
		err = errors.New("Size list is empty")
		return
	}

	return

}

// Calculate sizes priority
func calcSizesPriority(s map[string]Size) map[string]Size {

	type sizeN struct {
		name string
		x    int
	}

	sizeList := make([]sizeN, 0)

	for k, v := range s {
		sizeList = append(
			sizeList, sizeN{
				name: k,
				x:    v.X,
			},
		)
	}

	sort.Slice(sizeList, func(i, j int) bool { return sizeList[i].x < sizeList[j].x })

	for k, v := range sizeList {
		s[v.name] = Size{
			X:        v.x,
			priority: k,
		}
	}

	return s

}

// GetSize
func GetSize(rec string) (int, int, error) {

	if size, ok := sizeList[strings.ToLower(rec)]; ok {
		return size.X, size.priority, nil
	}

	return 0, 0, errSizeNotFound

}
