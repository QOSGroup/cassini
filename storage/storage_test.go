package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressBook(t *testing.T) {

	book := NewAddressBook()
	ok, err := book.Add("asd")
	assert.NoError(t, err)
	assert.Equal(t, ok, true)

	ok, err = book.Exist("123")
	assert.NoError(t, err)
	assert.Equal(t, ok, false)

	ok, err = book.Exist("aSd")
	assert.NoError(t, err)
	assert.Equal(t, ok, true)
}
