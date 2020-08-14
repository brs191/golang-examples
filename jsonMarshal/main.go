package main

import (
	"encoding/json"
	"fmt"
)

type fruit struct {
	Apple  string `json:"Apple"`
	Around string `json:"Around"`
}

// Alpha for apple
type Alpha struct {
	A      fruit
	Banana string `json:"Banana"`
}

func main() {

	fmt.Println("Hi ")

	msg := Alpha{
		A: fruit{
			Apple:  "apple",
			Around: "around",
		},
		Banana: "Banana",
	}
	byteArray, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(string(byteArray))

	newA := fruit{
		Apple:  "apple",
		Around: "around",
	}
	newmsg := Alpha{
		A:      newA,
		Banana: "Banana",
	}

	mA, _ := json.Marshal(newA)

	mB, _ := json.Marshal(newmsg)

	fmt.Println(string(mA) + string(mB))
}
