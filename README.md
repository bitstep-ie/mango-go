# mango-go
[![CodeQL](https://github.com/bitstep-ie/mango-go/actions/workflows/codeql.yml/badge.svg)](https://github.com/bitstep-ie/mango-go/actions/workflows/codeql.yml)
[![Dependabot](https://github.com/bitstep-ie/mango-go/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/bitstep-ie/mango-go/actions/workflows/dependabot/dependabot-updates)
[![codecov](https://codecov.io/github/bitstep-ie/mango-go/graph/badge.svg?token=L6EJH29N5L)](https://codecov.io/github/bitstep-ie/mango-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitstep-ie/mango-go)](https://goreportcard.com/report/github.com/bitstep-ie/mango-go)

mango-go is a collection of utility packages that can be used in other go projects.

## Packages 

Comprised of the following packages:
- compare
- env
- io
- mango_logger
- testutils
- time


## Documentation




## Consider contributing

Do you have a useful function? Do you have something that could be useful to others?

Yes please! [See our starting guidelines](contributing.md). Your help is very welcome!


## Make?

You can use `make all` to ensure all the checks are performed before you push the code on a remote branch and await github actions to run the CI.

The makefile is to help with local development of the library by giving you the exact steps that the ci will execute.

This makefile will NOT be used as part of builds.

Currently, mimics the mango CI steps that can be found here:

It is up to you if you deviate from the github actions, do so at your own risk and should not be committed back into the project.


## Structure

Try to keep common functions in existing packages, and follow the same pattern if required to create a new package. Update documentation as required, and make sure to note any breaking changes clearly in the PR and ofc debate on the mango teams channel if it requires and how to handle version increase.

### Filename convention 
For basic packages try to match to this convention:
`smallCaseUtils.go`
`smallCaseUtils_test.go`

Larger packages requiring multiple files (e.g: `mclogger`), has no current structure convention. 

### Gremlins coverage

Both `efficacy` & `mutant-coverage` sit at `95%`. Aim for this or higher, the build will fail if these thresholds are not met. Under committee review if necessary these thresholds will be reviewed.
