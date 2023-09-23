package templates

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

// TxtCustomFuncMap returns a FuncMap with the custom functions used in the templates.
// You can implement your own custom functions inside of this func.
func CustomTxtFuncMap(tmpl *template.Template) template.FuncMap {
	var funcMap template.FuncMap = map[string]interface{}{}

	// "include" function from Helm. Copied from:
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

	// toFile saves the result of executing a named template with data to a file at the specified path.
	// It accepts a 'path' for the file to be saved, a 'name' of the template to execute, and 'data' to
	// be provided as input to the template execution.
	//
	// The 'path' can be either a relative or an absolute file path. If the directory for the file does
	// not exist, it will be created automatically.
	funcMap["toFile"] = func(path, name string, data interface{}) (string, error) {
		// Execute the template.
		buf := bytes.NewBuffer(nil)
		if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}

		// If needed, create the directory where the file will be saved.
		directory := filepath.Dir(path)
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return "", err
		}

		// Create the file.
		err := os.WriteFile(path, buf.Bytes(), os.ModePerm)
		if err != nil {
			return "", err
		}

		return "", nil
	}

	// Implement your own custom functions here...

	return funcMap
}
