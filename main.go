package main

import (
	"fmt"
	"trees/common"
	"trees/history"
	"trees/util"
)

func main() {
	tree := history.NewHistoryTree(new(common.XorHasher), common.NewInMemoryStore())

	for i := uint64(0); i < 10; i++ {
		eventDigest := util.Uint64AsBytes(i)
		commitment := tree.Add(eventDigest, i)
		fmt.Println(commitment)
	}

}
