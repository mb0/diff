// Copyright 2009 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff_test

import (
	"fmt"
	"github.com/mb0/diff"
)

var names = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
}

// Diff on inputs with different representations
type MixedInput struct {
	A []int
	B []string
}

func (m *MixedInput) Equal(a, b int) bool {
	return m.A[a] == names[m.B[b]]
}

func ExampleInterface() {
	m := &MixedInput{
		[]int{1, 2, 3, 1, 2, 2, 1},
		[]string{"three", "two", "one", "two", "one", "three"},
	}
	changes := diff.Diff(len(m.A), len(m.B), m)
	for _, c := range changes {
		fmt.Println("change at", c.A, c.B)
	}
}
