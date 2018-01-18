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
			name:   "single-os",
			header: "application/octet-stream",
			Types: map[string]map[string]float64{
				"application": {
					"octet-stream": 1.,
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
			name:  "single-single",
			a:     "application/octet-strean",
			b:     "application/octet-strean",
			Types: map[string]map[string]float64{"application": {"octet-stream": 0}},
		},
		{
			name: "single-multi",
			a:    "text/html",
			b:    "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8",
			Types: map[string]map[string]float64{
				"text": {
					"html": 1.,
				},
			},
		},
		{
			name: "multi-multi",
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
						t.Logf("was expecting %v got %v instead", weight, g[group][subgroup])
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
