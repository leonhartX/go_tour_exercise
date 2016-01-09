package main

import "golang.org/x/tour/tree"
import "fmt"

func Walk(t *tree.Tree, ch chan int) {
	defer close(ch)
	var walker func(t *tree.Tree)
	walker = func(t *tree.Tree) {
		if t == nil {
			return
		}
		walker(t.Left)
		ch <- t.Value
		walker(t.Right)
	}
	walker(t)
}

func Same(t1, t2 *tree.Tree) (result bool) {
	result = false
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2
		if ok1 != ok2 || v1 != v2 {
			return
		}
		if !ok1 {
			break
		}
	}
	result = true
	return
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
