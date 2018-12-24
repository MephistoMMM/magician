// Copyright © 2018 Mephis Pheies <mephistommm@gmail.com>
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

import (
	"strings"
	"testing"
)

type defaultFileLineParser struct{}

// FilePath ...
func (dflp *defaultFileLineParser) FilePath() string {
	return "./fs.go"
}

// Parse ...
func (dflp *defaultFileLineParser) Parse(line string) (interface{}, error) {
	if strings.HasPrefix(line, "func") {
		return line, nil
	}

	return nil, nil
}

func TestScanLines(t *testing.T) {
	lines, _ := ScanLines(&defaultFileLineParser{})

	for _, line := range lines {
		t.Log(line)
		if l := line.(string); !strings.HasPrefix(l, "func") {
			t.Errorf("Result %s has no prefix '%s'.", l, "func")
		}
	}

}
