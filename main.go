package main

import (
	"app/game_mgr"
	"fmt"
)

func main() {
	fmt.Println("[Server Start]")
	mgr := game_mgr.NewGameMgr()
	err := mgr.Run()
	if err != nil {
		panic(err)
	}
}
