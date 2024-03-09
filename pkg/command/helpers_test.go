package command

import (
	"testing"
)

func TestSelectIdFromTextError(t *testing.T) {
	id, err := selectIdFromText("azaza")

	if err == nil {
		t.Error("expected error got nil")
	}

	if id != "" {
		t.Error("expected empty id got", id)
	}
}

func TestSelectIdFromTextSuccess(t *testing.T) {
	id, err := selectIdFromText("1. 82341644-7772-4745-a75f-ae08957d9008 (PC)")
	if err != nil {
		t.Error("expected nil error got", err)
	}

	if id != "82341644-7772-4745-a75f-ae08957d9008" {
		t.Error("expected 82341644-7772-4745-a75f-ae08957d9008 id got", id)
	}
}
