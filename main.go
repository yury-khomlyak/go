package main

import (
	"math/rand"
	"time"
)

func whatToDo(hive *Hive) ActDir {
	actions := make(ActDir)
	rand.Seed(time.Now().UnixNano())

	for id := range hive.Ants {

		//Only move or Eat randomly
		act := Action(0)
		switch rand.Intn(2) {
		case 0:
			act = Move
		case 1:
			act = Eat
		}

		//                       Random direction
		actions[id] = int(act) + rand.Intn(4)
	}

	time.Sleep(20 * time.Millisecond)

	return actions
}

