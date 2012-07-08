// Copyright 2009 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package diff implements a difference algorithm.
// The algorithm is described in "An O(ND) Difference Algorithm and its Variations", Eugene Myers, Algorithmica Vol. 1 No. 2, 1986, pp. 251-266.
package diff

// Diff returns the differences between two int slices.
func Diff(a, b []int) []Change {
	n := len(a)
	m := len(b)
	c := &context{a: a, b: b}
	if n > m {
		c.flags = make([]byte, n)
	} else {
		c.flags = make([]byte, m)
	}
	c.max = n + m + 1
	c.compare(0, 0, n, m)
	return c.result()
}

// A Change contains one or more deletions or inserts
// at one position in two sequences.
type Change struct {
	A, B int // position in input a and b
	Del  int // delete Del elements from input a
	Ins  int // insert Ins elements from input b
}

type context struct {
	a, b  []int  // inputs
	flags []byte // element bits 1 delete, 2 insert
	max   int
	// forward and reverse d-path endpoint x components
	forward, reverse []int
}

func (c *context) compare(aoffset, boffset, alimit, blimit int) {
	// eat common prefix
	for aoffset < alimit && boffset < blimit && c.a[aoffset] == c.b[boffset] {
		aoffset++
		boffset++
	}
	// eat common suffix
	for alimit > aoffset && blimit > boffset && c.a[alimit-1] == c.b[blimit-1] {
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
			for x < alimit && y < blimit && c.a[x] == c.b[y] {
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
			for x > aoffset && y > boffset && c.a[x-1] == c.b[y-1] {
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

func (c *context) result() (res []Change) {
	var x, y int
	n := len(c.a)
	m := len(c.b)
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
