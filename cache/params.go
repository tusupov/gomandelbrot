package cache

import (
	"os"
	"strconv"
	"strings"
)

const (

	// Значение по умолчанию для папки кэша
	defaultDirectoryPath = "./cache/" // folder

	// Значение по умолчанию для лимита кэша
	defaultDirectoryLimit = 16 << 30 // 16G

)

var directoryPath *string

// Возвращает папку кэша
func getDirectoryPath() string {

	if directoryPath == nil {

		// create new
		directoryPath = new(string)

		// load value from env
		*directoryPath = strings.TrimSpace(os.Getenv("CACHE_PATH"))

		if len(*directoryPath) == 0 {
			*directoryPath = defaultDirectoryPath
		}

	}

	return *directoryPath

}

var directoryLimit *int64

// Возвращает лимит ограничение для папки кэша
func getDirectoryLimit() int64 {

	if directoryLimit == nil {

		// create new
		directoryLimit = new(int64)

		// load value from env
		limitStr := strings.TrimSpace(os.Getenv("CACHE_LIMIT"))

		if len(limitStr) > 0 {
			// has limit

			// widthout unit
			limit, err := strconv.ParseInt(limitStr, 10, 0)
			if err == nil {

				*directoryLimit = limit

			} else {
				// with unit

				value := limitStr[:len(limitStr)-1]
				unit := limitStr[len(limitStr)-1]

				var unitValue byte = 0

				switch unit {

				case 'K', 'k':
					unitValue = 10

				case 'M', 'm':
					unitValue = 20

				case 'G', 'g':
					unitValue = 30

				case 'T', 't':
					unitValue = 40

				}

				limit, err := strconv.Atoi(value)
				if err == nil && unitValue > 0 {
					*directoryLimit = int64(limit) << unitValue
				}
			}

		}

		if *directoryLimit <= 0 {
			*directoryLimit = defaultDirectoryLimit
		}

	}

	return *directoryLimit

}
