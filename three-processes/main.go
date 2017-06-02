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
    fmt.Println("blocking until resource free")
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


  s.SetState(getGID(), "w")
  time.Sleep(5*time.Second)

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
	rnd := rand.New(rand.NewSource(int64(i)))
    go func(r *Resource, s *ThreadState, rnd *rand.Rand) {
      s.AddKey(getGID())

      for ; ; {
        wt := rnd.Intn(10)
        time.Sleep(time.Duration(wt)*time.Second)
        r.Access(state)
      }
    } (&r, state, rnd)
  }

  for ; ; {
    state.Print()
    time.Sleep(1*time.Second)
  }

  wg := &sync.WaitGroup{}
  wg.Add(1)
  wg.Wait()
}
