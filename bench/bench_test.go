package bench

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type alwaysGood struct{}

func (alwaysGood) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"event_id":"abc","event_type":"BOOM"}`)),
	}, nil
}

func BenchmarkGet(b *testing.B) {
	ctx := context.Background()
	c := client{c: &http.Client{Transport: alwaysGood{}}}

	b.Run("with generic", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			if _, err := c.GetEventGeneric(ctx, 1000); err != nil {
				b.Errorf("unexpected error: %v", err)
			}
		}
	})
	b.Run("no generic", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			if _, err := c.GetEvent(ctx, 1000); err != nil {
				b.Errorf("unexpected error: %v", err)
			}
		}
	})
}

func BenchmarkCreate(b *testing.B) {
	ctx := context.Background()
	c := client{c: &http.Client{Transport: alwaysGood{}}}
	event := Event{
		EventID:   "abc",
		EventType: "BOOM",
	}

	b.Run("with generic", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			if _, err := c.CreateEventGeneric(ctx, event); err != nil {
				b.Errorf("unexpected error: %v", err)
			}
		}
	})
	b.Run("no generic", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			if _, err := c.CreateEvent(ctx, event); err != nil {
				b.Errorf("unexpected error: %v", err)
			}
		}
	})
}
