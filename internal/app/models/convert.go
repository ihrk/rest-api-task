package models

type ConvertFunc[T, U any] func(*T) U

func ConvertSlice[T, U any](cf ConvertFunc[T, U], arr []T) []U {
	if len(arr) == 0 {
		return nil
	}

	result := make([]U, 0, len(arr))

	for i := range arr {
		result = append(result, cf(&arr[i]))
	}

	return result
}
