package main

import "object-storage/internal/app"

func main() {
	a, cleanup, err := app.Init()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	a.Run()
}
