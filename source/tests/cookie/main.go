//go:build wasm

package main

import (
	"fmt"
	"gooey/cookie"
	"time"
)

func main() {
	cookieStore := cookie.CookieStore
	fmt.Println("wasm loaded")
	err := cookieStore.Set(cookie.SetOptions{
		Name:  "hello",
		Value: "world",
	})
	if err != nil {
		panic(err)
	}

	_, err = cookieStore.Get(cookie.GetOptions{Name: "hello"})
	if err != nil {
		panic(err)
	}

	_, err = cookieStore.GetAll(nil)
	if err != nil {
		panic(err)
	}

	err = cookieStore.Delete(cookie.DeleteOptions{
		Name: "hello",
	})
	if err != nil {
		panic(err)
	}

	for true {
		// Do Nothing
		time.Sleep(100 * time.Millisecond)
	}

}
