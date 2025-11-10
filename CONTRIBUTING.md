# Contributing

We welcome both issue reports and pull requests! Please follow these guidelines to help maintainers respond effectively.

## 🧩 Issues

- **Before opening a new issue:**
  - Use the search tool to check for existing issues or feature requests.
  - Review existing issues and provide feedback or react to them.
  - Use English for all communications — it is the language all maintainers read and write.
  - For questions, configuration or deployment problems, please use the [Discussions Forum](https://github.com/bitstep-ie/mango-go/discussions).

- **Reporting a bug:**
  - Please provide a clear description of your issue, and a minimal reproducible code example if possible.
  - Include the mango-go version (or commit reference), Go version, and operating system.
  - Indicate whether you can reproduce the bug and describe steps to take.
  - Attach relevant logs.

- **Feature requests:**
  - Before opening a request, check that a similar idea hasn’t already been suggested.
  - Clearly describe your proposed feature and its benefits.

## 📥 Pull Requests

Please ensure your pull request meets the following requirements:

- Open your pull request against the `main` branch.
- All tests pass in available continuous integration systems (e.g., GitHub Actions) - to help you with this there is the <a href="#make">make file</a>.
- Add or modify tests to cover your code changes.
- If your pull request introduces a new feature, document it in [`docs/mango-go.md`](docs/mango-go.md), not in the README.
- Follow the checklist in the [Pull Request Template](.github/PULL_REQUEST_TEMPLATE.md:1).

Thank you for contributing!


### <a id="make"></a>🔧 Make

You can use `make all` to ensure all the checks are performed before you push the code on a remote branch and open PR which will execute the github actions.

The makefile is to help with local development of the library by giving you the exact steps that the ci will execute.

This makefile will NOT be used as part of builds.

It is up to you if you deviate from the github actions, do so at your own risk and should not be committed back into the project.

### <a id="structure"></a>📐 Structure

Try to keep common functions in existing packages, and follow the same pattern if required to create a new package. Update documentation as required, and make sure to note any breaking changes clearly in the PR and ofc debate on the mango teams channel if it requires and how to handle version increase.

### <a id="filename-convention"></a>🔖 Filename convention
For basic packages try to match to this convention:
`smallCaseUtils.go`
`smallCaseUtils_test.go`

Larger packages requiring multiple files (e.g: `logger`), has no current structure convention. We're open to discussion

### <a id="gremlins-coverage"></a>👹 Gremlins coverage

Both `efficacy` & `mutant-coverage` sit at `95%`. Aim for this or higher, the build will fail if these thresholds are not met. Under committee review if necessary these thresholds will be reviewed.

