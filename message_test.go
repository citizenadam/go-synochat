package synochat_test

import (
	"testing"

	"github.com/0ghny/go-synochat"
	"github.com/stretchr/testify/require"
)

const (
	sutBaseURL = "https://url"
	sutToken   = "atoken"
)

func TestSendMessageToChannel(t *testing.T) {
	t.Skip("This test requires a live server and token. Enable and configure as needed.")

	c, err := synochat.NewClient(sutBaseURL)
	require.NoError(t, err)
	require.NotNil(t, c)

	err = c.SendMessage(&synochat.ChatMessage{Text: "Hello from automated test"}, sutToken)
	require.NoError(t, err)
}
