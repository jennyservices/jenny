// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golang

import (
	"io"
	"io/ioutil"

	"golang.org/x/tools/imports"
)

func processFile(r io.Reader) ([]byte, error) {
	bytz, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	opts := imports.Options{
		Fragment:  false,
		AllErrors: false, // Report all errors (not just the first 10 on different lines)

		Comments:  true, // Print comments (true if nil *Options provided)
		TabIndent: true, // Use tabs for indent (true if nil *Options provided)
		TabWidth:  8,    // Tab width (8 if nil *Options provided)

		FormatOnly: false, // Disable the insertion and deletion of imports
	}
	_ = opts
	return imports.Process("FAKEFILE.go", bytz, &opts)
}
