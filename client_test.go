package webclient_test

import (
	"net/http"
	"testing"

	"github.com/fernandoocampo/webclient"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultClient(t *testing.T) {
	expectedHttpClient := http.Client{}
	got := webclient.New()
	assert.NotNil(t, got)
	assert.Equal(t, &expectedHttpClient, got.HTTPClient())
}
