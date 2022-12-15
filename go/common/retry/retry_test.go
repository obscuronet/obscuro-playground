package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoWithTimeoutStrategy_SuccessAfterRetries(t *testing.T) {
	var count int
	testFunc := func() error {
		count = count + 1
		fmt.Printf("c: %d\n", count)
		if count < 3 {
			return fmt.Errorf("attempt number %d", count)
		}
		return nil
	}
	err := Do(testFunc, NewTimeoutStrategy(1*time.Second, 100*time.Millisecond))
	if err != nil {
		assert.Fail(t, "Expected function to succeed before timeout but failed", err)
	}

	assert.Equal(t, 3, count, "expected function to be called 3 times before succeeding")
}

func TestDoWithTimeoutStrategy_UnsuccessfulAfterTimeout(t *testing.T) {
	var count int
	testFunc := func() error {
		count = count + 1
		fmt.Printf("c: %d\n", count)
		return fmt.Errorf("attempt number %d", count)
	}
	err := Do(testFunc, NewTimeoutStrategy(600*time.Millisecond, 100*time.Millisecond))
	if err == nil {
		assert.Fail(t, "expected failure from timeout but no err received")
	}

	assert.Greater(t, count, 5, "expected function to be called at least 5 times before timing out")
}

func TestDoWithFailFast_ReturnImmediately(t *testing.T) {
	// similar to unsuccessful test except our function returns a FailFast error and we expect no retries
	var count int
	testFunc := func() error {
		count = count + 1
		fmt.Printf("c: %d\n", count)
		return FailFast(fmt.Errorf("attempt number %d", count))
	}
	err := Do(testFunc, NewTimeoutStrategy(600*time.Millisecond, 100*time.Millisecond))
	if err == nil {
		assert.Fail(t, "expected failure from timeout but no err received")
	}

	assert.Equal(t, count, 1, "expected function to only be called once before breaking out of retry")
}

func TestDoublingBackoffStrategy_DoublingIntervalsAndRespectMaxRetries(t *testing.T) {
	var count int
	prevAttempt := time.Now()
	expectedInterval := 20
	testFunc := func() error {
		count = count + 1
		if count > 1 {
			// we check it waited at least long enough (not checking for too long as that could be flaky)
			assert.Greater(t, time.Since(prevAttempt), time.Duration(expectedInterval)*time.Millisecond)
			expectedInterval = 2 * expectedInterval
		}
		prevAttempt = time.Now()
		return fmt.Errorf("attempt number %d", count)
	}

	err := Do(testFunc, NewDoublingBackoffStrategy(20*time.Millisecond, 5))
	if err == nil {
		assert.Fail(t, "expected failure from hitting max retries but no err found")
	}

	assert.Equal(t, 5, count, "expected function to be called exactly 5 times before failing")
}
