# Cloney: Streamlining Git Template Customization

<br>
<p align="center">
  <img src="images/cloney-logo.png">
</p>
<br>


## Introduction

Have you ever used a template Git repository and found yourself in the tedious task of replacing values manually or making extensive adjustments to fit your specific needs? If you have, you're not alone. Traditional Git templates often leave you with the burden of customizing every detail, which can be time-consuming and error-prone. This is where Cloney comes into play, revolutionizing the way you work with Git repositories.

### The Pain of Manual Adjustments

Imagine you've stumbled upon a fantastic template Git repository on GitHub that promises to kickstart your project. Excited, you fork the repository or download the ZIP archive, only to realize that it's not quite ready for your unique requirements. You need to replace placeholder values, tweak configurations, and adapt the code to match your project's specifications.

This process can be both frustrating and error-prone. Manually searching and replacing values throughout the codebase can lead to mistakes and inconsistencies. What if you could automate this entire customization process and have a template repository that adapts itself to your needs effortlessly? This is precisely what Cloney is designed to do.

### Cloney: Redefining Git Templates

Cloney is not just another Git template manager; it's a unique tool that redefines how you work with template repositories. With Cloney, you can say goodbye to manual adjustments and hello to dynamic template creation and management.

### Cloney Template Repositories

Cloney Template Repositories are the foundation of Cloney's innovative approach. These repositories are enriched with a special `.cloney.yaml` metadata file that contains vital information about the template repository.

### Custom Variables

Cloney empowers you to define variables within your templates. These variables act as placeholders for values that can be customized during the cloning process.

## What Makes a Cloney Template Repository?

A Cloney Template Repository consists of the following components:

- **Git Repository**: The core of your project, containing all the files, directories, and code that you want to share as a template.

- **`.cloney.yaml` Metadata File**: This special YAML file serves as the template repository's control center. It stores essential details about the repository, such as variable definitions, descriptions, and other crucial information that Cloney relies on to generate dynamic templates.

By leveraging the `.cloney.yaml` metadata file, you enable users to customize the template during the cloning process, creating a unique project that adapts to specific requirements.

**Understanding `.cloney.yaml` Metadata**

To harness the full potential of Cloney Template Repositories, it's crucial to grasp the structure and content of the `.cloney.yaml` metadata file. This file serves as the blueprint for your template, defining its characteristics, variables, and default settings.

## Template Information

The `.cloney.yaml` metadata file begins with essential information about your template repository:

- **Name**: The name of your template, providing a clear identifier for users.

- **Description**: A brief but informative description of your template's purpose and functionality.

- **Authors**: A list of contributors or creators of the template, acknowledging their role in its development.

- **License**: The licensing information for your template, specifying how others can use and distribute it.

- **Template Version**: The version number of your template, allowing users to identify different releases.

- **Manifest Version**: The version of the Cloney manifest file used in the template, ensuring compatibility with Cloney's features.

### Template Variables

Within the `.cloney.yaml` metadata file, Cloney allows you to define variables that users can customize during the cloning process. These variables play a central role in enabling dynamic template creation, ensuring that users can tailor templates to their specific requirements.

- **Name**: Assign a unique identifier to each variable, making it easy for users to reference and customize.

- **Description**: Provide a clear and concise description of each variable's purpose. These descriptions guide users in understanding how each variable affects the template.

- **Default (Optional)**: Specify a default value for each variable. This value is used when a user doesn't provide a custom value during the cloning process. If this field is omitted, Cloney assumes that the variable is mandatory and must be informed by the user.

- **Example Value**: Every variable must include an example value that demonstrates how it should be formatted and used. This example serves as a practical reference for users, helping them correctly configure variables within their customized templates.

## Customizing Variables

When cloning a Cloney Template Repository, users can customize variables by providing values that align with the variable definitions in the `.cloney.yaml` metadata file. Here's an example of a `.cloney-vars.yaml` file that contains custom values for the variables:

```yaml
app_name: MyApp
enable_https: true
currencies:
  - "Real"
  - "US Dollar"
  - "Yene"
```

## Accessing Variables in Template Files

Cloney makes it effortless to access and utilize variables within your template files. It employs the Go template syntax, a powerful and flexible language for generating text and code.

To access and use variables within your template files, follow these steps:

1. **Enclose Variables with Double Curly Braces**: To indicate that a piece of text should be replaced with a variable's value, enclose the variable name within double curly braces. For example, `{{ .VariableName }}`.

2. **Use the Dot (`.`) to Access Variables**: In Go templates, you access variables by prefixing their names with a dot (`.`). This dot signifies the context in which the

 variable is defined.

**Example:**

Let's assume you have a Cloney variable named `app_name` and a variable named `currencies` defined in your `.cloney.yaml` metadata file. In your template files, you can use the variables as follows:

```go
appName := "{{ .app_name }}"
```

```go
currencies := []string{
  {{ range .currencies }}
    "{{ . }}",
  {{ end }}
}
```

When Cloney generates the customized template based on user input, it replaces `{{ .app_name }}` and `{{ .currencies }}` with the corresponding values, resulting in dynamic and adaptable template files.

### Go Template Tutorials

To help you make the most of Cloney's dynamic variables and the Go template syntax, we recommend exploring tutorials and documentation on Go templates. Go templates are a widely used tool for generating text and are well-documented within the Go programming language. You can find these resources at:

- [Go Template Package Documentation](https://pkg.go.dev/text/template): Official Go documentation for the `text/template` package, which provides an in-depth look at Go templates.

- [Go Text Templates](https://golang.org/pkg/text/template/): The official Go documentation on text templates.

- [Sprig Functions](https://masterminds.github.io/sprig/): Cloney includes the Sprig library, which adds a variety of useful functions to Go templates. Explore the Sprig documentation to take advantage of these functions in your Cloney templates.

With these resources, you'll be well-equipped to create dynamic and adaptable template files that make the most of Cloney's variable-driven customization. Cloney simplifies the process of turning generic templates into tailored solutions, helping you save time and reduce errors in your Git-based projects.
