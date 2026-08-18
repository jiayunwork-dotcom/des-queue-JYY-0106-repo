package sim

import (
	"math"
	"testing"
)

func TestValidateErrors(t *testing.T) {
	cases := []struct {
		name string
		c    Config
	}{
		{"lambda<=0", Config{Lambda: 0, Mu: 1, Servers: 1, Customers: 10}},
		{"mu<=0", Config{Lambda: 0.5, Mu: -1, Servers: 1, Customers: 10}},
		{"servers<1", Config{Lambda: 0.5, Mu: 1, Servers: 0, Customers: 10}},
		{"customers<1", Config{Lambda: 0.5, Mu: 1, Servers: 1, Customers: 0}},
		{"unstable MM1", Config{Lambda: 2, Mu: 1, Servers: 1, Customers: 10}},
	}
	for _, tc := range cases {
		if err := Validate(tc.c); err == nil {
			t.Errorf("%s: expected error, got nil", tc.name)
		}
	}
	if err := Validate(Config{Lambda: 2, Mu: 1, Servers: 3, Customers: 10}); err != nil {
		t.Errorf("MMc should not check stability, got %v", err)
	}
}

func TestMM1Unstable(t *testing.T) {
	_, err := RunMM1(Config{Lambda: 2, Mu: 1, Servers: 1, Customers: 100, Seed: 1})
	if err == nil {
		t.Fatal("expected unstable error, got nil")
	}
}

func TestMM1MatchesTheory(t *testing.T) {
	cfg := Config{Lambda: 0.5, Mu: 1.0, Servers: 1, Customers: 8000, Seed: 42}
	s, err := RunMM1(cfg)
	if err != nil {
		t.Fatalf("RunMM1: %v", err)
	}
	if math.Abs(s.Rho-0.5)/0.5 >= 0.10 {
		t.Errorf("rho = %.4f, want within 10%% of 0.5", s.Rho)
	}
	if math.Abs(s.Wq-1.0)/1.0 >= 0.10 {
		t.Errorf("Wq = %.4f, want within 10%% of 1.0", s.Wq)
	}
}

func TestMMcIdleServers(t *testing.T) {
	cfg := Config{Lambda: 0.5, Mu: 1.0, Servers: 1, Customers: 4000, Seed: 7}
	s, err := RunMMc(cfg, 3)
	if err != nil {
		t.Fatalf("RunMMc: %v", err)
	}
	if s.Rho >= 0.5 {
		t.Errorf("low-load MMc rho = %.4f, want clearly < 1", s.Rho)
	}
	if s.Served != 4000 {
		t.Errorf("served = %d, want 4000", s.Served)
	}
}

func TestExpRandDeterministic(t *testing.T) {
	a := NewExpRand(42, 1.0)
	b := NewExpRand(42, 1.0)
	for i := 0; i < 100; i++ {
		if a.Next() != b.Next() {
			t.Fatalf("sequence diverged at %d", i)
		}
	}
}
