package main

import (
  "sync"
  "time"
  "runtime"
  "math/rand"
  "fmt"
)

type Resource struct {
  numWorking int
  mutex1 sync.Mutex
  mutex2 sync.Mutex
}

type State struct {
  s map[int]sync.Mutex
}

func (r *Resource) Access(s *ThreadState) {
  r.mutex1.Lock()

  if r.numWorking >= 3 {
    //block until all are complete
    for ;  ; {
      time.Sleep(1*time.Millisecond)
      if r.numWorking == 0 {
        fmt.Println("all done")
        time.Sleep(1*time.Second)
        break
      }
    }
  }

  r.numWorking++

  r.mutex1.Unlock()

  // STATE working

  // do work.. wait a random amount of time

  workTime := rand.Intn(10-2)+2

  s.SetState(getGID(), "w")
  time.Sleep(time.Duration(workTime)*time.Second)

  r.mutex2.Lock()
  r.numWorking--
  r.mutex2.Unlock()

  s.SetState(getGID(), "s")
}

func main() {
  runtime.GOMAXPROCS(4)
  r := Resource{numWorking:0, mutex1:sync.Mutex{}, mutex2:sync.Mutex{}}

  state := NewThreadState(4)

  for i := 0; i < 5; i++ {
    go func(r *Resource, s *ThreadState) {
      s.AddKey(getGID())

      for ; ; {
        r.Access(state)
      }
    } (&r, state)
  }

  for ; ; {
    state.Print()
    time.Sleep(1*time.Second)
  }

  wg := &sync.WaitGroup{}
  wg.Add(1)
  wg.Wait()
}
