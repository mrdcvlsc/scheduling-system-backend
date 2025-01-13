package Utils

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkRngIntnDedicated(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Intn(10_000_000)
	}
}

func BenchmarkRngIntnDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Intn(10_000_000)
	}
}

func BenchmarkRngInt63nDedicated(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Int63n(10_000_000)
	}
}

func BenchmarkRngInt63nDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int63n(10_000_000)
	}
}

func BenchmarkRngInt31nDedicated(b *testing.B) {
	// fastest in my machine
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < b.N; i++ {
		rng.Int31n(10_000_000)
	}
}

func BenchmarkRngInt31nDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int31n(10_000_000)
	}
}
