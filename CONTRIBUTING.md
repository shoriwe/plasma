# Contributing

## Steps for contributing to this project

1. Fork the repository
2. Clone your fork
3. Make the changes to want
4. Create a pull request
5. Feedback loop until approval or revoke

## Language

Despite anyone native language, when contributing to the project, comments, function names, modules names, variables
names and any related word that could be written should be always be in `english` this because the easy of the language
and to maintain uniformity in the entire project.

## Important

- Always avoid hundreds of commits in your pull request, since it will take longer to review and approve.
- Avoid solving multiple issues in a single pull requests, since it will be easy to isolate the problem and the review of your solution.

## Changelog

Always include in your pull request in the section of changelog a list of those things you have modified, added or
removed

Example:

```
- Fixed mutex lock inside Value GetString
...
```

## Code

Everyone codes different but to maintain a bit of uniformity follow this rules in your coding process.

- ALWAYS use variable names that make sense.

Avoid doing:

```go
package my_package

import "github.com/shoriwe/gplasma/pkg/vm"

func f(plasma *vm.VM) (*vm.Value, error) { // Avoid cryptic function names
	x := vm.NewValue() // Avoid cryptic function variables
}
```