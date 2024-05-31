package main

import (
	"app/game_mgr"
)

func main() {
	mgr := game_mgr.NewGameMgr()
	err := mgr.Run()
	if err != nil {
		panic(err)
	}
}
