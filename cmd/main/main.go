package main

import "app/service/main/game_mgr"

func main() {
	go func() {
		err := game_mgr.RunHttp()
		if err != nil {
			panic(err)
		}
	}()

	err := game_mgr.NewGameMgr().Run()
	if err != nil {
		panic(err)
	}
}
