// Copyright © 2020 Mephis Pheies <mephistommm@gmail.com>
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
package main

import (
	"github.com/MephistoMMM/magician/lib"
	"flag"
)

var flagArgsFile string


// TaskArgs arguments for algorithm
type TaskArgs struct {
	// Welfare the max sum of invoices
	Welfare float64      `yaml:"welfare"`
	// Invoices invoices to reimburse
	Invoices []float64   `yaml:"invoices"`
}

type InvoiceSet struct {
	Sum float64
	Elems []int
}

func init() {
	flag.StringVar(&flagArgsFile, "args-file", "", "File of arguments for algorithm")
}

func yieldMaxSet(limit float64,
	elems []float64, count int,
	amount float64,  arr []float64) (float64, []float64) {

	if len(arr) == 0 {
		res := make([]float64, count)
		copy(res, elems[0:count])
		return amount, res
	}

	max := amount
	maxSet := elems[0:count]
	for i, v := range arr {
		if amount + v > limit {
			continue
		}

		elems[count] = v
		sum, res := yieldMaxSet(limit, elems, count+1, amount+v, arr[i+1:])
		if sum > max {
			max = sum
			maxSet = res
		}
	}

	if max == amount {
		maxSet = make([]float64, count)
		copy(maxSet, elems[0:count])
	}

	return max, maxSet
}

func YieldMaxSet(limit float64, arr []float64) (float64, []float64) {
	elems := make([]float64, len(arr))
	return yieldMaxSet(limit, elems, 0, 0, arr)
}


func main() {
	flag.Parse()
	if flagArgsFile == "" {
		flag.PrintDefaults()
		return
	}


	args := &TaskArgs{}
	if err := lib.YamlLoadFromFile(flagArgsFile, args); err != nil {
		lib.Fatalf("Failed to load yaml[%s] !", flagArgsFile)
	}

	lib.Debugf("Loaded yaml: %v", *args)
    max, elems := YieldMaxSet(args.Welfare, args.Invoices)
	lib.Infof("Max value: %v", max)
	lib.Infof("List: %v", elems)

}
