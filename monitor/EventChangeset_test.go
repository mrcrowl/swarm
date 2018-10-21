package monitor

import (
	"testing"

	"github.com/rjeczalik/notify"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	sut := NewEventChangeset()
	success := sut.Add(notify.Create, "abcd/efgh.js")
	assert.True(t, success)
}

func TestAddComposite(t *testing.T) {
	sut := NewEventChangeset()
	success := sut.Add(notify.Create|notify.Remove, "abcd/efgh.js")
	assert.False(t, success)
}

func TestAddDuplicate(t *testing.T) {
	sut := NewEventChangeset()
	sut.Add(notify.Create, "abcd/efgh.js")
	success := sut.Add(notify.Create, "abcd/efgh.js")
	assert.False(t, success)
}

func TestMakeEventKey(t *testing.T) {
	key := makeEventKey(notify.Create, "abcd")
	assert.Equal(t, "C:abcd", key)
}

var expectedStrings = map[notify.Event]string{
	notify.Create:              "C",
	notify.Write:               "W",
	notify.Remove:              "D",
	notify.Rename:              "M",
	notify.Create | notify.All: "?",
}

func TestEventToString(t *testing.T) {
	for e, str := range expectedStrings {
		eventString := eventToString(e)
		assert.Equal(t, str, eventString)
	}
}

func TestCompositeEvent(t *testing.T) {
	assert.False(t, isCompositeEvent(notify.Create))
	assert.True(t, isCompositeEvent(notify.Create|notify.Write))
	assert.False(t, isCompositeEvent(1))
	assert.False(t, isCompositeEvent(0))
	assert.False(t, isCompositeEvent(2))
	assert.False(t, isCompositeEvent(2048))
	assert.True(t, isCompositeEvent(3))
	assert.True(t, isCompositeEvent(17))
	assert.True(t, isCompositeEvent(2052))
}
