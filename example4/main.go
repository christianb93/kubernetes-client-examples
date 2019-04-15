package main

import (
	"errors"
	"fmt"

	"k8s.io/client-go/tools/cache"
)

type Book struct {
	Author string
	Title  string
}

func main() {
	fmt.Println("Creating store")
	// Create a store
	keyFunc := func(obj interface{}) (string, error) {
		book, ok := obj.(Book)
		if !ok {
			return "", errors.New("Not a book")
		}
		return book.Author + "/" + book.Title, nil
	}
	myStore := cache.NewStore(keyFunc)
	fmt.Println("Adding two books to store")
	mobyDick := Book{Author: "Melville", Title: "Moby Dick"}
	davidCopperfield := Book{Author: "Dickens", Title: "David Copperfield"}
	myStore.Add(mobyDick)
	myStore.Add(davidCopperfield)
	// Get a book
	fmt.Println("Get book again")
	item, exists, err := myStore.Get(mobyDick)
	if !exists || (err != nil) {
		fmt.Println("Hmm...this should exist!")
		panic("Could not get book")
	}
	fmt.Printf("Got book: %s\n", item.(Book))
	// Get book by key
	key, _ := keyFunc(mobyDick)
	fmt.Printf("Getting book by key, key is %s\n", key)
	item, exists, err = myStore.GetByKey(key)
	if !exists || (err != nil) {
		fmt.Println("Hmm...this should exist!")
		panic("Could not get book")
	}
	fmt.Printf("Got book: %s\n", item.(Book))
}
