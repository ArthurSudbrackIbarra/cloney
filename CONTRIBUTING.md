# Contributing to Cloney

Thank you for considering contributing to Cloney! We appreciate your interest in helping us improve this project. By contributing to Cloney, you not only make it better but also empower others in their Git repository management.

To contribute, follow these simple steps:

## Opening an Issue

If you encounter a bug, have a feature request, or want to discuss something related to Cloney, feel free to open an issue. This is the first step to get your ideas or concerns addressed.

1. Go to the [Cloney GitHub repository](https://github.com/arthursudbrackibarra/cloney).
2. Click on the "Issues" tab.
3. Click the "New Issue" button.
4. Fill out the issue template, providing as much detail as possible.
5. Submit the issue.

## Making Changes

If you'd like to contribute code or documentation, follow these steps to create a pull request:

1. Ensure you have Go version 1.21 or higher installed on your system. You can download it from the [official Go website](https://golang.org/dl/).

2. Clone the [Cloney repository](https://github.com/arthursudbrackibarra/cloney) to your local machine:

   ```sh
   git clone https://github.com/arthursudbrackibarra/cloney.git
   ```

3. Create a new branch for your changes:

   ```sh
   git checkout -b feature/your-feature
   ```

4. Make your changes and commit them:

   ```sh
   git commit -m "Add your commit message here"
   ```

5. Open a pull request on the [Cloney repository](https://github.com/arthursudbrackibarra/cloney) with a clear title and description of your changes.

## Testing the Cloney CLI

Before submitting your pull request, it's important to test the Cloney CLI to ensure your changes work as expected. You can do this by building the Cloney executable:

1. Navigate to the Cloney repository directory on your local machine:

   ```sh
   cd cloney
   ```

2. Build the Cloney executable:

   ```sh
   go build -o cloney # for Linux and macOS

   go build -o cloney.exe # for Windows
   ```

3. You can now test the Cloney CLI by running:

   ```sh
   ./cloney # for Linux and macOS

   ./cloney.exe # for Windows
   ```

   This will display the available commands and options, allowing you to verify that your changes haven't introduced any issues.

## Documentation Contributions

If your contribution includes changes to Cloney's code that introduce new features or changes that require documentation, please make sure to update the official Cloney documentation hosted at [Cloney Documentation](https://arthursudbrackibarra.github.io/cloney-documentation/). This ensures that users can easily access the latest information about the new features or changes.

## Code of Conduct

Please note that Cloney has a [Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project, you agree to abide by its terms.

## Questions and Support

If you have questions or need support with Cloney, feel free to reach out to us via [GitHub Discussions](https://github.com/arthursudbrackibarra/cloney/discussions) or the [Official Cloney Documentation](https://arthursudbrackibarra.github.io/cloney-documentation/).

Your contributions are valuable and help us make Cloney even better. Thank you for being a part of this project!
