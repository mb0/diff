diff
====

A difference algorithm package for go.

The algorithm is described by Eugene Myers in
["An O(ND) Difference Algorithm and its Variations"](http://www.xmailserver.org/diff2.pdf).

Example
-------
You can use diff.Ints and diff.Runes

    d := &diff.Runes{[]rune("sögen"), []rune("mögen")}
    d.Diff() // returns []Changes{{0,0,1,1}}

or you can implement diff.Interface and use diff.Diff(data)

    type MixedInput struct {
    	A []int
    	B []string
    }
    func (m *MixedInput) N() int {
    	return len(m.A)
    }
    func (m *MixedInput) M() int {
    	return len(m.B)
    }
    func (m *MixedInput) Equal(a, b int) bool {
    	return m.A[a] == len(m.B[b])
    }

Documentation at http://go.pkgdoc.org/github.com/mb0/diff
