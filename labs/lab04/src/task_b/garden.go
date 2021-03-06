package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	MaxActions = 20
)

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func gardener(garden [][]bool, rwMutex *sync.RWMutex, exit chan int) {
	for i := 0; i < MaxActions; i++ {
		rwMutex.Lock()
		fmt.Println("Gardener watering plants")
		for i := 0; i < len(garden); i++ {
			for j := 0; j < len(garden[0]); j++ {
				if garden[i][j] == false {
					garden[i][j] = true
				}
			}
		}
		rwMutex.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
	exit <- 1
}

func nature(garden [][]bool, rwMutex *sync.RWMutex, exit chan int) {
	for i := 0; i < MaxActions; i++ {
		rwMutex.Lock()
		fmt.Println("Some plants dying away")
		for i := 0; i < len(garden)*2; i++ {
			index1 := random.Intn(len(garden))
			index2 := random.Intn(len(garden[0]))
			garden[index1][index2] = !garden[index1][index2]
		}
		rwMutex.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
	exit <- 1
}

func main() {
	var garden [][]bool
	var rwMutex sync.RWMutex
	exit := make(chan int, 4)

	for i := 0; i < 10; i++ {
		var row []bool
		for j := 0; j < 10; j++ {
			random = rand.New(rand.NewSource(time.Now().UnixNano()))
			row = append(row, random.Intn(2) != 0)
		}
		garden = append(garden, row)
	}

	go Monitor2(garden, &rwMutex, exit)
	go Monitor1(garden, &rwMutex, exit)
	go nature(garden, &rwMutex, exit)
	go gardener(garden, &rwMutex, exit)

	for i := 0; i < 4; i++ {
		<-exit
	}
}
