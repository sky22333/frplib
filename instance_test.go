package frplib

import (
	"context"
	"strings"
	"testing"
)

type fakeService struct {
	closed bool
}

func (f *fakeService) Close() error {
	f.closed = true
	return nil
}

func TestRepeatedStartReturnsAlreadyRunning(t *testing.T) {
	m := newManager()
	factory := func(context.Context, string) (closeableService, error) {
		return &fakeService{}, nil
	}

	if err := m.start(instanceTypeClient, "a", "x=1", factory); err != nil {
		t.Fatalf("first start failed: %v", err)
	}
	if err := m.start(instanceTypeClient, "a", "x=1", factory); err == nil {
		t.Fatalf("second start should fail")
	} else if got := err.Error(); !strings.HasPrefix(got, ErrAlreadyRunning) {
		t.Fatalf("expected %s, got %q", ErrAlreadyRunning, got)
	}
}

func TestRepeatedStopIsNoop(t *testing.T) {
	m := newManager()
	factory := func(context.Context, string) (closeableService, error) {
		return &fakeService{}, nil
	}

	if err := m.start(instanceTypeClient, "a", "x=1", factory); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if err := m.stop(instanceTypeClient, "a"); err != nil {
		t.Fatalf("first stop failed: %v", err)
	}
	if err := m.stop(instanceTypeClient, "a"); err != nil {
		t.Fatalf("second stop should be no-op: %v", err)
	}
}

func TestInvalidID(t *testing.T) {
	if err := validateID("../bad"); err == nil {
		t.Fatalf("invalid id should fail")
	}
}
