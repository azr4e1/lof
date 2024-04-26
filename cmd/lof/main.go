package main

import (
	"fmt"
	"os"

	"github.com/azr4e1/lof"
)

func main() {
	node, err := lof.GetTree()
	if err != nil {
		os.Exit(1)
	}

	data, err := lof.OutputWindows(node, func(bn *lof.BaseNode) bool {
		return bn.IsTrueWindow()
	}).ToJSON()

	fmt.Println(string(data))
}
