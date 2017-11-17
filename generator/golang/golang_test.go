package golang

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Typeform/jenny/generator/internal/ir"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var (
	update = flag.Bool("u", false, "-u")
)

func TestSwagger(t *testing.T) {
	flag.Parse()
	tests := []struct {
		name     string
		service  ir.Service
		err      error
		template string
	}{
		{
			name:     "schema",
			err:      nil,
			template: "service",
			service: ir.Service{
				Name:             "TestService1",
				DiscoveryVersion: "v1",
				Schemas: map[string]ir.Schema{
					"Error": ir.Schema{
						Type: "object",
						Name: "Error",
						Properties: map[string]ir.Schema{
							// Number fields
							"numberField": ir.Schema{
								Name:        "numberField",
								Type:        "number",
								Description: String("This describes the first field"),
							},
							"numberFieldFloat": ir.Schema{
								Name:        "numberFieldFloat",
								Type:        "number",
								Format:      String("float"),
								Description: String("This describes the first field"),
							},
							"numberFieldDouble": ir.Schema{
								Name:        "numberFieldDouble",
								Type:        "number",
								Format:      String("double"),
								Description: String("This describes the first field"),
							},
							// Integer fields
							"integerField": ir.Schema{
								Name:        "integerField",
								Type:        "integer",
								Description: String("This describes the first field"),
							},
							"integerField32": ir.Schema{
								Name:        "integerField",
								Type:        "integer",
								Format:      String("int32"),
								Description: String("This describes the first field"),
							},
							"integerField64": ir.Schema{
								Name:        "integerField",
								Type:        "integer",
								Format:      String("int64"),
								Description: String("This describes the first field"),
							},
							// String Fields
							"stringField": ir.Schema{
								Name:        "stringField",
								Type:        "string",
								Description: String("This describes the first field"),
							},
							"stringFieldDate": ir.Schema{
								Name:        "stringFieldDate",
								Type:        "string",
								Format:      String("date"),
								Description: String("This describes the first field"),
							},
							"stringFieldByte": ir.Schema{
								Name:        "stringFieldByte",
								Type:        "string",
								Format:      String("byte"),
								Description: String("This describes the first field"),
							},
							"stringFieldPassword": ir.Schema{
								Name:        "stringFieldPassword",
								Type:        "string",
								Format:      String("password"),
								Description: String("This describes the first field"),
							},
							"stringFieldURL": ir.Schema{
								Name:        "stringFieldURL",
								Type:        "string",
								Format:      String("url"),
								Description: String("This describes the first field"),
							},
							// Bool fields
							"boolField": ir.Schema{
								Name:        "stringFieldURL",
								Type:        "boolean",
								Description: String("This describes the first field"),
							},
							// FuckedUP field
							"goField": ir.Schema{
								Name:        "stringFieldURL",
								Type:        "net.IP",
								Required:    false,
								Description: String("This describes the first field"),
							},
						},
					},
				},
			},
		},
		{
			name:     "resources",
			err:      nil,
			template: "service",
			service: ir.Service{
				Name:             "TestService2",
				DiscoveryVersion: "v1",
				Resources: map[string]ir.Method{
					"Login": ir.Method{
						Name:       "Login",
						ID:         "login",
						Path:       "/login",
						HTTPMethod: http.MethodGet,
						Scopes:     []string{"pets"},
						Responses: map[string]ir.Response{
							"Found": ir.Response{
								HTTPStatus: http.StatusFound,
								Returns: map[string]ir.Schema{
									"Location": ir.Schema{
										Name:     "Location",
										Location: String("header"),
										Type:     "string",
									},
								},
							},
						},
					},
					"GetUser": ir.Method{
						Name:       "GetUser",
						ID:         "block",
						Path:       "/users/{id}",
						HTTPMethod: http.MethodPost,
						Parameters: map[string]ir.Schema{
							"id": ir.Schema{
								Name:     "id",
								Type:     "number",
								Location: String("path"),
							},
							"First": ir.Schema{
								Name:     "first",
								Type:     "number",
								Location: String("query"),
								Required: true,
								Default:  String("10"),
							},
							"User": ir.Schema{
								Name:     "User",
								Type:     "user.User",
								Location: String("body"),
								Required: true,
								Default:  String("{ bla}"),
							},
						},
						Responses: map[string]ir.Response{
							"OK": ir.Response{
								HTTPStatus: 200,
								Default:    true,
							},
						},
					},
				},
			},
		},
		{
			name:     "defaultresponse",
			err:      nil,
			template: "service",
			service: ir.Service{
				Name:             "TestService3",
				DiscoveryVersion: "v3",
				Resources: map[string]ir.Method{
					"Login": ir.Method{
						Name:       "Login",
						ID:         "login",
						Path:       "/login",
						HTTPMethod: http.MethodGet,
						Responses: map[string]ir.Response{
							"OK": ir.Response{
								HTTPStatus: http.StatusOK,
								Default:    true,
								Returns: map[string]ir.Schema{
									"body": ir.Schema{
										Name:     "body",
										Location: String("body"),
										Type:     "user.User",
									},
								},
							},
							"Found": ir.Response{
								HTTPStatus: http.StatusFound,
								Returns: map[string]ir.Schema{
									"Location": ir.Schema{
										Name:     "Location",
										Location: String("header"),
										Type:     "string",
									},
								},
							},
						},
					},
					"GetUser": ir.Method{
						Name:       "GetUser",
						ID:         "block",
						Path:       "/users/{id}",
						HTTPMethod: http.MethodPost,
						Parameters: map[string]ir.Schema{
							"id": ir.Schema{
								Name:     "id",
								Type:     "number",
								Location: String("path"),
							},
							"First": ir.Schema{
								Name:     "first",
								Type:     "number",
								Location: String("query"),
								Required: true,
								Default:  String("10"),
							},
							"User": ir.Schema{
								Name:     "User",
								Type:     "user.User",
								Location: String("body"),
								Required: true,
								Default:  String("{ bla}"),
							},
						},
						Responses: map[string]ir.Response{
							"OK": ir.Response{
								HTTPStatus: 200,
								Default:    true,
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			enc := New(buf, test.template)
			err := enc.Encode(&test.service)
			if err != test.err {
				t.Log(err)
				t.Fail()
			}

			gold := fmt.Sprintf("testdata/%s.go", test.name)
			var goldFile *os.File

			if *update {
				goldFile, err = os.Create(gold)
				if err != nil {
					t.Fatal(err)
				}
				goldFile.Write(buf.Bytes())
				goldFile.Close()
				return
			}

			bytz, err := ioutil.ReadFile(gold)
			if err != nil {
				log.Fatal(err)
			}
			dmp := diffmatchpatch.New()

			diffs := dmp.DiffMain(string(bytz), buf.String(), false)

			if bytes.Compare(bytz, buf.Bytes()) != 0 {
				t.Log(dmp.DiffPrettyText(diffs))
				t.Fail()
			}
		})
	}
}
