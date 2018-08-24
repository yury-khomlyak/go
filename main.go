package main

import (
	"math/rand"
	"time"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
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

////////// DO NOT ALTER AFTER THIS LINE /////////

type Action int

const (
	Nothing Action = iota
	Move
	Load
	Unload
	Eat
)

type Direction int

const (
	Up Direction = iota
	UpRight
	Right
	DownRight
	Down
	DownLeftf
	Left
	UpLeft
)

type Hive struct {
	Id   string
	Skin uint8
	Ox   uint8
	Oy   uint8
	Ants map[int]*Ant
	Map  *Map
}

type Map struct {
	Width     uint8
	Height    uint8
	HiveLimit uint8
	Cells     [][]*Cell
}

type Cell struct {
	Food uint8
	CellType
}

type CellType uint8

const (
	Dark CellType = iota
	Grass
	HillBump
	HillConer
	HillHalf
	HillThird
)

type Ant struct {
	Wasted  int
	Age     uint8
	Health  uint8
	Payload uint8
	X       uint8
	Y       uint8
	Event
}

type Event uint8

const (
	Good Event = iota
	Birth
	NoAction
	Slow
	BadMove
	BadLoad
	BadUnload
	BadEat
	Collision
	Error
	Death
)

type ActDir map[int]int

func StartServer(){
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 501)
		return
	}

	if len(data) == 0 {
		http.Error(w, err.Error(), 502)
		return
	}

	var hive Hive
	err = json.Unmarshal(data, &hive)
	if err != nil {
		fmt.Println("Fail to convrt json to object", err)
		http.Error(w, err.Error(), 503)
		return
	}

	actions := whatToDo(&hive)
	fmt.Println(actions)

	output, err := json.Marshal(actions)
	if err != nil {
		http.Error(w, err.Error(), 504)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}