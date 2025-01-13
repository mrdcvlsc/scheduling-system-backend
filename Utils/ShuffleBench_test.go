package Utils

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkShuffleIntnDedicated(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
}

func BenchmarkShuffleIntnDefault(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	for i := 0; i < b.N; i++ {
		rand.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
}

func BenchmarkShuffleInt63nDedicated(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
}

func BenchmarkShuffleInt63nDefault(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	for i := 0; i < b.N; i++ {
		rand.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
}

func BenchmarkShuffleInt31nDedicated(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	// fastest in my machine
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
}

func BenchmarkShuffleInt31nDefault(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	for i := 0; i < b.N; i++ {
		rand.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
}
