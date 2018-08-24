package main

import (
	"math/rand"
	"time"
)

func main() {
	StartServer()
}

func whatToDo(hive *Hive) ActDir {
	actions := make(ActDir)
	rand.Seed(time.Now().UnixNano())
	for id, ant := range hive.Ants {

		//Default action if ant don't see Food
		action := Move

		//Default direction is Random
		direction := Direction(rand.Intn(4))

		food, hive, dir  := lookAround(ant, hive.Map)
		if hive && ant.Payload>0{
			direction = dir
			action = Unload
		}else if food{
			direction = dir
			if ant.Health<9{action = Eat}
			if ant.Payload<9 {action = Load}
		}

		actions[id] = int(action)*10 + int(direction)
	}

	return actions
}

func lookAround(ant *Ant, world *Map)(food, hive bool, dir Direction){

	if ant.Y > 0 {
		dir = Up
		food,hive = iSee(ant.Y-1,ant.X,world)
	}else if ant.Y < world.Height-1{
		dir = Down
		food,hive = iSee(ant.Y+1,ant.X,world)
	}else if ant.X < world.Width-1{
		dir = Right
		food,hive = iSee(ant.Y,ant.X+1,world)
	}else if ant.X >0{
		dir = Left
		food,hive = iSee(ant.Y,ant.X-1,world)
	}
	return
}

func iSee(y,x uint8, world *Map) (food, hive bool) {
	if world.Cells[y][x].Food>0{
		food = true
	}
	if world.Cells[y][x].CellType > 2{
		hive = true
	}
	return
}