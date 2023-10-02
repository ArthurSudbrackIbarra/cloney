package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

	// "toFile" function is a custom function provided by Cloney, which allows users to dynamically
	// create files from a template. This function has 2 hidden parameters 'templateDir' and 'fileDir',
	// representing the directory of the template being processed and the directory of the file currently
	// being processed, respectively. This parameter is utilized to determine the absolute path where the
	// generated file will be saved, relative to the file being processed.
	//
	// During template execution, 'templateDir' and 'fileDir' are automatically injected into the function.
	funcMap["toFile"] = func(templateDir, fileDir, relativePath, name string, data interface{}) (string, error) {
		// Execute the template.
		buf := bytes.NewBuffer(nil)
		if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}

		// Calculate the absolute path of the file.
		absPath := filepath.Join(fileDir, relativePath)

		// If on Windows, replace forward slashes with backslashes.
		if os.PathSeparator == '\\' {
			absPath = strings.ReplaceAll(absPath, "/", "\\")
		}

		// If the path is out of the scope of the template directory, return an error.
		if !strings.HasPrefix(absPath, templateDir) {
			return "", fmt.Errorf("Cannot create file outside the scope of the template directory: %s", relativePath)
		}

		// If needed, create the directory where the file will be saved.
		directory := filepath.Dir(absPath)
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return "", err
		}

		// Create the file.
		err := os.WriteFile(absPath, buf.Bytes(), os.ModePerm)
		if err != nil {
			return "", err
		}

		return "", nil
	}

	// "os" function is a custom function provided by Cloney, which returns the user's operating system.
	// This function has no parameters. It is useful for generating OS-specific parts.
	funcMap["os"] = func() (string, error) {
		return runtime.GOOS, nil
	}

	// "arch" function is a custom function provided by Cloney, which returns the user's operating system architecture.
	// This function has no parameters. It is useful for generating OS-specific parts.
	funcMap["arch"] = func() (string, error) {
		return runtime.GOARCH, nil
	}

	// Implement your own custom functions here...

	return funcMap
}
