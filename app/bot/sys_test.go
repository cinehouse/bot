package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSys_OnMessage(t *testing.T) {
	bot, err := NewSys("./../../data")
	require.NoError(t, err)
	assert.Equal(t, Response{Text: "_pong_", Send: true}, bot.OnMessage(Message{Text: "ping"}))
}

func TestSys_Help(t *testing.T) {
	bot, err := NewSys("./../../data")
	require.NoError(t, err)
	assert.Equal(t, "ping _– answer pong_\n", bot.Help())
}

func TestSys_Failed(t *testing.T) {
	_, err := NewSys("/tmp/no-such-place")
	require.Error(t, err)
}
