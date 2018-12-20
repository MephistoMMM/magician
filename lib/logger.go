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
	"os"

	log "github.com/sirupsen/logrus"
)

var Logger = log.New()

func init() {
	val, ok := os.LookupEnv("GRAFTER_ENV")
	if !ok {
		val = "pro"
	}

	if val == "pro" {
		Logger.SetLevel(log.InfoLevel)
	} else {
		Logger.SetLevel(log.DebugLevel)
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
}

var (
	Debug   = Logger.Debug
	Debugf  = Logger.Debugf
	Debugln = Logger.Debugln

	Info   = Logger.Info
	Infof  = Logger.Infof
	Infoln = Logger.Infoln

	Warn   = Logger.Warn
	Warnf  = Logger.Warnf
	Warnln = Logger.Warnln

	Error   = Logger.Error
	Errorf  = Logger.Errorf
	Errorln = Logger.Errorln

	Fatal   = Logger.Fatal
	Fatalf  = Logger.Fatalf
	Fatalln = Logger.Fatalln

	Panic   = Logger.Panic
	Panicf  = Logger.Panicf
	Panicln = Logger.Panicln
)
