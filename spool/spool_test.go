package spool

import (
	"testing"

	"github.com/nubesk/binn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpool(t *testing.T) {
	s := New()
	s.Reset("123456")
	s.Publish("123456", &binn.Bottle{
		Msg: "sample message",
	})
	bs, err := s.Get("123456")
	require.NoError(t, err)
	assert.Equal(t, len(bs), 1)
	bs, err = s.Get("123456")
	require.NoError(t, err)
	assert.Equal(t, len(bs), 0)
}
