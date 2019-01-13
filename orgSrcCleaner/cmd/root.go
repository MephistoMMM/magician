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
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/MephistoMMM/magician/lib"
	"github.com/MephistoMMM/magician/orgSrcCleaner/fileIterator"
	"github.com/MephistoMMM/magician/orgSrcCleaner/parser"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var log = lib.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "orgSrcCleaner <path>",
	Short: "Clean up linked static sources in org file.",
	Long: `orgSrcCleaner clean up static sources, like pictures and pdfs,
linked in your org files. It could restore your static sources to one directory,
to directories named by the filename of org file, or to directories named by the
headline.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		// the argument should be the path of actual file or directory
		src := args[0]
		if lib.IsNotExist(src) {
			return fmt.Errorf("src is not exist: %s", src)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		directory, err := filepath.Abs(args[0])
		if err != nil {
			log.Fatalln(err)
		}
		log.Debug(directory)
		iterator, err := fileIterator.NewFileIterator(directory)
		if err != nil {
			log.Fatalln(err)
		}

		var wg, wg2 sync.WaitGroup
		wg.Add(1)
		orgLinks := make(chan *parser.OrgLink, 10)
		go func(links <-chan *parser.OrgLink) {
			for link := range links {
				log.Debug(link)
				log.Infoln(link)
			}

			wg.Done()
		}(orgLinks)

		for iterator.HasNext() {
			file, err := iterator.Next()
			if err != nil {
				log.Fatalln(err)
			}
			log.Debug(file)

			wg2.Add(1)
			go func(file string) {
				links, err := lib.ScanLines(parser.NewOrgLinkParser(file))
				if err != nil {
					log.Errorln(err)
					wg2.Done()
					return
				}

				log.Debug(links)
				for _, link := range links {
					orgLinks <- link.(*parser.OrgLink)
				}
				wg2.Done()
			}(file)
		}

		wg2.Wait()
		close(orgLinks)

		wg.Wait()
	},
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.orgSrcCleaner.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".orgSrcCleaner" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".orgSrcCleaner")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
