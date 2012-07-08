// Copyright 2009 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package diff implements a difference algorithm.
// The algorithm is described in "An O(ND) Difference Algorithm and its Variations", Eugene Myers, Algorithmica Vol. 1 No. 2, 1986, pp. 251-266.
package diff

// A type that satisfies diff.Interface can be diffed by this package.
// It typically has two sequences A and B of comparable elements.
type Interface interface {
	// N is the number of elements in A, called once
	N() int
	// M is the number of elements in B, called once
	M() int
	// Equal returns whether the elements at a and b are considered equal.
	// Called repeatedly with 0<=a<N and 0<=b<M
	Equal(a, b int) bool
}

// Ints attaches diff.Interface methods to an array of two int slices
type Ints [2][]int

func (i *Ints) N() int {
	return len(i[0])
}
func (i *Ints) M() int {
	return len(i[1])
}
func (i *Ints) Equal(a, b int) bool {
	return i[0][a] == i[1][b]
}
func (i *Ints) Diff() []Change {
	return Diff(i)
}

// Runes attaches diff.Interface methods to an array of two rune slices
type Runes [2][]rune

func (r *Runes) N() int {
	return len(r[0])
}
func (r *Runes) M() int {
	return len(r[1])
}
func (r *Runes) Equal(a, b int) bool {
	return r[0][a] == r[1][b]
}
func (r *Runes) Diff() []Change {
	return Diff(r)
}

// Diff returns the differences of data.
func Diff(data Interface) []Change {
	n := data.N()
	m := data.M()
	c := &context{data: data}
	if n > m {
		c.flags = make([]byte, n)
	} else {
		c.flags = make([]byte, m)
	}
	c.max = n + m + 1
	c.compare(0, 0, n, m)
	return c.result(n, m)
}

// A Change contains one or more deletions or inserts
// at one position in two sequences.
type Change struct {
	A, B int // position in input a and b
	Del  int // delete Del elements from input a
	Ins  int // insert Ins elements from input b
}

type context struct {
	data  Interface
	flags []byte // element bits 1 delete, 2 insert
	max   int
	// forward and reverse d-path endpoint x components
	forward, reverse []int
}

func (c *context) compare(aoffset, boffset, alimit, blimit int) {
	// eat common prefix
	for aoffset < alimit && boffset < blimit && c.data.Equal(aoffset, boffset) {
		aoffset++
		boffset++
	}
	// eat common suffix
	for alimit > aoffset && blimit > boffset && c.data.Equal(alimit-1, blimit-1) {
		alimit--
		blimit--
	}
	// both equal or b inserts
	if aoffset == alimit {
		for boffset < blimit {
			c.flags[boffset] |= 2
			boffset++
		}
		return
	}
	// a deletes
	if boffset == blimit {
		for aoffset < alimit {
			c.flags[aoffset] |= 1
			aoffset++
		}
		return
	}
	x, y := c.findMiddleSnake(aoffset, boffset, alimit, blimit)
	c.compare(aoffset, boffset, x, y)
	c.compare(x, y, alimit, blimit)
}

func (c *context) findMiddleSnake(aoffset, boffset, alimit, blimit int) (int, int) {
	// midpoints
	fmid := aoffset - boffset
	rmid := alimit - blimit
	// correct offset in d-path slices
	foff := c.max - fmid
	roff := c.max - rmid
	isodd := (rmid-fmid)&1 != 0
	maxd := (alimit - aoffset + blimit - boffset + 2) / 2
	// allocate when first used
	if c.forward == nil {
		c.forward = make([]int, 2*c.max)
		c.reverse = make([]int, 2*c.max)
	}
	c.forward[c.max+1] = aoffset
	c.reverse[c.max-1] = alimit
	var x, y int
	for d := 0; d <= maxd; d++ {
		// forward search
		for k := fmid - d; k <= fmid+d; k += 2 {
			if k == fmid-d || k != fmid+d && c.forward[foff+k+1] < c.forward[foff+k-1] {
				x = c.forward[foff+k+1] // down
			} else {
				x = c.forward[foff+k-1] + 1 // right
			}
			y = x - k
			for x < alimit && y < blimit && c.data.Equal(x, y) {
				x++
				y++
			}
			c.forward[foff+k] = x
			if isodd && k > rmid-d && k < rmid+d {
				if c.reverse[roff+k] <= c.forward[foff+k] {
					return x, x - k
				}
			}
		}
		// reverse search x,y correspond to u,v
		for k := rmid - d; k <= rmid+d; k += 2 {
			if k == rmid+d || k != rmid-d && c.reverse[roff+k-1] < c.reverse[roff+k+1] {
				x = c.reverse[roff+k-1] // up
			} else {
				x = c.reverse[roff+k+1] - 1 // left
			}
			y = x - k
			for x > aoffset && y > boffset && c.data.Equal(x-1, y-1) {
				x--
				y--
			}
			c.reverse[roff+k] = x
			if !isodd && k >= fmid-d && k <= fmid+d {
				if c.reverse[roff+k] <= c.forward[foff+k] {
					// lookup opposite end
					x = c.forward[foff+k]
					return x, x - k
				}
			}
		}
	}
	panic("should never be reached")
}

func (c *context) result(n, m int) (res []Change) {
	var x, y int
	for x < n || y < m {
		if x < n && y < m && c.flags[x]&1 == 0 && c.flags[y]&2 == 0 {
			x++
			y++
		} else {
			a := x
			b := y
			for x < n && (y >= m || c.flags[x]&1 != 0) {
				x++
			}
			for y < m && (x >= n || c.flags[y]&2 != 0) {
				y++
			}
			if a < x || b < y {
				res = append(res, Change{a, b, x - a, y - b})
			}
		}
	}
	return
}
