package templates

import (
	"bytes"
	"text/template"
)

// TxtCustomFuncMap returns a FuncMap with the custom functions used in the templates.
// You can implement your own custom functions inside of this func.
func CustomTxtFuncMap(tmpl *template.Template) template.FuncMap {
	var funcMap template.FuncMap = map[string]interface{}{}

	// 'include' function from Helm. Copied from:
	// https://github.com/helm/helm/blob/8648ccf5d35d682dcd5f7a9c2082f0aaf071e817/pkg/engine/engine.go#L148
	//
	// The include function allows you to bring in another template,
	// and then pass the results to other template functions.
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	// Implement your own custom functions here...

	return funcMap
}
