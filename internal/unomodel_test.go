package internal

import (
	"testing"
)

func TestUnoGame(t *testing.T) {
	mygame := UnoGame{}
	mygame.Init()
	mygame.AddPlayer("test1", 10086)
	mygame.AddPlayer("test2", 10000)
	mygame.DealCard()
	t.Log(mygame)
}
