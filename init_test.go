package ini

import (
	"errors"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (c *mockClient) Test() {
	c.Called()
}

func TestOnce(t *testing.T) {
	var init Once

	client := &mockClient{}
	client.On("Test")

	err1 := init.Do(func() error {
		client.Test()
		return errors.New("error 1")
	})

	err2 := init.Do(func() error {
		client.Test()
		return errors.New("error 2")
	})

	client.AssertNumberOfCalls(t, "Test", 1)
	assert.Same(t, err1, err2)
	assert.Error(t, err1, "error 1")
}

func TestOnce_NoError(t *testing.T) {
	var init Once

	client := &mockClient{}
	client.On("Test")

	err1 := init.Do(func() error {
		client.Test()
		return nil
	})

	err2 := init.Do(func() error {
		client.Test()
		return errors.New("error 2")
	})

	client.AssertNumberOfCalls(t, "Test", 1)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestOnce_Concurrency(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var init Once

		client := &mockClient{}
		client.On("Test")

		var (
			err1, err2 error
			called     int
		)
		go func() {
			called++
			err1 = init.Do(func() error {
				client.Test()
				return errors.New("error 1")
			})
		}()
		go func() {
			called++
			err2 = init.Do(func() error {
				client.Test()
				return errors.New("error 2")
			})
		}()

		synctest.Wait()
		client.AssertNumberOfCalls(t, "Test", 1)
		assert.Equal(t, 2, called)

		// we know that err1 and err2 must be the same error, but we don't know if it's "error 1" or "error 2".
		assert.Same(t, err1, err2)
	})
}

func TestSuccessOnce(t *testing.T) {
	var init SuccessOnce

	client := &mockClient{}
	client.On("Test")

	err1 := init.Do(func() error {
		client.Test()
		return errors.New("error 1")
	})

	err2 := init.Do(func() error {
		client.Test()
		return nil
	})

	client.AssertNumberOfCalls(t, "Test", 2)
	assert.Error(t, err1, "error 1")
	assert.NoError(t, err2)
}

func TestSuccessOnce_Concurrency(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var init SuccessOnce

		client := &mockClient{}
		client.On("Test")

		var (
			called int
		)
		go func() {
			called++
			_ = init.Do(func() error {
				client.Test()
				return nil
			})
		}()
		go func() {
			called++
			_ = init.Do(func() error {
				client.Test()
				return errors.New("error 2")
			})
		}()

		synctest.Wait()

		// due to concurrency nature, we don't know how many times client.Test was called.
		// client.AssertNumberOfCalls(t, "Test", 1)

		// we do know init.Do was called twice.
		assert.Equal(t, 2, called)

		// best we can assert that the third init.Do will return nil and will not execution f.
		assert.NoError(t, init.Do(func() error { panic("don't call me maybe") }))
	})
}
