package utils

import "github.com/pkg/errors"

func ReduceSlice[T any](inputs []T, reduceFunc func(prev T, cur T) (T, error)) (T, error) {
	if len(inputs) == 0 {
		return *new(T), nil
	}
	res := inputs[0]
	for _, v := range inputs[1:] {
		temp, err := reduceFunc(res, v)
		if err != nil {
			return *new(T), errors.WithStack(err)
		}
		res = temp
	}
	return res, nil
}

func MaxSlice[T any](inputs []T, isAGreater func(a T, b T) (bool, error)) (T, error) {
	return ReduceSlice(inputs, func(prev T, cur T) (T, error) {
		isPrevGreater, err := isAGreater(prev, cur)
		if err != nil {
			return *new(T), errors.WithStack(err)
		}
		if isPrevGreater {
			return prev, nil
		}
		return cur, nil
	})
}
