package main

import (
	"fmt"
	"mygithub/multipackage/bye"
	"mygithub/multipackage/hello"
)

func main() {
	fmt.Println("in main")
	hello.Sayhello()
	bye.Saybye()
	fmt.Println("exit main")
}
