# Cloney

<br>
<p align="center">
  <img src="images/cloney-logo.png">
</p>
<br>

## Table of Contents

- 🧬 [Introduction](#introduction)

  - [The Pain of Manual Adjustments](#the-pain-of-manual-adjustments)
  - [Cloney Template Repositories](#cloney-template-repositories)
  - [Custom Variables](#custom-variables)

- 🧬 [What Makes a Cloney Template Repository?](#what-makes-a-cloney-template-repository)

- 🧬 [Cloney Example Template Repository](#cloney-example-template-repository)

- 🧬 [Understanding the `.cloney.yaml` Metadata File](#understanding-the-cloneyyaml-metadata-file)

  - [Template Variables](#template-variables)

- 🧬 [Providing Custom Variables](#providing-custom-variables)

- 🧬 [Accessing Variables in Template Files](#accessing-variables-in-template-files)

- 🧬 [Dynamically Generating Files](#dynamically-generating-files)

- 🧬 [Go Template Tutorials](#go-template-tutorials)

- 🧬 [Cloney Command-Line Interface (CLI)](#cloney-command-line-interface-cli)
  - [Cloney CLI Commands](#cloney-cli-commands)

## Introduction

Have you ever used a template Git repository and found yourself in the tedious task of replacing values manually or making extensive adjustments to fit your specific needs? If you have, you're not alone. Traditional Git templates often leave you with the burden of customizing every detail, which can be time-consuming and error-prone. This is where Cloney comes into play, revolutionizing the way you work with Git repositories.

### The Pain of Manual Adjustments

Imagine you've stumbled upon a fantastic template Git repository on GitHub that promises to kickstart your project. Excited, you fork the repository only to realize that it's not quite ready for your unique requirements. You need to replace placeholder values, tweak configurations, and adapt the code to match your project's specifications.

This process can be both frustrating and error-prone. Manually searching and replacing values throughout the codebase can lead to mistakes and inconsistencies. What if you could automate this entire customization process and have a template repository that adapts itself to your needs effortlessly? This is precisely what Cloney is designed to do, with it, you can say goodbye to manual adjustments and hello to dynamic template creation and management.

### Cloney Template Repositories

Cloney Template Repositories are the foundation of Cloney's innovative approach. These repositories are enriched with a special `.cloney.yaml` metadata file that contains information about the template repository.

### Custom Variables

Cloney empowers you to define variables within your templates. These variables act as placeholders for values that can be customized during the cloning process.

## What Makes a Cloney Template Repository?

A Cloney Template Repository consists of the following components:

- **Git Repository**: The core of your project, containing all the files, directories, and code that you want to share as a template.

- **`.cloney.yaml` Metadata File**: This special YAML file serves as the template repository's control center and is located in the root of the repository. It stores essential details about the repository, such as variable definitions, descriptions, and other crucial information that Cloney relies on to generate dynamic templates.

By adding the `.cloney.yaml` file to your repository's root, you effectively transform your standard repository into a Cloney Template Repository. This simple addition unlocks the power of dynamic template creation and customization, allowing users to tailor the project to their specific needs with ease.

## Cloney Example Template Repository

Visit the [Cloney Example Template Repository](https://github.com/ArthurSudbrackIbarra/cloney-test) to see Cloney templates in action. Here, you can:

- Explore its structure.
- Check out its `.cloney.yaml` metadata file.
- Customize variables to see Cloney in action.

Discover the simplicity of project setup and customization by exploring this example repository.

## Understanding the `.cloney.yaml` Metadata File

To harness the full potential of Cloney Template Repositories, it's crucial to grasp the structure and content of the `.cloney.yaml` metadata file. This file serves as the blueprint for your template, defining its characteristics, variables, and default settings. Below is an illustrative example of a `.cloney.yaml` file:

```yaml
manifest_version: v1

name: Bank Site
description: A template to create a site for a bank.
authors:
  - Chico Buarque
  - Freddie Mercury
  - John Lennon
license: MIT
template_version: "1.1.0"

configuration:
  ignore_paths:
    - development/*.html
    - images/*.jpg

variables:
  - name: app_name
    description: The name of your application.
    default: my_app
    example: my_app

  - name: enable_https
    description: Wheter to enable HTTPS or not.
    example: true

  - name: currencies
    description: List of currencies to use.
    example:
      - Real
      - US Dollar
      - Yen
```

- **Manifest Version**: The version of the Cloney manifest file used in the template, ensuring compatibility with different versions of Cloney.

- **Name**: The name of your template, providing a clear identifier for users.

- **Description**: A brief but informative description of your template's purpose and functionality.

- **Authors**: A list of contributors or creators of the template, acknowledging their role in its development.

- **License**: The licensing information for your template, specifying how others can use and distribute it.

- **Template Version**: The version number of your template, allowing users to identify different releases.

- **Variables**: A list of variables that users can customize during the cloning process.

### Template Configuration

The `.cloney.yaml` metadata file also allows you to configure some aspects of your template. These configurations are optional and can be omitted if you don't require them.

- **Ignore Paths**: A list of paths to ignore when cloning the template. This is useful for excluding files that are not relevant to the template's customization process.

### Template Variables

Within the `.cloney.yaml` metadata file, Cloney allows you to define variables that users can customize during the cloning process. These variables play a central role in enabling dynamic template creation, ensuring that users can tailor templates to their specific requirements.

- **Name**: Assign a unique identifier to each variable, making it easy for users to reference and customize.

- **Description**: Provide a clear and concise description of each variable's purpose. These descriptions guide users in understanding how each variable affects the template.

- **Default (Optional)**: Specify a default value for each variable. This value is used when a user doesn't provide a custom value during the cloning process. If this field is omitted, Cloney assumes that the variable is mandatory and must be informed by the user.

- **Example Value**: Every variable must include an example value that demonstrates how it should be formatted and used. This example serves as a practical reference for users, helping them correctly configure variables within their customized templates.

## Providing Custom Variables

When cloning a Cloney Template Repository, users can customize variables by providing values that align with the variable definitions in the `.cloney.yaml` metadata file. Here's an example of a `.cloney-vars.yaml` file that contains custom values for the variables:

```yaml
app_name: My App
enable_https: true
currencies:
  - Real
  - US Dollar
  - Yen
```

## Accessing Variables in Template Files

Cloney makes it effortless to access and utilize variables within your template files. It employs the Go template syntax, a powerful and flexible language for generating text and code. It's essential to note that you **don't require any knowledge of Go programming** to utilize Cloney effectively. The term "Go template" is simply a historical reference to the language where these templates initially emerged.

To access and use variables within your template files, follow these steps:

1. **Enclose Variables with Double Curly Braces**: To indicate that a piece of text should be replaced with a variable's value, enclose the variable name within double curly braces. For example, `{{ .VariableName }}`.

2. **Use the Dot (`.`) to Access Variables**: In Go templates, you access variables by prefixing their names with a dot (`.`). This dot signifies the context in which the variable is defined.

**Example:**

Imagine you're working with a Cloney Template Repository for a billing system. In this template, you want to create a customized HTML file based on various parameters. Let's explore how Cloney's variables can make this process effortless.

Here's a snippet from your template HTML file, `index.html`:

```html
<!-- File: index.html -->

<h1>
  {{ .app_name }} is a fantastic app that supports the following currencies:
</h1>
<ul>
  {{- range .currencies }}
  <li>{{ . }}</li>
  {{- end }}
</ul>

{{- if .enable_https }}
<p>HTTPS is enabled.</p>
{{- end }}
```

Now, let's break down what's happening here:

- `{{ .app_name }}`: This placeholder will be replaced with the user-defined name of the application when you generate your customized template.

- `{{- range .currencies }} ... {{- end }}`: This construct iterates through the list of currencies provided by the user during the cloning process. It dynamically generates list items for each currency.

- `{{- if .enable_https }} ... {{- end }}`: This conditional block checks if the user has chosen to enable HTTPS during cloning. If so, it includes a paragraph indicating that HTTPS is enabled.

When you use Cloney to create your template based on user input, it automatically populates these placeholders with the values you specified during the cloning process. This results in a tailored HTML file that perfectly suits your requirements.

Here's an example of what the generated HTML could look like:

```html
<!-- File: index.html -->

<h1>MyApp is a fantastic app that supports the following currencies:</h1>
<ul>
  <li>Real</li>
  <li>US Dollar</li>
  <li>Yen</li>
</ul>

<p>HTTPS is enabled.</p>
```

**PS**: It's important to mention that, although this example uses HTML, you can use Cloney variables in any type of file. This includes YAML, JSON, Markdown, and even source code files.

## Dynamically Generating Files

Cloney, allows you to not only replace placeholders with variable values but also generate files and directories dynamically. This feature is particularly useful for creating multiple files based on a template and data. Let's dive into the process using a practical example.

### Example Scenario

Suppose you have a YAML file containing currency information like this:

```yaml
app_name: MyApp

currencies:
  - name: USD
    symbol: $
    description: "United States Dollar"
  - name: EUR
    symbol: €
    description: "Euro"
  - name: GBP
    symbol: £
    description: "British Pound"
```

Your goal is to create a home page that lists all the currencies and links to a dedicated page for each currency.

### Creating the `home.html` Template

You can start by creating a `home.html` file that dynamically generates the content for each currency and generates links to their respective pages.

```html
<!-- File: home.html -->

{{- define "currency-file-content" -}}
<!-- File: {{ .name }}.html -->

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Currency: {{ .name }}</title>
</head>
<body>
    <h1>Currency: {{ .name }}</h1>
    <p>Symbol: {{ .symbol }}</p>
    <p>Description: {{ .description }}</p>
</body>
</html>
{{- end -}}

<p>Check out the different currencies:</p>

<ul>
{{- range .currencies }}
{{- $fileName := printf "%s.html" .name }}
{{- toFile $fileName "currency-file-content" . }}
  <li>
    <a href="{{ $fileName }}">{{ .name }}</a>
  </li>
{{- end }}
</ul>
```

Let's break down what's happening within this template:

1. `{{- define "currency-file-content" -}} ... {{- end -}}`: This defines a template named `currency-file-content` that will be used to generate the content of each currency's page.

2. `{{- range .currencies }} ... {{- end }}`: This iterates through the list of currencies provided in your YAML data.

3. `{{- $fileName := printf "%s.html" .name }}`: This creates a variable named `fileName` to store the name of the file to be generated. The file name will be constructed using the currency's name with the `.html` extension.

4. `{{- toFile $fileName "currency-file-content" . }}`: This generates a file named `fileName` using the `currency-file-content` template and the current currency as the context. This is where the dynamic file generation happens.

### Generated Output

When Cloney processes this template, it generates several files. The `home.html` file will have the following content:

```html
<!-- File: home.html -->

<p>Check out the different currencies:</p>
<ul>
  <li>
    <a href="USD.html">USD</a>
  </li>
  <li>
    <a href="EUR.html">EUR</a>
  </li>
  <li>
    <a href="GBP.html">GBP</a>
  </li>
</ul>
```

Additionally, Cloney generates three other files, namely `USD.html`, `EUR.html`, and `GBP.html`, each containing content specific to the respective currency:

```html
<!-- File: USD.html -->

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Currency: USD</title>
</head>
<body>
    <h1>Currency: USD</h1>
    <p>Symbol: $</p>
    <p>Description: United States Dollar</p>
</body>
</html>
```

## Known Limitations and Workarounds

While Cloney's dynamic file generation is a powerful feature, there are certain limitations and workarounds to consider when working with it.

### Limitation: Accessing Root-Level Variables

In the example provided, accessing root-level variables like `app_name` from within the `currency-file-content` template can be challenging. This is because the context within the template is limited to the data passed to it, and it does not have direct access to variables defined outside the template.

### Workaround: Using a Dictionary

To access root-level variables within a dynamically generated template, you can use a dictionary to store both the root-level variables and the parameters needed for the `currency-file-content` template. Here's how to do it:

```html
{{- $dict := dict "root" $ "param" . -}}
{{- toFile $fileName "currency-file-content" $dict -}}
```

In this code, a dictionary named `$dict` is created, which contains two key-value pairs: `"root"` and `"param"`. `"root"` is assigned the value of the root-level context (`$`), and `"param"` is assigned the value of the current currency data (`.`).

### Accessing Variables in the Template

With the dictionary approach, you can now access root-level variables and parameters in the `currency-file-content` template as follows:

```html
{{- define "currency-file-content" -}}
<p>The app name is: {{ .root.app_name }}</p>
<p>The currency name is: {{ .param.name }}</p>
{{- end -}}
```

By using this workaround, you can access both root-level and local variables within your dynamically generated files.

## Go Template Tutorials

To help you make the most of Cloney's dynamic variables and the Go template syntax, we recommend exploring tutorials and documentation on Go templates. Go templates are a widely used tool for generating text and are well-documented within the Go programming language. You can find these resources at:

- **Official Go Documentation**: Dive into the official Go documentation for the `text/template` package, which provides an in-depth look at Go templates. Explore the intricacies of template creation and customization.

  - [Go Template Package Documentation](https://pkg.go.dev/text/template)

- **Go Text Templates**: Check out the official Go documentation dedicated to text templates. It offers comprehensive insights into working with Go templates, including advanced topics and best practices.

  - [Go Text Templates](https://golang.org/pkg/text/template/)

- **Sprig Functions**: Cloney leverages the Sprig library, which extends Go templates with a variety of useful functions. Familiarize yourself with these functions to enhance the capabilities of your Cloney templates and streamline your projects.

  - [Explore Sprig Functions](https://masterminds.github.io/sprig/)

- **Helm Template Functions**: Helm is a widely-used package manager for Kubernetes, utilizing Go templates and introducing custom functions. Cloney integrates with Helm's `include` function, enabling its use within your templates.

  - [Explore Helm's Include Function](https://helm.sh/docs/howto/charts_tips_and_tricks/#:~:text=The%20include%20function%20allows%20you,results%20to%20other%20template%20functions.)

## Cloney Command-Line Interface (CLI)

The Cloney Command-Line Interface (CLI) is the tool you use to interact with Cloney Template Repositories. It offers a set of straightforward commands that enable you to manage, customize, and work with Cloney templates.

### Cloney CLI Commands

The `cloney` CLI is used to interact with Cloney Template Repositories, allowing you to clone, customize, and manage templates. Below, you'll find examples of various `cloney` commands:

1. **Start Command**

   - To create a new Cloney Template Repository in the current directory:

     ```bash
     cloney start
     ```

    <br>
    <p align="center">
      <img src="images/commands/start.png">
    </p>
    <br>
    <br>

1. **Info Command**

   - To retrieve information about a Cloney Template Repository in the current directory:

     ```bash
     cloney info
     ```

   - To get information about a local Cloney Template Repository at a specific path:

     ```bash
     cloney info /path/to/local/repository
     ```

   - To fetch information about a remote Cloney Template Repository using its URL:

     ```bash
     cloney info https://github.com/username/template-repo.git
     ```

     <br>
     <p align="center">
       <img src="images/commands/info.png">
     </p>
     <br>
     <br>

1. **Clone Command**

   - To clone a Cloney Template Repository with custom variable values from a variables file (e.g., `my-vars.yaml`) in the same directory:

     ```bash
     cloney clone https://github.com/username/template-repo.git -v my-vars.yaml
     ```

     If you don't specify a variables file using `-v`, Cloney will search for a file named `.cloney-vars.yaml` in the current directory by default.

     You can also specify variables using inline YAML:

     ```bash
     cloney clone https://github.com/username/template-repo.git -v '{ app_name: "MyApp", enable_https: true }'
     ```

    <br>
    <p align="center">
      <img src="images/commands/clone.png">
    </p>
    <br>
    <br>

1. **Dry Run Command**

   - To debug your own Cloney Template Repository, seeing how variables are replaced, and checking for errors:

     ```bash
     cloney dry-run -v my-vars.yaml
     ```

   Just like the `clone` command, you can pass variables in the same way.

    <br>
    <p align="center">
      <img src="images/commands/dry-run.png">
    </p>
    <br>
    <br>
