package mypkg1

import "fmt"

func MyFunc(inputValue int) (int, error) {
	if inputValue > 1000 {
		return 0, fmt.Errorf("do not accept inputValue > 1000: %v", inputValue)
	}
	return inputValue * 2, nil
}
