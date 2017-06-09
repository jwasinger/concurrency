package main

import (
	"fmt"
	"math/rand"
	"time"
	"runtime"
)

type Resource struct {
	name string
}

type Table struct {
	resources [3]Resource
	smoked chan bool
	round chan bool
	activeResources [2]*Resource
}

func (t *Table) Smoked() chan bool {
	return t.smoked
}

func (t *Table) Get() [2]*Resource {
	return t.activeResources
}

func (t *Table) Round() chan bool {
	return t.smoked
}


func (t *Table) Place() {
	num1 := rand.Intn(3)
	num2 := 0
	for ; ; {
		num2 = rand.Intn(3)
		if num1 != num2 {
			break	
		}
	}

	t.activeResources[0] = &t.resources[num1]
	t.activeResources[1] = &t.resources[num2]
	
	fmt.Printf("%s and %s placed on table\n", *t.activeResources[0], *t.activeResources[1])
}

func NewTable() *Table {
	table := &Table{
		resources: [3]Resource{NewResource("tobacco"), NewResource("matches"), NewResource("papers")},
		smoked: make(chan bool),
		round: make(chan bool),
		activeResources: [2]*Resource{nil, nil},
	}

	table.activeResources = [2]*Resource{&table.resources[0], &table.resources[1]}
	return table
}

func NewResource(s string) Resource {
	return Resource{name: s}
}

func (r *Resource) Name() string {
	return r.name
}

func smoker(resource Resource, table *Table) {
	for ; ; {
		<- table.Round()
		resources := table.Get()
		grab := true

		for i := range resources {
			if resources[i].Name() == resource.Name() {
				grab = false
			}
		}
		
		if !grab {
			continue
		}

		fmt.Printf("goroutine %d is making cigarette with %s\n", getGID(), resource.Name())
		time.Sleep(1*time.Second)
		table.Smoked() <- true
	}
}

func agent(t *Table) {
	for ; ; {
		t.Place()
		t.Round() <- true
		t.Round() <- true
		t.Round() <- true
		<- t.Smoked()
		fmt.Println("agent sees that cigarette was smoked\n")
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	
	table := NewTable()
	go agent(table)

	go smoker(table.resources[0], table)
	go smoker(table.resources[1], table)
	go smoker(table.resources[2], table)

	time.Sleep(5*time.Minute)
}
