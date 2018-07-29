package mandelbrot

import "testing"

var testSizeList = []string{"small", "medium", "big"}

func TestLoadSize(t *testing.T) {

	err := LoadSize("../sizes.json")
	if err != nil {
		t.Fatal(err)
	}

}

func TestSetWorkers(t *testing.T) {

	err := SetWorkers(10)
	if err != nil {
		t.Fatal(err)
	}

	err = SetWorkers(0)
	if err == nil {
		t.Fatal("Error expected")
	}

}

func TestNew(t *testing.T) {

	for _, s := range testSizeList {
		_, err := New(s)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestGetSize(t *testing.T) {

	for _, s := range testSizeList {
		_, _, err := GetSize(s)
		if err != nil {
			t.Fatal(err)
		}
	}

}
