package stats

import (
	"strings"
	"testing"

	"des-queue/internal/sim"
)

func TestTheoryMM1(t *testing.T) {
	th, err := TheoryMM1(0.5, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !near(th.Rho, 0.5) || !near(th.L, 1.0) || !near(th.W, 2.0) || !near(th.Wq, 1.0) {
		t.Errorf("theory = %+v, want {0.5 1 2 1}", th)
	}
	for _, l := range []float64{0, -1, 2} {
		if _, err := TheoryMM1(l, 1.0); err == nil {
			t.Errorf("TheoryMM1(%g, 1) expected error", l)
		}
	}
	if _, err := TheoryMM1(0.5, 0); err == nil {
		t.Error("TheoryMM1(0.5, 0) expected error")
	}
}

func TestCompare(t *testing.T) {
	d := Compare(sim.Stats{Rho: 0.55, Wq: 0.9}, Theory{Rho: 0.5, Wq: 1.0})
	if !near(d.RhoErr, 10) || !near(d.WqErrPct, -10) {
		t.Errorf("diff = %+v, want {10 -10}", d)
	}
}

func TestReportNonEmpty(t *testing.T) {
	c := sim.Config{Lambda: 0.5, Mu: 1.0, Servers: 1, Customers: 100, Seed: 42}
	s := sim.Stats{Rho: 0.5, L: 1, W: 2, Wq: 1, Throughput: 0.5, Served: 100}
	r := Report(c, s, Diff{RhoErr: 1.5, WqErrPct: -2.5})
	if len(r) == 0 {
		t.Fatal("report is empty")
	}
	for _, key := range []string{"lambda", "rho", "Wq", "served"} {
		if !strings.Contains(r, key) {
			t.Errorf("report missing %q", key)
		}
	}
}

func near(a, b float64) bool {
	const eps = 1e-9
	d := a - b
	return d < eps && d > -eps
}
