package swagger

import (
	"fmt"
	"net/http"
	"path"

	"github.com/jennyservices/jenny/generator/internal/ir"
	"github.com/jennyservices/jenny/generator/util"
	"github.com/go-openapi/spec"
	"sevki.org/x/errors"
)

func deDupe(a []string) (r []string) {
	seen := make(map[string]bool)
	for _, s := range a {
		if _, ok := seen[s]; !ok {
			r = append(r, s)
			seen[s] = true
		}
	}
	return
}

var ErrNoPaths = errors.New("no paths")

func (s *swaggerDecoder) decodeMethods() *errors.Group {
	g := errors.NewGroup("methods")
	if s.spec.Paths == nil {
		return g.Add(ErrNoPaths)
	}

	for pathName, p := range s.spec.Paths.Paths {
		methods := map[string]*spec.Operation{
			http.MethodGet: p.Get,
			"POST":         p.Post,
			"DELETE":       p.Delete,
			"PUT":          p.Put,
			"PATCH":        p.Patch,
			"HEAD":         p.Head,
		}
		for methodName, method := range methods {
			if method == nil {
				continue
			}
			x := ir.Method{
				ID:             path.Join(s.svc.Name, util.NormalizeName(method.ID)),
				Name:           util.NormalizeName(method.ID),
				Path:           path.Join(s.svc.BasePath, pathName),
				HTTPMethod:     methodName,
				Description:    method.Summary,
				Consumes:       deDupe(append(method.Consumes, s.spec.Consumes...)),
				Produces:       deDupe(append(method.Produces, s.spec.Produces...)),
				ParameterOrder: make([]string, len(method.Parameters)),
				Parameters:     make(map[string]ir.Schema),
			}
			for _, sec := range method.Security {
				for secName, scopes := range sec {
					if secScheme, ok := s.spec.SecurityDefinitions[secName]; ok && secScheme.Type == "oauth2" {
						x.Scopes = deDupe(append(x.Scopes, scopes...))
					}
				}
			}

			for i, param := range method.Parameters {
				name := util.NormalizeName(param.Name)
				x.ParameterOrder[i] = name
				paramSchema := ir.Schema{
					Name:           name,
					CannonicalName: param.Name,
				}
				if len(param.Description) > 0 {
					paramSchema.Description = String(param.Description)
				}
				if len(param.Type) < 1 {
					paramSchema.Type = name
				} else {
					paramSchema.Type = param.Type
				}
				if param.Schema != nil {
					schme, eg := s.swaggerToIR(name, methodName, param.Schema)
					if eg.Errored() {
						g.Add(eg)
						continue
					}
					paramSchema = *schme
				}

				paramSchema.Required = param.Required
				if param.Default != nil {
					paramSchema.Default = String(fmt.Sprint(param.Default))
				}
				paramSchema.Order = i
				paramSchema.ID = path.Join(x.ID, name)
				paramSchema.Location = String(param.In)
				paramSchema.CannonicalName = param.Name
				x.Parameters[name] = paramSchema
			}
			if method.Responses != nil {
				x.Responses = make(map[string]ir.Response)
				hasDefault := false

				if method.Responses.Default != nil {
					x.Responses["Default"] = ir.Response{
						Default:    true,
						Error:      true,
						HTTPStatus: http.StatusInternalServerError,
					}

					s.decodeResponse(*method.Responses.Default, "Default")
					hasDefault = true
				}

				for respCode, response := range method.Responses.StatusCodeResponses {
					r := ir.Response{
						Default:    false,
						Error:      false,
						HTTPStatus: respCode,
					}
					if respCode > 400 {
						r.Error = true
					}
					if !hasDefault && (respCode == 200 || respCode == 201) {
						hasDefault = true
						r.Default = true
					}
					if len(response.Headers) > 0 || response.Schema != nil {
						r.Returns = make(map[string]ir.Schema)
					}
					if response.Schema != nil {
						schema, _ := s.swaggerToIR("body", x.ID, response.Schema)
						schema.Required = false
						schema.Location = String("body")
						r.Returns["body"] = *schema
					}

					for headerName, s := range response.Headers {
						r.Returns[headerName] = ir.Schema{
							ID:       path.Join(x.ID, "Headers", headerName),
							Location: String("header"),
							Name:     headerName,
							Type:     s.Type,
							Required: true,
						}
					}
					x.Responses[http.StatusText(respCode)] = r
					_, _ = respCode, response
				}

			}

			s.svc.Resources[x.ID] = x
		}
	}
	return g
}

func (s *swaggerDecoder) decodeResponse(response spec.Response, respName string) (ir.Response, *errors.Group) {
	g := errors.NewGroup("response")
	name := util.NormalizeName(respName)
	r := ir.Response{
		Default: false,
	}

	x, err := s.swaggerToIR("", name, response.Schema)
	g = g.Add(err)
	if !err.Errored() {
		r.Returns = map[string]ir.Schema{
			respName: *x,
		}
	}
	return r, err
}
