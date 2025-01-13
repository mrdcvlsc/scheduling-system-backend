package Utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkShuffleIntnDedicated(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	fmt.Printf("\nIntnDedicatedShuffle | Before: %v\n", arr)

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}

	fmt.Printf("IntnDedicatedShuffle | After: %v\n", arr)
}

func BenchmarkShuffleIntnDefault(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	fmt.Printf("\nIntnDefaultShuffle | Before: %v\n", arr)

	for i := 0; i < b.N; i++ {
		rand.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}

	fmt.Printf("IntnDefaultShuffle | After: %v\n", arr)
}

func BenchmarkShuffleInt63nDedicated(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	fmt.Printf("\nInt63nDedicatedShuffle | Before: %v\n", arr)

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}

	fmt.Printf("Int63nDedicatedShuffle | After: %v\n", arr)
}

func BenchmarkShuffleInt63nDefault(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	fmt.Printf("\nInt63nDefaultShuffle | Before: %v\n", arr)

	for i := 0; i < b.N; i++ {
		rand.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}

	fmt.Printf("Int63nDefaultShuffle | After: %v\n", arr)
}

func BenchmarkShuffleInt31nDedicated(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	fmt.Printf("\nInt31nDedicatedShuffle | Before: %v\n", arr)

	// fastest in my machine
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}

	fmt.Printf("Int31nDedicatedShuffle | After: %v\n", arr)
}

func BenchmarkShuffleInt31nDefault(b *testing.B) {
	arr := make([]uint64, 0, 8)

	for i := 0; i < 10; i++ {
		arr = append(arr, uint64(i+1))
	}

	fmt.Printf("\nInt31nDefaultShuffle | Before: %v\n", arr)

	for i := 0; i < b.N; i++ {
		rand.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}

	fmt.Printf("Int31nDefaultShuffle | After: %v\n", arr)
}
