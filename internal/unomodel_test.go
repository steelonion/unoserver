package internal

import (
	"testing"
)

func TestUnoGame(t *testing.T) {
	mygame := UnoGame{}
	mygame.Init()
	t.Log(mygame)
}
