package mime

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		Types  Types
		header string
	}{
		{
			name:   "single",
			header: "text/html",
			Types: map[string]map[string]float64{
				"text": {
					"html": 1.,
				},
			},
		},
		{
			name:   "complicated",
			header: "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8",
			Types: map[string]map[string]float64{
				"text": {
					"html": 1.0,
				},
				"application": {
					"xhtml": 1,
					"xml":   0.9,
				},
				"*": {"*": 0.8},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewTypes(test.header)
			for group, subgroups := range g {
				for subgroup, weight := range subgroups {
					if test.Types[group][subgroup] != weight {
						t.Fail()
					}
				}
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		name  string
		Types Types
		a, b  string
	}{
		{
			name: "single",
			a:    "text/html",
			b:    "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8",
			Types: map[string]map[string]float64{
				"text": {
					"html": 1.,
				},
			},
		},
		{
			name: "double",
			a:    "text/html, application/xml",
			b:    "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8",
			Types: map[string]map[string]float64{
				"text": {
					"html": 1.,
				},
				"application": {
					"xml": .9,
				},
			},
		},
		{
			name:  "empty",
			a:     "text/html",
			b:     "image/gif",
			Types: map[string]map[string]float64{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := Intersect(NewTypes(test.a), NewTypes(test.b))
			for group, subgroups := range test.Types {
				for subgroup, weight := range subgroups {
					if g[group][subgroup] != weight {
						t.Fail()
					}
				}
			}
		})
	}
}

func TestParseList(t *testing.T) {
	tests := []struct {
		name  string
		Types Types
		list  []Type
	}{
		{
			name: "single",
			list: []Type{"text/html"},
			Types: map[string]map[string]float64{
				"text": {
					"html": 1.,
				},
			},
		},
		{
			name: "multiple",
			list: []Type{TextPlain, ApplicationJSON, ApplicationXML},
			Types: map[string]map[string]float64{
				"text": {
					"plain": 1.,
				},
				"application": {
					"json": 1.,
					"xml":  1.,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := Aggregate(test.list)
			for group, subgroups := range g {
				for subgroup, weight := range subgroups {
					if test.Types[group][subgroup] != weight {
						t.Fail()
					}
				}
			}
		})
	}
}
