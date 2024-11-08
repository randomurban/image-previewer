package cache

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Parallel()
	t.Run("empty cache", func(t *testing.T) {
		t.Parallel()
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		t.Parallel()
		c := NewCache(5)

		_, wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		_, wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		_, wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCache2(t *testing.T) {
	t.Parallel()
	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		deletedVal, wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		require.Nil(t, deletedVal)

		deletedVal, wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		require.Nil(t, deletedVal)

		deletedVal, wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)
		require.Nil(t, deletedVal)

		deletedVal, wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)
		require.Equal(t, 100, deletedVal)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})

	t.Run("complex purge logic", func(t *testing.T) {
		t.Parallel()
		c := NewCache(3)

		deletedVal, wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		require.Nil(t, deletedVal)

		deletedVal, wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		require.Nil(t, deletedVal)

		deletedVal, wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)
		require.Nil(t, deletedVal)

		val, ok := c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		deletedVal, wasInCache = c.Set("aaa", 500)
		require.True(t, wasInCache)
		require.Nil(t, deletedVal)

		deletedVal, wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)
		require.Equal(t, 300, deletedVal)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 500, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})

	t.Run("clear cache", func(t *testing.T) {
		t.Parallel()
		c := NewCache(3)

		_, wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		_, wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		_, wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		_, wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		c.Clear()

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("ddd")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Parallel()
	c := NewCache(10)
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100_000; i++ {
			number, err := rand.Int(rand.Reader, big.NewInt(100_000))
			require.NoError(t, err)
			c.Get(Key(number.String()))
		}
	}()
	wg.Wait()

	val, ok := c.Get("99999")
	require.True(t, ok)
	require.Equal(t, 99999, val)

	val, ok = c.Get("99990")
	require.True(t, ok)
	require.Equal(t, 99990, val)

	val, ok = c.Get("0")
	require.False(t, ok)
	require.Nil(t, val)
}
