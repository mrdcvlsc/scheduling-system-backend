package Utils

import "fmt"

func RemoveChunkInSlice[T any](slice []T, start, size int) ([]T, error) {
	end := start + size

	if start < 0 || size < 0 || end > len(slice) {
		return nil, fmt.Errorf("RemoveChunkInSlice invalid parameters: start=%d, size=%d, len=%d", start, size, len(slice))
	}

	return append(slice[:start], slice[end:]...), nil
}
