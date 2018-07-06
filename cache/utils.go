package cache

import (
"fmt"
"os"


)

// get directory size
func diskUsage(directory string) (size int64) {

	dir, err := os.Open(directory)
	if err != nil {
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			size += diskUsage(fmt.Sprintf("%s/%s", directory, file.Name()))
		} else {
			size += file.Size()
		}
	}

	return

}
