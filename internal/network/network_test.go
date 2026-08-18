package network

import "testing"

func TestTandemBasic(t *testing.T) {
	cfg := TandemConfig{
		Nodes:     []Node{{Name: "A", Servers: 1, Mu: 2}, {Name: "B", Servers: 1, Mu: 2}},
		Lambda:    0.5,
		Customers: 500,
		Seed:      42,
	}
	res, err := RunTandem(cfg)
	if err != nil {
		t.Fatalf("tandem: %v", err)
	}
	if len(res.Nodes) != 2 {
		t.Fatalf("nodes = %d", len(res.Nodes))
	}
	if res.TotalW <= 0 {
		t.Fatal("total W should be > 0")
	}
	if res.Throughput <= 0 {
		t.Fatal("throughput should be > 0")
	}
}

func TestTandemNoNodes(t *testing.T) {
	_, err := RunTandem(TandemConfig{Lambda: 1, Customers: 10, Seed: 1})
	if err != ErrNoNodes {
		t.Fatalf("expected ErrNoNodes, got %v", err)
	}
}

func TestParallelBasic(t *testing.T) {
	cfg := ParallelConfig{
		Nodes:     []Node{{Name: "S1", Servers: 1, Mu: 1}, {Name: "S2", Servers: 1, Mu: 1}},
		Lambda:    1.5,
		Customers: 500,
		Seed:      42,
	}
	res, err := RunParallel(cfg)
	if err != nil {
		t.Fatalf("parallel: %v", err)
	}
	if len(res.Nodes) != 2 {
		t.Fatalf("nodes = %d", len(res.Nodes))
	}
	if res.AvgW < 0 {
		t.Fatal("avg W should be >= 0")
	}
}

func TestParallelNoNodes(t *testing.T) {
	_, err := RunParallel(ParallelConfig{Lambda: 1, Customers: 10, Seed: 1})
	if err != ErrNoNodes {
		t.Fatalf("expected ErrNoNodes, got %v", err)
	}
}
