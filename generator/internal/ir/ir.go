// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ir

// Type is the
type Type struct {
	Type   string  `json:"type"`
	Format *string `json:"format,omitempty"`
}

// Schema is the Intermediary representation of a Schema
type Schema struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	CannonicalName   string            `json:"cannonical_name"`
	Type             string            `json:"type"`
	Required         bool              `json:"required"`
	Default          *string           `json:"default"`
	Order            int               `json:"order"`
	Format           *string           `json:"format,omitempty"`
	Description      *string           `json:"description,omitempty"`
	Location         *string           `json:"location,omitempty"`
	Enum             []string          `json:"enum,omitempty"`
	EnumDescriptions []string          `json:"enumDescriptions,omitempty"`
	Properties       map[string]Schema `json:"properties,omitempty"`
	Item             *Schema           `json:"items,omitempty"`
}

// Method is the Intermediary representation of a Method
type Method struct {
	ID             string              `json:"id"`
	Name           string              `json:"Name"`
	Path           string              `json:"path"`
	HTTPMethod     string              `json:"httpMethod"`
	Description    string              `json:"description"`
	Produces       []string            `json:"produces,omitempty"`
	Consumes       []string            `json:"consumes,omitempty"`
	Parameters     map[string]Schema   `json:"parameters"`
	ParameterOrder []string            `json:"parameterOrder"`
	Responses      map[string]Response `json:"responses,omitempty"`
	Scopes         []string            `json:"scopes,omitempty"`
}

// Response is the Intermediary representation of a Service
type Response struct {
	HTTPStatus int               `json:"http_status"`
	Returns    map[string]Schema `json:"returns,omitempty"`
	Default    bool              `json:"default"`
	Error      bool              `json:"error"`
}

// Service is the Intermediary representation of a Service
type Service struct {
	Kind              string            `json:"kind"`
	Etag              string            `json:"etag"`
	DiscoveryVersion  string            `json:"discoveryVersion"`
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	CanonicalName     string            `json:"canonicalName"`
	Version           string            `json:"version"`
	Revision          string            `json:"revision"`
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	OwnerDomain       string            `json:"ownerDomain"`
	OwnerName         string            `json:"ownerName"`
	DocumentationLink string            `json:"documentationLink"`
	Protocol          string            `json:"protocol"`
	BaseURL           string            `json:"baseUrl"`
	BasePath          string            `json:"basePath"`
	RootURL           string            `json:"rootUrl"`
	ServicePath       string            `json:"servicePath"`
	BatchPath         string            `json:"batchPath"`
	Parameters        map[string]Schema `json:"parameters"`
	Schemas           map[string]Schema `json:"schemas"`
	Resources         map[string]Method `json:"resources"`
}

// Decoder is an interface that takes raw data bytes and converts them to a
// service definition
type Decoder interface {
	Decode() (*Service, error)
}

// Encoder is an interface that takes a service definition and prints the code
// for a language.
type Encoder interface {
	Encode(*Service) error
}
