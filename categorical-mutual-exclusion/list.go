package main

import (
  "fmt"
  "errors"
)

type LLNode struct {
  Value int
  Next *LLNode 
}

type LinkedList struct {
  Root *LLNode  
}

func (l *LinkedList) Print() {
  fmt.Printf("[")

  if l.Root != nil {
    n := l.Root
    fmt.Printf("%d, ", n.Value)
    n = n.Next

    for ; n != nil; n = n.Next {
      fmt.Printf("%d, ", n.Value)
    }
  }

  fmt.Printf("]\n")
}

func (l *LinkedList) PushEnd(val int) {
  var prev *LLNode = nil

  if l.Root == nil {
    l.Root = &LLNode{val, nil}
    return
  }

  for n := l.Root; n != nil; n = n.Next { 
    prev = n
  }

  prev.Next = &LLNode{val, nil}  
}

/*
 * Find the location of an element within the list
 */
func (l *LinkedList) Search(val int) int {
  if l.Root == nil {
    return -1
  }

  n := l.Root

  for i := 0; n != nil; i++ {
    if n.Value == val {
      return i
    }

    n = n.Next
  }
  return -1
}

/*
 * Delete an element at location i
 */
func (l *LinkedList) Delete(index int) error {
  // lock searchers
  // lock inserters

  var prev *LLNode
  n := l.Root

  if index == 0 {
    if l.Root.Next != nil {
      l.Root = l.Root.Next
    } else {
      l.Root = nil
    }

    return nil
  }

  i := 0

  for ; n != nil && i != index; i++ {
    prev = n
    n = n.Next
  }

  if i != index {
    return errors.New("index does not exist")
    fmt.Println("error")
  }

  if n.Next != nil {
    prev.Next = n.Next
  } else {
    prev.Next = nil
  }

  return nil
}

func (l *LinkedList) Size() int {
  n := l.Root
  i := 0

  for ; n!=nil; i++ {
    n = n.Next
  }

  return i
}
