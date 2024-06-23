package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	env, err := NewEnv()
	if err != nil {
		panic(err)
	}

	fmt.Println("Proxy running at", env.Addr())

	proxy := NewProxy(env)

	proxy.Run()
}
