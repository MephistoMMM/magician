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

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/MephistoMMM/magician/lib"
	"github.com/spf13/cobra"
)

var log = lib.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "npmPackageParser DEPENDENCIES",
	Short: "npmPackageParser parser result of `npm list --long --parseable` to json represented dependencies.",
	Long: `npmPackageParser parser result of 'npm list --long --parseable' to json represented dependencies.

Before call this command, you should run 'npm list --long --parseable > dependencies' to generate a file stored raw dependencies data. Then you could use this command to parser the generated file to get its json format.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		// the argument should be the path of actual dependencies file
		dependencies := args[0]
		if !lib.IsFile(dependencies) {
			return fmt.Errorf(
				"DEPENDENCIES is not exist or not a regular file: %s",
				dependencies)
		}

		return nil
	},
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type packageParser struct {
	filePath string
	match    *regexp.Regexp
	template string
}

func NewPackageParser(filePath string, match *regexp.Regexp, template string) lib.FileLineParser {
	return &packageParser{
		filePath: filePath,
		match:    match,
		template: template,
	}
}

// FilePath ...
func (pp *packageParser) FilePath() string {
	return pp.filePath
}

// Parse ...
func (pp *packageParser) Parse(line string) (interface{}, error) {
	if !pp.match.MatchString(line) {
		log.Infof("Not match: %s\n", line)
		return nil, nil
	}

	dst := []byte{}

	dst = pp.match.ExpandString(
		dst,
		pp.template,
		line,
		pp.match.FindStringSubmatchIndex(line),
	)
	return dst, nil
}

func run(cmd *cobra.Command, args []string) {
	dependencies := args[0]

	pp := NewPackageParser(
		dependencies,
		regexp.MustCompile(
			"^/.+?:(?P<package>.+?)@(?P<version>.+?):.+$"),
		"\"$package\": \"$version\"")

	results, err := lib.ScanLines(pp)
	if err != nil {
		log.Fatalf("Scan lines error: %s\n", err)
	}

	slimContent := make(map[string]bool, len(results)*2)

	for _, result := range results {
		slimContent[string(result.([]byte))] = true
	}

	var buf bytes.Buffer
	buf.WriteString("{\n")
	for k, _ := range slimContent {
		if len(k) > 0 {
			buf.WriteString(fmt.Sprintf("\t%s,\n", k))
		}
	}
	buf.WriteString("}\n")

	fmt.Println(buf.String())
}
