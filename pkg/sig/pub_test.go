package sig

import (
	"testing"
)

func TestPub(t *testing.T) {
	pub := NewPub()
	sub := pub.Subscribe()

	select {
	case <-sub.Channel():
		t.Error("sub.Channel() should not have received a signal")
	default:
	}

	pub.Publish()

	select {
	case <-sub.Channel():
	default:
		t.Error("sub.Channel() should have received a signal")
	}

	pub.Publish()
	pub.Publish()

	select {
	case <-sub.Channel():
	default:
		t.Error("sub.Channel() should have received a signal")
	}

	select {
	case <-sub.Channel():
		t.Error("sub.Channel() should not have received a signal")
	default:
	}

	pub.Unsubscribe(sub)
	pub.Publish()

	select {
	case <-sub.Channel():
		t.Error("sub.Channel() should not have received a signal")
	default:
	}
}
