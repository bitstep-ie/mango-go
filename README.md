[![CodeQL](https://github.com/bitstep-ie/mango-go/actions/workflows/codeql.yml/badge.svg)](https://github.com/bitstep-ie/mango-go/actions/workflows/codeql.yml)
[![Dependabot](https://github.com/bitstep-ie/mango-go/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/bitstep-ie/mango-go/actions/workflows/dependabot/dependabot-updates)
[![codecov](https://codecov.io/github/bitstep-ie/mango-go/graph/badge.svg?token=L6EJH29N5L)](https://codecov.io/github/bitstep-ie/mango-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitstep-ie/mango-go)](https://goreportcard.com/report/github.com/bitstep-ie/mango-go)

<br />
<div align="center">
    <a href="https://github.com/bitstep-ie/mango-go">
    <picture>
        <source srcset="images/mango-with-text-black.png" media="(prefers-color-scheme: light)">
        <source srcset="images/mango-with-text-white.png" media="(prefers-color-scheme: dark)">
        <img src="images/mango-with-text-black.png" alt="mango Logo">
    </picture>
    </a>
    <h3 align="center">mango-go</h3>
    <p align="center">
        A collection of utility packages for go
        <br />
        <a href="#"><strong>📚 Explore the docs »</strong></a>
        <br />
        <br />
        <a href="#">🔎 View Examples</a>
        &middot;
        <a href="https://github.com/bitstep-ie/mango-go/issues/new?labels=bug&template=bug-report---.md">
            🐛 Report Bug
        </a>
        &middot;
        <a href="https://github.com/bitstep-ie/mango-go/issues/new?labels=enhancement&template=feature-request---.md">
            💡 Request Feature
        </a>
    </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>📂 Table of Contents</summary>
  <ol>
    <li>
        <a href="#📋-about-the-project">📋 About The Project</a>
    </li>
    <li>
        <a href="#🚀-getting-started">🚀 Getting Started</a>
        <ul>
            <li><a href="#🛠️-installation">🛠️ Installation</a></li>
            <li><a href="#📦-packages">📦 Packages</a></li>
        </ul>
    </li>
    <li>
        <a href="#👨‍💻-usage">👨‍💻 Usage</a>
    </li>
    <li>
        <a href="#📝-contributing">📝 Contributing</a>
        <ul>
            <li><a href="#🔧-make">🔧 Make</a></li>
            <li><a href="#📐-structure">📐 Structure</a></li>
            <li><a href="#🔖-filename-convention">🔖 Filename convention</a></li>
            <li><a href="#👹-gremlins-coverage">👹 Gremlins coverage</a></li>
        </ul>
    </li>
    <li>
        <a href="#📜-license">📜 License</a>
    </li>
    <li>
        <a href="#👏-acknowledgments">👏 Acknowledgments</a>
        <ul>
            <a href="#🤝-contributors">🤝 Contributors</a>
        </ul>
    </li>
  </ol>
</details>

## 📋 About the project

## 🚀 Getting started

### 🛠️ Installation


### 📦 Packages 

- compare
- env
- io
- mango_logger
- testutils
- time


## 👨‍💻 Usage




## 📝 Contributing

Do you have a useful function? Do you have something that could be useful to others?

Yes please! [See our starting guidelines](contributing.md). Your help is very welcome!

### 🔧 Make

You can use `make all` to ensure all the checks are performed before you push the code on a remote branch and open PR which will execute the github actions.

The makefile is to help with local development of the library by giving you the exact steps that the ci will execute.

This makefile will NOT be used as part of builds.

It is up to you if you deviate from the github actions, do so at your own risk and should not be committed back into the project.

### 📐 Structure

Try to keep common functions in existing packages, and follow the same pattern if required to create a new package. Update documentation as required, and make sure to note any breaking changes clearly in the PR and ofc debate on the mango teams channel if it requires and how to handle version increase.

### 🔖 Filename convention
For basic packages try to match to this convention:
`smallCaseUtils.go`
`smallCaseUtils_test.go`

Larger packages requiring multiple files (e.g: `mclogger`), has no current structure convention.

### 👹 Gremlins coverage

Both `efficacy` & `mutant-coverage` sit at `95%`. Aim for this or higher, the build will fail if these thresholds are not met. Under committee review if necessary these thresholds will be reviewed.


## 📜 License


## 👏 Acknowledgments


### 🤝 Contributors

Thanks goes to these wonderful people:

<table align="center">
  <tr>
    <td align="center"><a href="https://github.com/Ronan-L-OByrne"><img src="https://github.com/Ronan-L-OByrne.png?size=100" width="100px;" alt="Ronan"/><br /><sub><b>Ronan</b></sub></a></td>
    <td align="center"><a href="https://github.com/bencarroll1"><img src="https://github.com/bencarroll1.png?size=100" width="100px;" alt="Ben"/><br /><sub><b>Ben</b></sub></a></td>
  </tr>
</table>

