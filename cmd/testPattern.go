// Copyright (c) 2023-2024, Open Communications Security
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
//  3. Neither the name of the copyright holder nor the names of its
//     contributors may be used to endorse or promote products derived from
//     this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	testPattern string
)

// testPatternCmd represents the testPattern command
var testPatternCmd = &cobra.Command{
	Use:   "test_pattern --pattern <pattern> [<test path>]...",
	Short: "Tests a path pattern regular expression.",
	Long: `Tests a path pattern regular expression.

All non flag parammeters, if any will be tested agains the pattern,
	`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if testPattern == "" {
			return fmt.Errorf("pattern is required")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		// Compile the pattern
		p, err := regexp.Compile(testPattern)
		if err != nil {
			return err
		}
		fmt.Printf("Pattern: %s\n", p.String())
		for _, s := range args {
			fmt.Printf(" %s => %v\n", s, p.Match([]byte(s)))
		}
		fmt.Print(args)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testPatternCmd)

	testPatternCmd.Flags().StringVarP(&testPattern, "pattern", "p", "", "Path pattern regex to be tested.")
}
