package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter_withRoot(t *testing.T) {
	var r Router[string]

	r.Add("/", "root")
	r.Add("/foo", "foo")
	r.Add("/bar", "bar")
	r.Add("/foo/bar", "foobar")

	t.Run("root-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"root"})
		assert.Equal(t, sub, "")
	})

	t.Run("root-sub", func(t *testing.T) {
		h, sub, ok := r.Match("/hello/world")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"root"})
		assert.Equal(t, sub, "hello/world")
	})

	t.Run("foo-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/foo")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"foo"})
		assert.Equal(t, sub, "")
	})

	t.Run("foo-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/foo/hello/world")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"foo"})
		assert.Equal(t, sub, "hello/world")
	})

	t.Run("bar-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/bar")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"bar"})
		assert.Equal(t, sub, "")
	})

	t.Run("bar-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/bar/hello/world")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"bar"})
		assert.Equal(t, sub, "hello/world")
	})

	t.Run("foobar-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/foo/bar")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"foobar"})
		assert.Equal(t, sub, "")
	})

	t.Run("foobar-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/foo/bar/hello/world")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"foobar"})
		assert.Equal(t, sub, "hello/world")
	})
}

func TestRouter_withoutRoot(t *testing.T) {
	var r Router[string]

	r.Add("/foo", "foo")
	r.Add("/bar", "bar")
	r.Add("/foo/bar", "foobar")

	t.Run("root-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/")
		assert.False(t, ok)
		assert.Equal(t, h, []string(nil))
		assert.Equal(t, sub, "")
	})

	t.Run("root-sub", func(t *testing.T) {
		h, sub, ok := r.Match("/hello/world")
		assert.False(t, ok)
		assert.Equal(t, h, []string(nil))
		assert.Equal(t, sub, "")
	})

	t.Run("foo-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/foo")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"foo"})
		assert.Equal(t, sub, "")
	})

	t.Run("foo-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/foo/hello/world")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"foo"})
		assert.Equal(t, sub, "hello/world")
	})

	t.Run("bar-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/bar")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"bar"})
		assert.Equal(t, sub, "")
	})

	t.Run("bar-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/bar/hello/world")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"bar"})
		assert.Equal(t, sub, "hello/world")
	})

	t.Run("foobar-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/foo/bar")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"foobar"})
		assert.Equal(t, sub, "")
	})

	t.Run("foobar-direct", func(t *testing.T) {
		h, sub, ok := r.Match("/foo/bar/hello/world")
		assert.True(t, ok)
		assert.Equal(t, h, []string{"foobar"})
		assert.Equal(t, sub, "hello/world")
	})
}
