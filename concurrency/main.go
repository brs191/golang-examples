package main

import (
	"sync"
)

func waitGroup() {
	println("waitGroupWithChannel enter")
	var wg sync.WaitGroup

	nums := [...]string{"one", "two", "three"}

	for _, num := range nums {
		wg.Add(1)
		go func(num string) {
			defer wg.Done()
			println(num)
		}(num)
	}

	wg.Wait()
	println("waitGroupWithChannel exit")

}

func waitGroupWithChannel() {
	println("waitGroupWithChannel enter")
	var wg sync.WaitGroup

	nums := [...]string{"one", "two", "three"}
	resCh := make(chan string)

	for _, num := range nums {
		wg.Add(1)
		go func(num string, resch chan string) {
			println(num)
			resCh <- num
			wg.Done()
		}(num, resCh)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	for n := range resCh {
		println("final ", n)
	}
	println("waitGroupWithChannel exit")
}

func waitGroupWithChannels() {
	println("waitGroupWithChannel enter")
	var wg sync.WaitGroup

	nums := [...]string{"one", "two", "three"}
	resCh := make(chan string)

	for _, num := range nums {
		wg.Add(1)
		go func(num string, resch chan string) {
			println(num)
			resCh <- num
			wg.Done()
		}(num, resCh)
	}

	Nums := [...]string{"One", "Two", "Three"}
	resCH := make(chan string)
	for _, num := range Nums {
		wg.Add(1)
		go func(num string, resch chan string) {
			println(num)
			resCH <- num
			wg.Done()
		}(num, resCh)
	}

	go func() {
		wg.Wait()
		close(resCh)
		close(resCH)
	}()

	for {
		select {
		case num, ok := <-resCh:
			if !ok {
				return
			}
			println("final ", num)
		case Num, ok := <-resCH:
			if !ok {
				return
			}
			println("final ", Num)
		}
	}

	// println("waitGroupWithChannels exit")
}

func mergeChannels(ch ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	merge := func(c <-chan string) {
		defer wg.Done()
		for {
			select {
			case num, ok := <-c:
				if !ok {
					return
				}
				out <- num
			}
		}
	}

	wg.Add(len(ch))
	for _, c := range ch {
		go merge(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func waitGroupMergeChannels() {
	println("waitGroupWithChannel enter")
	var wg sync.WaitGroup

	nums := [...]string{"one", "two", "three"}
	resCh := make(chan string)

	for _, num := range nums {
		wg.Add(1)
		go func(num string, resch chan string) {
			println("resch ", num)
			resCh <- num
			wg.Done()
		}(num, resCh)
	}

	Nums := [...]string{"One", "Two", "Three"}
	resCH := make(chan string)
	for _, num := range Nums {
		wg.Add(1)
		go func(num string, resch chan string) {
			println("resCH ", num)
			resCH <- num
			wg.Done()
		}(num, resCh)
	}

	go func() {
		wg.Wait()
		close(resCh)
		close(resCH)
	}()

	mergedCh := mergeChannels(resCh, resCH)
	for n := range mergedCh {
		println("merged: ", n)
	}

	// println("waitGroupWithChannels exit")
}

func main() {

	//waitGroup()
	//waitGroupWithChannels()
	waitGroupMergeChannels()

	println("main done!!")
}
