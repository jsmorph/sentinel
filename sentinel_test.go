package sentinel

import (
	"context"
	"testing"
)

func TestSentinel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		want = make([]string, 0, 4)
		s    = NewSentinel()
		next *Sentinel
		err  error
	)

	t.Run("write1", func(t *testing.T) {
		v := "hello"
		if next, err = s.Write(v); err != nil {
			t.Fatal(err)
		}
		want = append(want, v)
	})

	t.Run("write2", func(t *testing.T) {
		if _, err = s.Write("again"); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("write3", func(t *testing.T) {
		v := "world"
		if _, err = next.Write(v); err != nil {
			t.Fatal(err)
		}
		want = append(want, v)
	})

	t.Run("consume", func(t *testing.T) {
		for _, v := range want {
			select {
			case <-s.C:
				if s.Value.(string) != v {
					t.Fatal(s.Value)
				}
				s = s.Next
			case <-ctx.Done():
				t.Fatal("no value")
			}
		}
	})
}
