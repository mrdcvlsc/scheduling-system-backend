package geneticalgorithm_test

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkIntnDedicated(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Intn(10_000_000)
	}
}

func BenchmarkIntnDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Intn(10_000_000)
	}
}

func BenchmarkInt63nDedicated(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Int63n(10_000_000)
	}
}

func BenchmarkInt63nDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int63n(10_000_000)
	}
}

func BenchmarkInt31nDedicated(b *testing.B) {
	// fastest in my machine
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Int31n(10_000_000)
	}
}

func BenchmarkInt31nDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int31n(10_000_000)
	}
}
