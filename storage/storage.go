package storage

import (
	"strings"
)

// AddressBook for cassini
type AddressBook interface {
	Add(address string) (bool, error)
	Exist(address string) (bool, error)
}

// AddressCacheBook implements a address memcache
type AddressCacheBook struct {
	addrs map[string]bool
}

// Add address to memcache
func (c *AddressCacheBook) Add(address string) (bool, error) {
	c.addrs[strings.ToUpper(address)] = true
	return true, nil
}

// Exist determine whether an address exists in memcache
func (c *AddressCacheBook) Exist(address string) (bool, error) {
	return c.addrs[strings.ToUpper(address)], nil
}

// NewAddressBook creates an AddressCacheBook
func NewAddressBook() AddressBook {
	book := &AddressCacheBook{
		addrs: make(map[string]bool)}
	return book
}
