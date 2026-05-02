# sbomcalc

`sbomcalc` は、SBOMのコンポーネントに対して集合演算を行うCLIツールです。

SPDX JSONとCycloneDX JSONを読み取り、コンポーネントを集合として扱い、query、diff、changedの結果を出力します。現在の対象はコンポーネント単位の比較のみです。依存関係グラフ、SPDX relationships、CycloneDX dependencies、脆弱性、semantic version比較は v0.1 では扱いません。

[English README](README.md)

## 主な機能

- `and`、`or`、`minus`、`xor` によるSBOMコンポーネント集合の問い合わせ
- `diff` によるSBOM間の比較
- `changed` によるバージョン変更のみの表示
- コンポーネント識別レベルの指定
  - L1: `name`
  - L2: `name + version`
- ストリーミングJSONデコードによるSBOM読み取り
- `table`、`txt`、`json`、SPDX JSON、CycloneDX JSON出力

## 対応形式

入力:

- SPDX JSON 2.2
- SPDX JSON 2.3
- CycloneDX JSON 1.5
- CycloneDX JSON 1.6
- CycloneDX JSON 1.7

出力:

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

`spdx-json` は `spdx-json@2.3` として扱います。
`cyclonedx-json` は `cyclonedx-json@1.7` として扱います。

## インストール

ソースからインストールする場合:

```bash
go install github.com/aeg/sbomcalc/cmd/sbomcalc@latest
```

ローカルでビルドする場合:

```bash
go build -o sbomcalc ./cmd/sbomcalc
```

## 使い方

```bash
sbomcalc query [-l1|-l2] "EXPR" [-o FORMAT[=FILE] ...]
sbomcalc diff old.json new.json [-o FORMAT[=FILE] ...]
sbomcalc changed old.json new.json [-o FORMAT[=FILE] ...]
```

`-o` を省略した場合は、`table` を標準出力に出力します。

### query

```bash
sbomcalc query -l1 "a.json and b.json"
sbomcalc query -l2 "(a.json and b.json) minus c.json"
sbomcalc query -l2 "new.json minus old.json" -o table -o cyclonedx-json@1.7=added.cdx.json
```

演算子:

| 演算子 | 意味 |
| --- | --- |
| `and` | 積集合 |
| `or` | 和集合 |
| `minus` | 差集合 |
| `xor` | 対称差 |

括弧を使えます。同じ階層の演算子は左から右に評価します。

`&`、`|`、`-`、`^` のような記号演算子は v0.1 では対応していません。

### レベル

L1はコンポーネント名だけを使います。

```text
openssl
curl
zlib
```

L2はコンポーネント名とバージョンを使います。

```text
openssl@1.1.1
curl@7.81.0
```

デフォルトはL2です。

SBOM形式の出力は `query -l2` のみ対応しています。`query -l1` では `table`、`txt`、`json` を出力できます。

### diff

```bash
sbomcalc diff old.json new.json
```

`diff` は以下を出力します。

- `added`: 新しいSBOMにだけ存在する名前
- `removed`: 古いSBOMにだけ存在する名前
- `changed`: 両方のSBOMに存在するが、バージョン集合が異なる名前
- `unchanged`: 両方のSBOMに存在し、バージョン集合も同じ名前

### changed

```bash
sbomcalc changed old.json new.json
```

`changed` は、古いSBOMと新しいSBOMでバージョン集合が異なる名前だけを出力します。

## 出力

tableを標準出力に出力する例:

```bash
sbomcalc query -l2 "new.json minus old.json" -o table
```

JSONをファイルに出力する例:

```bash
sbomcalc diff old.json new.json -o json=result.json
```

複数の出力を指定する例:

```bash
sbomcalc query -l2 "new.json minus old.json" \
  -o table \
  -o cyclonedx-json@1.7=added.cdx.json
```

標準出力に出力できる指定は1つだけです。同じコマンド内で同じ出力ファイルを複数回指定するとエラーになります。

## 例

このリポジトリのテストデータを使う例:

```bash
go run ./cmd/sbomcalc query --l1 "testdata/old.spdx.json and testdata/new.cdx.json"
```

出力:

```text
NAME
curl
openssl
```

```bash
go run ./cmd/sbomcalc diff testdata/old.spdx.json testdata/new.cdx.json
```

出力:

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

このリポジトリには、より複雑なバージョン集合を持つテストデータも含まれています。

```bash
go run ./cmd/sbomcalc diff testdata/complex-old.spdx.json testdata/complex-new.cdx.json
```

このテストデータでは、同じ名前の中に共通バージョン、削除されたバージョン、追加されたバージョンが同時に含まれるケースを確認できます。

## 注意事項

- 標準入力からの読み取りには対応していません。入力はファイルパスで指定します。
- v0.1では、query式内で空白を含むファイルパスは使えません。
- v0.1では、query式内でクォートしたファイルパスは使えません。
- 生成するSBOMは新規の最小SBOMです。入力SBOMのメタデータ、relationships、dependencies、脆弱性は引き継ぎません。
- 空のコンポーネント名は無視します。
- 空のバージョンは許容し、`""` として扱います。

## 開発

テストを実行します。

```bash
go test ./...
```

環境によってGoの標準ビルドキャッシュに書き込めない場合:

```bash
GOCACHE=/tmp/sbomcalc-go-build go test ./...
```

ローカルでCLIを実行します。

```bash
go run ./cmd/sbomcalc --help
```
