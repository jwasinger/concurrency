package main

import (
  "fmt"
  "time"
  "sync"
)

func main() {
  l := &LinkedList{}

  insert_wg := &sync.WaitGroup{}
  delete_wg := &sync.WaitGroup{}
  search_wg := &sync.WaitGroup{}

  im := &sync.Mutex{}
  dm := &sync.Mutex{}

  numSearchers := 2
  numDeleters := 2
  numInserters := 2

  //inserters
  for i := 0; i < numInserters; i++ {
    go func() {
      for ; ; {
        //dont run while deleters are running
        delete_wg.Wait()

        insert_wg.Add(1)

        fmt.Println("inserting...")
        time.Sleep(1*time.Second)

        im.Lock()
        l.PushEnd(1)
        fmt.Println(l.Size())
        im.Unlock()

        insert_wg.Done()
        time.Sleep(1*time.Second)
      }
    } ()
  }

  //deleters
  for i := 0; i < numDeleters; i++ {
    go func() {
      for ; ; {
        insert_wg.Wait()
        search_wg.Wait()
        
        delete_wg.Add(1)

        fmt.Println("deleting....")
        time.Sleep(time.Second*1)

        dm.Lock()
        l.Delete(l.Size()-1)
        fmt.Println(l.Size())
        dm.Unlock()

        delete_wg.Done()
        time.Sleep(1*time.Second)
      }
    } ()
  }

  for i := 0; i < numSearchers; i++ {
    go func() {
      for ; ; {
        delete_wg.Wait()
        search_wg.Add(1)

        fmt.Println("searching... ")
        time.Sleep(time.Second*1)

        search_wg.Done()
        time.Sleep(1*time.Second)
      }
    } ()
  }

  time.Sleep(5*time.Minute)
}

