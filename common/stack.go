package common

import (
	spn "github.com/RenatoGeh/gospn/spn"
)

// StackSPN is a stack of SPNs (nodes).
type StackSPN struct {
	data []spn.SPN
}

// Push puts element e on top of the pointer stack s.
func (s *StackSPN) Push(e spn.SPN) {
	s.data = append(s.data, e)
}

// Pop removes and returns the last element of pointer stack s.
func (s *StackSPN) Pop() spn.SPN {
	n := len(s.data) - 1
	e := s.data[n]
	s.data[n] = nil
	s.data = s.data[:n]
	return e
}

// Peek returns the top of the stack.
func (s *StackSPN) Peek() spn.SPN {
	return s.data[len(s.data)-1]
}

// Get returns the i-th element of the stack. Strongly discouraged, since this is a stack.
func (s *StackSPN) Get(i int) spn.SPN {
	return s.data[i]
}

// Size returns the size of pointer stack s.
func (s *StackSPN) Size() int { return len(s.data) }

// Empty returns whether pointer stack s is empty or not.
func (s *StackSPN) Empty() bool { return len(s.data) == 0 }
