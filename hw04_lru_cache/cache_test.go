package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func fillCache(c *Cache, capacity int) {
	for i := 0; i < capacity; i++ {
		(*c).Set(Key(strconv.Itoa(i)), i)
	}
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("exceed capacity", func(t *testing.T) {
		capacity := 3
		c := NewCache(capacity)
		fillCache(&c, capacity)
		c.Set("abc", 404)
		_, ok := c.Get("0")
		require.False(t, ok, "Exceeding cache capacity didn't push last element from cache")
	})

	t.Run("test value changed", func(t *testing.T) {
		c := NewCache(3)
		oldVal := 404
		expectedVal := 777
		key := Key("abc")
		c.Set(key, oldVal)
		c.Set(key, expectedVal)
		value, ok := c.Get(key)
		require.True(t, ok, "Invalid cache key presence")
		require.Equal(t, expectedVal, value, "Invalid value from cache")
	})

	t.Run("test value changed", func(t *testing.T) {
		c := NewCache(3)
		oldVal := 404
		expectedVal := 777
		key := Key("abc")
		c.Set(key, oldVal)
		c.Set(key, expectedVal)
		value, ok := c.Get(key)
		require.True(t, ok)
		require.Equal(t, expectedVal, value)
	})

	t.Run("test clear", func(t *testing.T) {
		capacity := 4
		key := Key("0")
		expectedVal := 0
		c := NewCache(capacity)
		fillCache(&c, capacity)
		val, ok := c.Get(key)
		require.Equal(t, expectedVal, val)
		require.True(t, ok)
		c.Clear()
		val, ok = c.Get(key)
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("test list works after clear", func(t *testing.T) {
		capacity := 3
		key := Key("2")
		backKey := Key("0")
		expectedKey := Key("abc")
		expectedVal := "Unusual val"
		c := NewCache(capacity)
		fillCache(&c, capacity)
		c.Clear()
		fillCache(&c, capacity)
		_, ok := c.Get(key)
		require.True(t, ok)
		c.Set(expectedKey, expectedVal)
		_, ok = c.Get(backKey)
		require.False(t, ok)
		val, _ := c.Get(expectedKey)
		require.Equal(t, expectedVal, val)
	})

	t.Run("elements are properly pushed front", func(t *testing.T) {
		capacity := 3
		keys := []Key{"0", "1"}
		pushingKey := Key("abc")
		oldestKey := Key("2")
		c := NewCache(capacity)
		fillCache(&c, capacity)
		for _, key := range keys {
			c.Get(key)
		}
		c.Set(pushingKey, "")
		_, ok := c.Get(oldestKey)
		require.False(t, ok, "Oldest element is still present in cache")
		for _, key := range keys {
			c.Get(key)
		}
		c.Set(oldestKey, "")
		for _, key := range keys {
			c.Set(key, "")
		}
		c.Set(pushingKey, "")
		_, ok = c.Get(oldestKey)
		require.False(t, ok, "Oldest element is still present in cache")
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
