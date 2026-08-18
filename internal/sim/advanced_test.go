package sim

import "testing"

func TestMMcKRejects(t *testing.T) {
	cfg := MMcKConfig{Lambda: 2.0, Mu: 1.0, Servers: 1, Capacity: 3, Customers: 1000, Seed: 42}
	st, err := RunMMcK(cfg)
	if err != nil {
		t.Fatalf("MMcK: %v", err)
	}
	if st.Rejected == 0 {
		t.Fatal("expected some rejections with high load")
	}
	if st.RejectRate <= 0 || st.RejectRate >= 1 {
		t.Fatalf("reject rate = %f", st.RejectRate)
	}
}

func TestMMcKInvalidConfig(t *testing.T) {
	_, err := RunMMcK(MMcKConfig{Lambda: -1, Mu: 1, Servers: 1, Capacity: 5, Customers: 10, Seed: 1})
	if err == nil {
		t.Fatal("expected error")
	}
	_, err = RunMMcK(MMcKConfig{Lambda: 1, Mu: 1, Servers: 3, Capacity: 1, Customers: 10, Seed: 1})
	if err == nil {
		t.Fatal("expected error for capacity < servers")
	}
}

func TestImpatienceAbandons(t *testing.T) {
	cfg := ImpatienceConfig{Lambda: 2.0, Mu: 1.0, Servers: 1, Patience: 0.5, Customers: 500, Seed: 42}
	st, err := RunWithImpatience(cfg)
	if err != nil {
		t.Fatalf("impatience: %v", err)
	}
	if st.Abandoned == 0 {
		t.Fatal("expected some abandons with high load and low patience")
	}
	if st.AbandonRate <= 0 {
		t.Fatalf("abandon rate = %f", st.AbandonRate)
	}
}

func TestImpatienceInvalid(t *testing.T) {
	_, err := RunWithImpatience(ImpatienceConfig{Lambda: 1, Mu: 1, Servers: 1, Patience: -1, Customers: 10, Seed: 1})
	if err == nil {
		t.Fatal("expected error")
	}
}
