package gnet

import "testing"

func TestServer(t *testing.T) {
	s := NewServer("ginx")
	s.Server()
}
