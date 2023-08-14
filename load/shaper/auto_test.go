package shaper

import (
	"math"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestAutoShaper_GrowsAdditive(t *testing.T) {
	ctrl := gomock.NewController(t)
	info := NewMockLoadInfoSource(ctrl)

	info.EXPECT().GetSentTransactions().AnyTimes().Return(uint64(120), nil)
	info.EXPECT().GetReceivedTransactions().AnyTimes().Return(uint64(120), nil)

	shaper := NewAutoShaper(10, 0.2)

	start := time.Now()
	shaper.Start(start, info)

	for i := 0; i < 10; i++ {
		want := float64(i * 10)
		if got := shaper.GetNumMessagesInInterval(start, time.Second); math.Abs(got-want) > 1e-6 {
			t.Errorf("invalid number of messages, wanted %f, got %f", want, got)
		}
		start = start.Add(time.Second)
	}
}

func TestAutoShaper_ShrinksMultiplicative(t *testing.T) {
	ctrl := gomock.NewController(t)
	info := NewMockLoadInfoSource(ctrl)

	info.EXPECT().GetSentTransactions().AnyTimes().Return(uint64(100000), nil)
	info.EXPECT().GetReceivedTransactions().AnyTimes().Return(uint64(0), nil)

	rate := 1000.0
	shaper := NewAutoShaper(10, 0.2)
	shaper.(*autoShaper).rate = rate

	start := time.Now()
	shaper.Start(start, info)

	for i := 0; i < 10; i++ {
		want := float64(rate)
		if got := shaper.GetNumMessagesInInterval(start, time.Second); math.Abs(got-want) > 1e-6 {
			t.Errorf("invalid number of messages, wanted %f, got %f", want, got)
		}
		rate *= 0.8
		start = start.Add(time.Second)
	}
}
