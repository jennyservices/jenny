// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golang

import (
	"testing"
)

// Bug report: https://github.com/jennyservices/jenny/issues/20
// Reported by: @marc-gr
// Issue: Now, if a field type is number, but its format is float, jenny ignores
// the format and set it as int. We should take into account the format when
// generating the code.
func TestIssue20(t *testing.T) {
	tests := []struct {
		name     string
		typ      string
		format   *string
		expected string
	}{
		{
			name:     "number",
			typ:      "number",
			format:   nil,
			expected: "int",
		},
		{
			name:     "float",
			typ:      "number",
			format:   String("float"),
			expected: "float64",
		},
		{
			name:     "string",
			typ:      "string",
			format:   nil,
			expected: "string",
		},
		{
			name:     "bytes",
			typ:      "string",
			format:   String("binary"),
			expected: "[]byte",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := getSimpleType(test.typ, test.format)
			if actual != test.expected {
				t.Logf("was expecting %q got %q instead", test.expected, actual)
				t.Fail()
			}
		})
	}
}
