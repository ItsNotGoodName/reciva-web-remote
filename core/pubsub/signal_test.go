package pubsub

import (
	"testing"
)

func TestSignalPub(t *testing.T) {
	pub := NewSignalPub()
	sub, unsub := pub.Subscribe()
	defer unsub()

	select {
	case <-sub:
		t.Error("sub should not have received a signal")
	default:
	}

	pub.Publish()

	select {
	case <-sub:
	default:
		t.Error("sub should have received a signal")
	}

	pub.Publish()
	pub.Publish()

	select {
	case <-sub:
	default:
		t.Error("sub should have received a signal")
	}

	select {
	case <-sub:
		t.Error("sub should not have received a signal")
	default:
	}

	unsub()
	pub.Publish()

	select {
	case <-sub:
		t.Error("sub should not have received a signal")
	default:
	}
}
