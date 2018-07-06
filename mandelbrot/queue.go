package mandelbrot

import (
	"os"
	"runtime"
	"strconv"
)

// Количество одновременно обрабатываемых процессов
var workers = make(chan struct{}, getWorkersCnt())

// Получит количество одновременно обрабатываемых процессов
func getWorkersCnt() int {

	// Берется и параметров environment
	cnt, err := strconv.Atoi(os.Getenv("WORKERS"))
	if err == nil && cnt > 0 {
		return cnt
	}

	// Если параметр пустой то
	// берем значение по умолчанию как количество процессоров
	return runtime.GOMAXPROCS(0)

}
