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
package fileIterator

import (
	"github.com/MephistoMMM/magician/lib"
)

func orgFilterChain() lib.FilterSupport {
	var filterChain lib.FilterSupport
	tmp, _ := lib.NewFilterIgnoreDotSupport()
	filterChain = tmp
	tmp, _ = lib.NewFilterIgnoreUnregularSupport()
	filterChain.SetNext(tmp)
	tmp, err := lib.NewFilterRegexpMatchSupport(`\.org$`)
	if err != nil {
		panic(err)
	}
	filterChain.SetNext(tmp)
	return filterChain
}

// NewOrgFileIterator create a file iterator to return org files one by one
func NewOrgFileIterator(directory string) (lib.FileIterator, error) {
	iterator, _ := lib.NewFileIterator(directory, orgFilterChain())
	return iterator, nil
}
