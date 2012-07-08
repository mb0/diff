// Copyright 2009 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"testing"
)

type testcase struct {
	name string
	data *Ints
	res  []Change
}

var tests = []testcase{
	{"shift",
		&Ints{
			{1, 2, 3},
			{0, 1, 2, 3},
		},
		[]Change{{0, 0, 0, 1}},
	},
	{"push",
		&Ints{
			{1, 2, 3},
			{1, 2, 3, 4},
		},
		[]Change{{3, 3, 0, 1}},
	},
	{"unshift",
		&Ints{
			{0, 1, 2, 3},
			{1, 2, 3},
		},
		[]Change{{0, 0, 1, 0}},
	},
	{"pop",
		&Ints{
			{1, 2, 3, 4},
			{1, 2, 3},
		},
		[]Change{{3, 3, 1, 0}},
	},
	{"all changed",
		&Ints{
			{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			{10, 11, 12, 13, 14},
		},
		[]Change{
			{0, 0, 10, 5},
		},
	},
	{"all same",
		&Ints{
			{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		[]Change{},
	},
	{"wrap",
		&Ints{
			{1},
			{0, 1, 2, 3},
		},
		[]Change{
			{0, 0, 0, 1},
			{1, 2, 0, 2},
		},
	},
	{"snake",
		&Ints{
			{0, 1, 2, 3, 4, 5},
			{1, 2, 3, 4, 5, 6},
		},
		[]Change{
			{0, 0, 1, 0},
			{6, 5, 0, 1},
		},
	},
	// note: input is ambiguous
	// first two traces differ from fig.1
	// it still is a lcs and ses path
	{"paper fig. 1",
		&Ints{
			{1, 2, 3, 1, 2, 2, 1},
			{3, 2, 1, 2, 1, 3},
		},
		[]Change{
			{0, 0, 1, 1},
			{2, 2, 1, 0},
			{5, 4, 1, 0},
			{7, 5, 0, 1},
		},
	},
}

func TestDiffAB(t *testing.T) {
	for _, test := range tests {
		res := test.data.Diff()
		if len(res) != len(test.res) {
			t.Error(test.name, "expected length", len(test.res), "for", res)
			continue
		}
		for i, c := range test.res {
			if c != res[i] {
				t.Error(test.name, "expected ", c, "got", res[i])
			}
		}
	}
}

func TestDiffBA(t *testing.T) {
	// interesting: fig.1 Diff(b, a) results in the same path as `diff -d a b`
	tests[len(tests)-1].res = []Change{
		{0, 0, 2, 0},
		{3, 1, 1, 0},
		{5, 2, 0, 1},
		{7, 5, 0, 1},
	}
	for _, test := range tests {
		res := Diff(&Ints{test.data[1], test.data[0]})
		if len(res) != len(test.res) {
			t.Error(test.name, "expected length", len(test.res), "for", res)
			continue
		}
		for i, c := range test.res {
			// flip change data also
			rc := Change{c.B, c.A, c.Ins, c.Del}
			if rc != res[i] {
				t.Error(test.name, "expected ", rc, "got", res[i])
			}
		}
	}
}

func BenchmarkDiff(b *testing.B) {
	t := tests[len(tests)-1]
	for i := 0; i < b.N; i++ {
		t.data.Diff()
	}
}

func BenchmarkDiffRunes(b *testing.B) {
	data := &Runes{[]rune("1231221"), []rune("321213")}
	for i := 0; i < b.N; i++ {
		Diff(data)
	}
}
