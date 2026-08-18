package metrics

import "testing"

func TestRecorderBasic(t *testing.T) {
	r := NewRecorder(1.0)
	r.Record(Sample{Time: 1, QueueLen: 2, BusyCount: 1, IdleCount: 1, InSystem: 3})
	r.Record(Sample{Time: 2, QueueLen: 4, BusyCount: 2, IdleCount: 0, InSystem: 6})
	if r.Count() != 2 {
		t.Fatalf("count = %d", r.Count())
	}
	if r.AvgQueueLen() != 3 {
		t.Fatalf("avg queue = %f", r.AvgQueueLen())
	}
	if r.MaxQueueLen() != 4 {
		t.Fatalf("max queue = %d", r.MaxQueueLen())
	}
}

func TestUtilization(t *testing.T) {
	r := NewRecorder(1.0)
	r.Record(Sample{BusyCount: 1})
	r.Record(Sample{BusyCount: 2})
	avg := r.AvgUtilization(2)
	if avg != 0.75 {
		t.Fatalf("avg util = %f, want 0.75", avg)
	}
	peak := r.PeakUtilization(2)
	if peak != 1.0 {
		t.Fatalf("peak util = %f, want 1.0", peak)
	}
}

func TestPercentile(t *testing.T) {
	r := NewRecorder(1.0)
	for i := 0; i < 100; i++ {
		r.Record(Sample{QueueLen: i + 1})
	}
	p50 := r.QueueLenPercentile(50)
	if p50 < 49 || p50 > 51 {
		t.Fatalf("p50 = %f", p50)
	}
}

func TestTimeSeries(t *testing.T) {
	r := NewRecorder(1.0)
	r.Record(Sample{Time: 1, QueueLen: 5})
	r.Record(Sample{Time: 2, QueueLen: 3})
	ts := r.QueueLenTimeSeries()
	if len(ts) != 2 || ts[0].Value != 5 {
		t.Fatalf("ts = %v", ts)
	}
}
