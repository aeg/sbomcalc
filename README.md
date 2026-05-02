# sbomcalc

`sbomcalc` is a CLI tool for set operations on SBOM components.

It reads SPDX JSON and CycloneDX JSON files, treats components as sets, and prints query, diff, and changed results. The current scope is component-level comparison only. Dependency graphs, SPDX relationships, CycloneDX dependencies, vulnerabilities, and semantic version comparison are intentionally out of scope for v0.1.

[日本語版 README](README_ja.md)

## Features

- Query SBOM component sets with `and`, `or`, `minus`, and `xor`.
- Compare SBOMs with `diff`.
- Show only version changes with `changed`.
- Select component identity level:
  - L1: `name`
  - L2: `name + version`
- Read SBOM files with streaming JSON decoding.
- Write `table`, `txt`, `json`, SPDX JSON, and CycloneDX JSON output.

## Supported Formats

Input:

- SPDX JSON 2.2
- SPDX JSON 2.3
- CycloneDX JSON 1.5
- CycloneDX JSON 1.6
- CycloneDX JSON 1.7

Output:

- `table`
- `txt`
- `json`
- `spdx-json`
- `spdx-json@2.2`
- `spdx-json@2.3`
- `cyclonedx-json`
- `cyclonedx-json@1.5`
- `cyclonedx-json@1.6`
- `cyclonedx-json@1.7`

`spdx-json` defaults to `spdx-json@2.3`.
`cyclonedx-json` defaults to `cyclonedx-json@1.7`.

## Install

From source:

```bash
go install github.com/eiji/sbomcalc/cmd/sbomcalc@latest
```

For local development:

```bash
go build -o sbomcalc ./cmd/sbomcalc
```

## Usage

```bash
sbomcalc query [-l1|-l2] "EXPR" [-o FORMAT[=FILE] ...]
sbomcalc diff old.json new.json [-o FORMAT[=FILE] ...]
sbomcalc changed old.json new.json [-o FORMAT[=FILE] ...]
```

If `-o` is omitted, `table` is written to stdout.

### Query

```bash
sbomcalc query -l1 "a.json and b.json"
sbomcalc query -l2 "(a.json and b.json) minus c.json"
sbomcalc query -l2 "new.json minus old.json" -o table -o cyclonedx-json@1.7=added.cdx.json
```

Operators:

| Operator | Meaning |
| --- | --- |
| `and` | intersection |
| `or` | union |
| `minus` | difference |
| `xor` | symmetric difference |

Parentheses are supported. Operators at the same level are evaluated left to right.

Symbol operators such as `&`, `|`, `-`, and `^` are not supported in v0.1.

### Levels

L1 uses only component names:

```text
openssl
curl
zlib
```

L2 uses component names and versions:

```text
openssl@1.1.1
curl@7.81.0
```

The default level is L2.

SBOM format output is supported only for `query -l2`. `query -l1` can write `table`, `txt`, and `json`.

### Diff

```bash
sbomcalc diff old.json new.json
```

`diff` reports:

- `added`: names present only in the new SBOM
- `removed`: names present only in the old SBOM
- `changed`: names present in both SBOMs but with different version sets
- `unchanged`: names present in both SBOMs with the same version set

### Changed

```bash
sbomcalc changed old.json new.json
```

`changed` prints only names whose version sets differ between the old and new SBOM.

## Output

Write table output to stdout:

```bash
sbomcalc query -l2 "new.json minus old.json" -o table
```

Write JSON output to a file:

```bash
sbomcalc diff old.json new.json -o json=result.json
```

Write multiple outputs:

```bash
sbomcalc query -l2 "new.json minus old.json" \
  -o table \
  -o cyclonedx-json@1.7=added.cdx.json
```

Only one output may target stdout. Reusing the same output file in one command is an error.

## Examples

Using the test data in this repository:

```bash
go run ./cmd/sbomcalc query --l1 "testdata/old.spdx.json and testdata/new.cdx.json"
```

Output:

```text
NAME
curl
openssl
```

```bash
go run ./cmd/sbomcalc diff testdata/old.spdx.json testdata/new.cdx.json
```

Output:

```text
ADDED
  nginx@1.24.0

REMOVED
  log4j@2.14.1

CHANGED
  openssl
    old: 1.1.1
    new: 3.0.0

UNCHANGED
  curl@7.81.0
```

The repository also includes more complex version-set fixtures:

```bash
go run ./cmd/sbomcalc diff testdata/complex-old.spdx.json testdata/complex-new.cdx.json
```

These fixtures include names whose version sets have a shared version, a removed version, and an added version.

## Notes

- stdin input is not supported. Inputs must be file paths.
- File paths with spaces are not supported in query expressions in v0.1.
- Quoted file paths are not supported in query expressions in v0.1.
- Generated SBOM output is a new minimal SBOM. Input metadata, relationships, dependencies, and vulnerabilities are not preserved.
- Empty component names are ignored.
- Empty versions are allowed and treated as `""`.

## Development

Run tests:

```bash
go test ./...
```

If the default Go build cache is not writable in your environment:

```bash
GOCACHE=/tmp/sbomcalc-go-build go test ./...
```

Run the CLI locally:

```bash
go run ./cmd/sbomcalc --help
```
