package main

import (
	"app/game_mgr"
	"fmt"
)

func main() {
	fmt.Println("server start")
	mgr := game_mgr.NewGameMgr()
	err := mgr.Run()
	if err != nil {
		panic(err)
	}
}
