// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/jennyservices/jenny/generator"
	"github.com/spf13/cobra"
)

var (
	file, out, pkg string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the code from Swagger definition",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {

		gen := generator.New(file, out)

		err := gen.Generate()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&file, "file", "f", "./transport/v1/swagger.yaml", "name of the file to read")
	generateCmd.Flags().StringVarP(&out, "out", "o", "./transport/v1", "where to output the files ")
	generateCmd.Flags().StringVarP(&pkg, "pkg", "p", "v1", "denotes the name of the pakage")
}
