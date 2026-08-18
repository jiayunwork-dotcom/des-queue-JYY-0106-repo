package stats

import (
	"math"
	"testing"
)

func TestTheoryMMc(t *testing.T) {
	th, err := TheoryMMc(1.5, 1.0, 2)
	if err != nil {
		t.Fatalf("MMc: %v", err)
	}
	if th.Rho <= 0 || th.Rho >= 1 {
		t.Fatalf("rho = %f", th.Rho)
	}
	if th.Wq <= 0 {
		t.Fatalf("Wq = %f", th.Wq)
	}
}

func TestTheoryMMcUnstable(t *testing.T) {
	_, err := TheoryMMc(3.0, 1.0, 2)
	if err != ErrUnstable {
		t.Fatalf("expected ErrUnstable, got %v", err)
	}
}

func TestErlangBLoss(t *testing.T) {
	// For c=1, a=1: B = 1/(1+1) = 0.5.
	b := ErlangBLoss(1.0, 1)
	if math.Abs(b-0.5) > 1e-6 {
		t.Fatalf("B(1,1) = %f, want 0.5", b)
	}
}

func TestLittlesLaw(t *testing.T) {
	// L = lambda * W => 1 = 0.5 * 2.
	diff := LittlesLaw(0.5, 1.0, 2.0)
	if diff > 1e-9 {
		t.Fatalf("Little's law violation: %f", diff)
	}
}
