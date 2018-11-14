package common

import (
	"testing"

	//catypes "github.com/QOSGroup/cassini/types"

	"github.com/stretchr/testify/assert"
)

func TestParseUrls(t *testing.T) {

	f, s, err := ParseUrls("http://127.0.0.1:2379", "")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(f))
	assert.Equal(t, "127.0.0.1", f[0].Hostname())
	assert.Equal(t, "2379", s[0].Port())

	f, s, err = ParseUrls("http://192.168.1.100:8080", "http://localhost/,http://127.0.0.1:9001")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(f))
	assert.Equal(t, "192.168.1.100", f[0].Hostname())
	assert.Equal(t, "8080", f[0].Port())

	assert.Equal(t, 2, len(s))
	assert.Equal(t, "localhost", s[0].Hostname())
	assert.Equal(t, "", s[0].Port())
	assert.Equal(t, "127.0.0.1", s[1].Hostname())
	assert.Equal(t, "9001", s[1].Port())

}
