// Copyright © 2019 Mephis Pheies <mephistommm@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package lib

type AlignCharIndex struct {
	Char    byte
	Indices []int
}

func newAlignCharIndex(c byte, num int) *AlignCharIndex {
	return &AlignCharIndex{
		Char:    c,
		Indices: make([]int, num),
	}
}

func FindAlignChar(strs []string) []*AlignCharIndex {
	summaries := make([]map[byte][]int, len(strs))

	// get summary of each str
	for i := range strs {
		summaries[i] = make(map[byte][]int)
		summary := summaries[i]
		for j := range strs[i] {
			indices, ok := summary[strs[i][j]]
			if !ok {
				indices = []int{}
			}
			indices = append(indices, j)
			summary[strs[i][j]] = indices
		}
	}

	result := []*AlignCharIndex{}
	// calculate result according to summaries
	for j := range strs[0] {
		c := strs[0][j]
		isAlian := true
		for i := 1; i < len(strs); i++ {
			if _, ok := summaries[i][c]; !ok {
				isAlian = false
				break
			}
		}

		// c is not a alian character
		if !isAlian {
			for _, summary := range summaries {
				delete(summary, c)
			}
			continue
		}

		// c is a alian character
		aci := newAlignCharIndex(c, len(strs))
		for i, summary := range summaries {
			aci.Indices[i] = summary[c][0]
			if len(summary[c]) == 1 {
				// the index is the last element
				delete(summary, c)
			} else {
				summary[c] = summary[c][1:]
			}
		}
		result = append(result, aci)
	}

	return result
}
