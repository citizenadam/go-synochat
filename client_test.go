package synochat_test

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestNewClient_WithValidUrl_ShouldCreateNewClient(t *testing.T) {
	c, err := synochat.NewClient("http://syno.local")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "http://syno.local", c.BaseURL.String())
}

func TestNewClient_WithInvalidBaseURL_ShouldReturnsAnError(t *testing.T) {
	c, err := synochat.NewClient("syno.local")
	assert.Nil(t, c)
	assert.Error(t, err)
	assert.Equal(t, "invalid URL provided", err.Error())
}

func TestNewClient_WithEmptyBaseURL_ShouldReturnsAnError(t *testing.T) {
	c, err := synochat.NewClient("")
	assert.Nil(t, c)
	assert.Error(t, err)
	assert.Equal(t, "url cannot be empty or just whitespaces", err.Error())
}

func TestNewCustomClient_WithValidUrlAndNilHttpClient_ShouldCreateNewClient(t *testing.T) {
	c, err := synochat.NewCustomClient("http://syno.local", nil)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "http://syno.local", c.BaseURL.String())
}
