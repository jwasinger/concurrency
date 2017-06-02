package main

import (
  "sync"
  "bytes"
  "fmt"
)

type ThreadState struct {
  m sync.Mutex
  state map[uint64]string
}

func (t *ThreadState) State(tid uint64) string {
  t.m.Lock()
  defer t.m.Unlock()

  return t.state[tid]
}

func sortUint64(l []uint64) {
  for i := 0; i < len(l); i++ {
    for j := i+1; j < len(l); j++ {
      if l[i] > l[j] {
        tmp := l[i]
        l[i] = l[j]
        l[j] = tmp
      }
    }
  }
}

func (t *ThreadState) SetState(tid uint64, val string) {
  t.m.Lock()
  defer t.m.Unlock()

  t.state[tid] = val
}

func (t *ThreadState) Print() {
  t.m.Lock()
  defer t.m.Unlock()

  strBuf := &bytes.Buffer{}

  keys := make([]uint64, len(t.state))

  i := 0
  for k := range t.state {
      keys[i] = k
      i++
  }

  sortUint64(keys)

  for _, v := range keys {
    s := fmt.Sprintf("%d:%s ", v, t.state[v])
    strBuf.WriteString(s)
  }

  strBuf.WriteString("\n")
  
  /*
  resetBuf := &bytes.Buffer{}
  for i := 0; i < len(strBuf.String()); i++ {
    resetBuf.WriteString("\b")
  }

  resetBuf.WriteString("\b")
  */

  fmt.Printf(strBuf.String())
  //fmt.Printf(resetBuf.String())
}

func (t *ThreadState) AddKey(val uint64) {
  t.m.Lock()
  defer t.m.Unlock()

  t.state[val] = "s"
}

func NewThreadState(num uint64) *ThreadState {
  return &ThreadState{m: sync.Mutex{}, state: make(map[uint64]string)}
}
