package templates

var Operations = `
{{range .}}
	{{$privateType := .Name | makePrivate}}
	{{$exportType := .Name | makePublic}}

	type {{$exportType}} interface {
		{{range .Operations}}
			{{$faults := len .Faults}}
			{{$soapAction := findSOAPAction .Name $privateType}}
			{{$requestType := findMessageType .Input.Message | replaceReservedWords | makePublic}}
			{{$responseType := findMessageType .Output.Message | replaceReservedWords | makePublic}}

			{{/*if ne $soapAction ""*/}}
			{{if gt $faults 0}}
			// Error can be either of the following types:
			// {{range .Faults}}
			//   - {{.Name}} {{.Doc}}{{end}}{{end}}
			{{if ne .Doc ""}}/* {{.Doc}} */{{end}}
			{{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error)
			{{/*end*/}}
			{{makePublic .Name | replaceReservedWords}}Context (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error)
			{{/*end*/}}
		{{end}}
	}

	type {{$privateType}} struct {
		client *proxy.Client
	}

	func New{{$exportType}}(client *proxy.Client) {{$exportType}} {
		return &{{$privateType}}{client: client}
	}

	{{range .Operations}}
		{{$requestType := findMessageType .Input.Message | replaceReservedWords | makePublic}}
		{{$soapAction := findSOAPAction .Name $privateType}}
		{{$responseType := findMessageType .Output.Message | replaceReservedWords | makePublic}}
		func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}}Context (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error) {
			{{if ne $responseType ""}}response := new({{$responseType}}){{end}}
			err := service.client.CallContext(ctx, "{{if ne $soapAction ""}}{{$soapAction}}{{else}}''{{end}}", {{if ne $requestType ""}}request{{else}}nil{{end}}, {{if ne $responseType ""}}response{{else}}struct{}{}{{end}})
			if err != nil {
				return {{if ne $responseType ""}}nil, {{end}}err
			}

			return {{if ne $responseType ""}}response, {{end}}nil
		}

		func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error) {
			return service.{{makePublic .Name | replaceReservedWords}}Context(context.Background(), {{if ne $requestType ""}}request{{end}})
		}
	{{end}}
{{end}}
`