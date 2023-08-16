package main

import (
	"log"

	"github.com/horockey/go-toolbox/datastructs/avl_tree"
)

func main() {
	t := avl_tree.New[int]()

	err := t.Insert("k5", 1)
	fatalOnErr(err)
	err = t.Insert("k2", 2)
	fatalOnErr(err)
	err = t.Insert("k10", 10)
	fatalOnErr(err)
	err = t.Insert("k6", 6)
	fatalOnErr(err)
	err = t.Insert("k4", 4)
	fatalOnErr(err)
	err = t.Insert("k11", 4)
	fatalOnErr(err)
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
