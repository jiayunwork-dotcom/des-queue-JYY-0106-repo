package event

import "testing"

func TestFELOrder(t *testing.T) {
	q := &Queue{}
	if _, ok := q.Pop(); ok {
		t.Fatal("empty Pop should return false")
	}
	if q.Len() != 0 {
		t.Fatalf("empty Len = %d, want 0", q.Len())
	}

	q.Push(Event{Time: 3.0, Kind: Departure, Customer: 1})
	q.Push(Event{Time: 1.0, Kind: Arrival, Customer: 2})
	q.Push(Event{Time: 2.0, Kind: Arrival, Customer: 3})
	q.Push(Event{Time: 2.0, Kind: Departure, Customer: 4})

	want := []Event{
		{Time: 1.0, Kind: Arrival, Customer: 2},
		{Time: 2.0, Kind: Arrival, Customer: 3},
		{Time: 2.0, Kind: Departure, Customer: 4},
		{Time: 3.0, Kind: Departure, Customer: 1},
	}
	for i, w := range want {
		e, ok := q.Pop()
		if !ok {
			t.Fatalf("Pop %d: expected event, got empty", i)
		}
		if e != w {
			t.Fatalf("Pop %d = %+v, want %+v", i, e, w)
		}
	}
	if _, ok := q.Pop(); ok {
		t.Fatal("Pop after drain should return false")
	}
}
