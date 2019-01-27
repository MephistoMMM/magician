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
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/MephistoMMM/magician/lib"
	"github.com/spf13/cobra"
)

var installCmd string
var log = lib.Logger

const funcTemplate = `
function install_{{.FuncName}}() {
    echo "================================> install {{.Name}} ..."
    {{range .Dependences}}install_{{.FuncName}} $1
    {{end}}{{installCmd}} $1/{{.Pkg}}
    echo "================================> done"
}
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "depInstallScriptGenerator YAML_SRC",
	Short: "generate a shell script to install dependences.",
	Long:  `depInstallScriptGenerator generate a shell script include many function, which are used to install dependences defined in YAML_SRC.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		// the argument should be the path of actual yaml file
		src := args[0]
		if !lib.IsFile(src) {
			return fmt.Errorf(
				"src is not exist, is not a regular file: %s",
				src)
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

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&installCmd, "cmd", "yum install -y", "dependence install command")

}

type RenderFunc func(wr io.Writer, data Dependence) error

type Dependence struct {
	Name        string        `yaml:"name"`
	Pkg         string        `yaml:"pkg"`
	Dependences []*Dependence `yaml:"dependences"`
}

// FuncName return name used as function name
func (dp Dependence) FuncName() string {
	return strings.Replace(dp.Name, "-", "_", -1)
}

// GenScript ...
func (dp *Dependence) GenScript(render RenderFunc, buf io.Writer) error {
	for _, v := range dp.Dependences {
		if err := v.GenScript(render, buf); err != nil {
			return err
		}
	}

	if err := render(buf, *dp); err != nil {
		return err
	}
	return nil
}

type Requirements struct {
	Pkgs []*Dependence `yaml:"pkgs"`
}

// GenScript ...
func (rs Requirements) GenScript(installCmd string) (*bytes.Buffer, error) {
	renderedMap := make(map[string]bool, len(rs.Pkgs))
	funcMap := template.FuncMap{
		"installCmd": func() string { return installCmd },
	}

	tmpl, err := template.New(
		"DependenceTemp").Funcs(funcMap).Parse(funcTemplate)
	if err != nil {
		return nil, err
	}
	render := func(wr io.Writer, data Dependence) error {
		if v := renderedMap[data.Name]; v {
			return nil
		}

		renderedMap[data.Name] = true
		return tmpl.Execute(wr, data)
	}

	var buf bytes.Buffer
	// write shell header
	buf.WriteString(`#!/bin/bash`)
	for _, v := range rs.Pkgs {
		if err = v.GenScript(render, &buf); err != nil {
			return nil, err
		}
	}
	return &buf, nil
}

func run(cmd *cobra.Command, args []string) {
	requirements := &Requirements{}
	if err := lib.YamlLoadFromFile(args[0], requirements); err != nil {
		log.Fatalf("Parse yaml file error: %s", err)
	}

	buf, err := requirements.GenScript(installCmd)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if _, err = buf.WriteTo(os.Stdout); err != nil {
		log.Fatalf("error: %v", err)
	}
}
