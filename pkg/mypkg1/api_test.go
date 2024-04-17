package mypkg1

import (
	"cd-slack-notification-bot/go/pkg/testutils"
	"testing"
)

func TestMyFuncNew(t *testing.T) {
	t.Run("Valid case", func(it *testing.T) {
		value, err := MyFunc(10)
		if err != nil {
			it.Errorf("Expected no error, got %v", err)
		}
		if value != 20 {
			it.Errorf("Expected 20, got %v", value)
		}
	})

	t.Run("Invalid case", func(it *testing.T) {
		_, err := MyFunc(1001)
		if err == nil {
			it.Errorf("Expected error, got no error: %v", err)
		}
	})

	t.Run("Integration test case", func(it *testing.T) {
		testutils.SkipCI(it)
		t.Log("Test 1 integration")
		value, err := MyFunc(20)
		if err != nil {
			it.Errorf("Expected no error, got %v", err)
		}
		if value != 40 {
			it.Errorf("Expected 40, got %v", value)
		}
	})
}
