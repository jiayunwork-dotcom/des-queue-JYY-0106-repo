package queue

import "testing"

func TestFIFOBasic(t *testing.T) {
	q := NewFIFO(0)
	_ = q.Enqueue(Customer{ID: 1})
	_ = q.Enqueue(Customer{ID: 2})
	if q.Len() != 2 {
		t.Fatalf("len = %d", q.Len())
	}
	c, _ := q.Dequeue()
	if c.ID != 1 {
		t.Fatalf("dequeue = %d, want 1", c.ID)
	}
}

func TestFIFOCapacity(t *testing.T) {
	q := NewFIFO(2)
	_ = q.Enqueue(Customer{ID: 1})
	_ = q.Enqueue(Customer{ID: 2})
	if !q.IsFull() {
		t.Fatal("should be full")
	}
	if err := q.Enqueue(Customer{ID: 3}); err != ErrFull {
		t.Fatalf("expected ErrFull, got %v", err)
	}
}

func TestFIFOEmpty(t *testing.T) {
	q := NewFIFO(0)
	if !q.IsEmpty() {
		t.Fatal("should be empty")
	}
	_, err := q.Dequeue()
	if err != ErrEmpty {
		t.Fatalf("expected ErrEmpty, got %v", err)
	}
}

func TestPriorityQueue(t *testing.T) {
	q := NewPriority(0)
	_ = q.Enqueue(Customer{ID: 1, Priority: 3})
	_ = q.Enqueue(Customer{ID: 2, Priority: 1})
	_ = q.Enqueue(Customer{ID: 3, Priority: 2})
	c, _ := q.Dequeue()
	if c.ID != 2 {
		t.Fatalf("dequeue = %d, want 2 (highest priority)", c.ID)
	}
}

func TestSJF(t *testing.T) {
	q := NewSJF(0)
	_ = q.Enqueue(Customer{ID: 1, ServiceDur: 5.0})
	_ = q.Enqueue(Customer{ID: 2, ServiceDur: 1.0})
	_ = q.Enqueue(Customer{ID: 3, ServiceDur: 3.0})
	c, _ := q.Dequeue()
	if c.ID != 2 {
		t.Fatalf("SJF dequeue = %d, want 2 (shortest)", c.ID)
	}
}

func TestStrategyInterface(t *testing.T) {
	strategies := []Strategy{
		NewFIFO(10),
		NewPriority(10),
		NewSJF(10),
	}
	for _, s := range strategies {
		_ = s.Enqueue(Customer{ID: 1})
		if s.Len() != 1 {
			t.Fatal("len should be 1")
		}
		_, _ = s.Dequeue()
		if !s.IsEmpty() {
			t.Fatal("should be empty after dequeue")
		}
	}
}
