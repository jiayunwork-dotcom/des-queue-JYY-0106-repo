package arrival

import (
	"math"
	"testing"
)

func TestPoissonMean(t *testing.T) {
	p := NewPoisson(42, 2.0)
	sum := 0.0
	n := 10000
	for i := 0; i < n; i++ {
		sum += p.Next()
	}
	mean := sum / float64(n)
	// Expected mean = 1/rate = 0.5.
	if math.Abs(mean-0.5) > 0.05 {
		t.Fatalf("mean = %f, want ~0.5", mean)
	}
}

func TestBatchBasic(t *testing.T) {
	b := NewBatch(42, 1.0, 3)
	// First call gets interval > 0, then 2 zeros.
	first := b.Next()
	if first <= 0 {
		t.Fatalf("first interval = %f, want > 0", first)
	}
	if b.Next() != 0 {
		t.Fatal("second in batch should be 0")
	}
	if b.Next() != 0 {
		t.Fatal("third in batch should be 0")
	}
	// Next batch starts with positive interval.
	fourth := b.Next()
	if fourth <= 0 {
		t.Fatalf("new batch interval = %f, want > 0", fourth)
	}
}

func TestDeterministic(t *testing.T) {
	d := NewDeterministic(2.5)
	for i := 0; i < 10; i++ {
		if d.Next() != 2.5 {
			t.Fatal("deterministic should return constant")
		}
	}
	if d.Name() != "deterministic" {
		t.Fatalf("name = %q", d.Name())
	}
}

func TestErlangMean(t *testing.T) {
	e := NewErlang(42, 3, 2.0)
	sum := 0.0
	n := 10000
	for i := 0; i < n; i++ {
		sum += e.Next()
	}
	mean := sum / float64(n)
	// Expected mean = k/rate = 3/2 = 1.5.
	if math.Abs(mean-1.5) > 0.1 {
		t.Fatalf("mean = %f, want ~1.5", mean)
	}
}

func TestTimeVaryingPositive(t *testing.T) {
	tv := NewTimeVarying(42, 1.0, 0.5, 10.0)
	for i := 0; i < 100; i++ {
		interval := tv.Next()
		if interval < 0 {
			t.Fatalf("interval %d = %f < 0", i, interval)
		}
	}
}

func TestGeneratorInterface(t *testing.T) {
	generators := []Generator{
		NewPoisson(1, 1.0),
		NewBatch(1, 1.0, 2),
		NewDeterministic(1.0),
		NewErlang(1, 2, 1.0),
		NewTimeVarying(1, 1.0, 0.3, 5.0),
	}
	for _, g := range generators {
		if g.Name() == "" {
			t.Fatal("empty name")
		}
		if g.Next() < 0 {
			t.Fatal("negative interval")
		}
	}
}
