package main

import "github.com/luuuweiii/RiceDouyin/controller"

func main() {
	handler, err := controller.BuildInjector()
	if err != nil {
		panic(err)
	}
	r := initRouter(handler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
