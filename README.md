# function-comment-linter

`funccomment` is a module plugin for `golangci-lint` that requires every Go
function and method to have a doc comment whose first line starts with its
name and contains a Chinese description. `init` functions are ignored.

The separate `constcomment` rule requires every constant to have a
name-prefixed comment with Chinese text. Enable it together with `funccomment`
when both checks are wanted.

The Chinese description may appear on a later comment line, so an English
suffix on the first line can be completed by a later Chinese line.

## Requirements

- Go 1.25.4 or newer
- `golangci-lint` 2.6.2 (the version in `.custom-gcl.yml`)
- A network connection that can download the tagged module version `v0.1.2`

Use the same Go version for the build and lint runs. The published plugin is
loaded from the `v0.1.2` Git tag in `.custom-gcl.yml`; do not add `path: .` to
that shared configuration. If the local Go cache was created by another Go
version, use a fresh writable cache, for example:

```bash
GOCACHE=/private/tmp/function-comment-linter-gocache golangci-lint custom -v
```

The command writes the custom binary as `./custom-gcl`; it is ignored by Git.

## Build and run in another project

Copy `.custom-gcl.yml` and `.golangci.example.yml` into the project that will
use this plugin. Rename the latter to `.golangci.yml`, then build the custom
binary:

```bash
golangci-lint custom -v
```

Run the generated binary:

```bash
./custom-gcl run --config .golangci.yml ./...
```

The `custom` entry in the configuration is required so the generated binary
knows how to register the module plugin.

The generated `custom-gcl` is platform-specific. Build it separately on each
machine and keep it out of version control.

## Local plugin development

This repository keeps a separate `.custom-gcl.local.yml` for testing the
working tree before publishing a tag. It uses `path: .` so local source edits
are loaded directly:

```bash
cp .custom-gcl.local.yml .custom-gcl.yml
golangci-lint custom -v
./custom-gcl run --config .golangci.yml ./...
```

After publishing a new tag, update the `version` in `.custom-gcl.yml` and in
the consuming projects.

## Examples

Valid comments include:

```go
// ProcessData 加载数据。
func ProcessData() {}

// pathAnalysisConcurrencyKey
// 构造用户路径分析并发租约使用的 Redis key。
func pathAnalysisConcurrencyKey() {}

/*
processData 加载数据。
*/
func processData() {}
```

These are reported:

```go
func processData() {}

// ProcessData loads the data.
func ProcessData() {}
```

The description must contain at least one Chinese character. A comment that
only contains the function name or only English text is reported.

Run the unit and analyzer tests with:

```bash
GOCACHE=/private/tmp/function-comment-linter-gocache go test ./...
```
