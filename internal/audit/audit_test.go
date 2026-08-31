package audit

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestPublisherNotifyAllObservers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := Event{Timestamp: 1, Metrics: []string{"Alloc"}, IPAddress: "127.0.0.1"}

	first := NewMockObserver(ctrl)
	first.EXPECT().Notify(event).Return(nil)
	second := NewMockObserver(ctrl)
	second.EXPECT().Notify(event).Return(nil)

	p := NewPublisher(nil, first, second)
	p.Notify(event)
}

func TestPublisherNotifyOneFailingDoesNotStopOthers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := Event{Timestamp: 1, Metrics: []string{"Alloc"}, IPAddress: "127.0.0.1"}

	failing := NewMockObserver(ctrl)
	failing.EXPECT().Notify(event).Return(errors.New("boom"))
	ok := NewMockObserver(ctrl)
	ok.EXPECT().Notify(event).Return(nil)

	p := NewPublisher(nil, failing, ok)
	p.Notify(event)
}

func TestPublisherRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := Event{Timestamp: 1, Metrics: []string{"Alloc"}, IPAddress: "127.0.0.1"}

	observer := NewMockObserver(ctrl)
	observer.EXPECT().Notify(event).Return(nil)

	p := NewPublisher(nil)
	p.Register(observer)
	p.Notify(event)
}
