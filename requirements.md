# sbomcalc 仕様書 v0.1

## 目的

SBOMファイル群に対して、コンポーネント単位の集合演算を行うCLIツールをGoで実装する。

対象はまずコンポーネント単位のみとし、依存関係グラフは扱わない。

---

## 対象SBOM形式

入力対応:

- SPDX JSON 2.2
- SPDX JSON 2.3
- CycloneDX JSON 1.5
- CycloneDX JSON 1.6
- CycloneDX JSON 1.7

形式判定:

- SPDX JSONはトップレベルの `spdxVersion` で判定する。
  - `SPDX-2.2` を SPDX JSON 2.2 として扱う。
  - `SPDX-2.3` を SPDX JSON 2.3 として扱う。
- CycloneDX JSONはトップレベルの `bomFormat` と `specVersion` で判定する。
  - `bomFormat` が `CycloneDX` で、`specVersion` が `1.5`、`1.6`、`1.7` のいずれかであること。
- バージョンが対象外の場合は未対応SBOM形式としてエラーにする。

出力対応:

- table
- txt
- json
- spdx-json
- spdx-json@2.2
- spdx-json@2.3
- cyclonedx-json
- cyclonedx-json@1.5
- cyclonedx-json@1.6
- cyclonedx-json@1.7

デフォルト:

- `spdx-json` は `spdx-json@2.3`
- `cyclonedx-json` は `cyclonedx-json@1.7`

生成するSBOM:

- `query -l2` の結果だけをSBOM形式で出力できる。
- 生成するSBOMは新規の最小SBOMとし、入力SBOMのメタデータ、依存関係、relationship、vulnerabilityは引き継がない。
- 入力から取得できた `name`、`version`、`purl`、`supplier`、`licenses`、`hashes` だけをコンポーネント情報として反映する。
- 同じ `name` と `version` のコンポーネントが複数入力に存在する場合は、最初に読み取ったレコードを代表としてSBOM出力に使う。

---

## CLI

### query

```bash
sbomcalc query [-l1|-l2] "EXPR" [-o FORMAT[=FILE] ...]
```

例:

```bash
sbomcalc query -l1 "a.json and b.json"
sbomcalc query -l2 "(a.json and b.json) minus c.json"
sbomcalc query -l2 "new.json minus old.json" -o table -o cyclonedx-json@1.7=added.cdx.json
```

### diff

```bash
sbomcalc diff old.json new.json [-o FORMAT[=FILE] ...]
```

### changed

```bash
sbomcalc changed old.json new.json [-o FORMAT[=FILE] ...]
```

---

## レベル指定

### `-l1`

パッケージ名のみで集合化する。

```text
ComponentKey = name
```

例:

```text
openssl
curl
zlib
```

### `-l2`

パッケージ名 + バージョンで集合化する。

```text
ComponentKey = name + version
```

例:

```text
openssl@1.1.1
curl@7.81.0
```

### デフォルト

`-l2`

### 制約

`-l1` と `-l2` の同時指定は禁止。

---

## query 式文法

### 演算子

| 演算子     | 意味  |
| ------- | --- |
| `and`   | 積集合 |
| `or`    | 和集合 |
| `minus` | 差集合 |
| `xor`   | 対称差 |

記号演算子 `&`, `|`, `-`, `^` は v0.1 では対応しない。

### 括弧

`(` `)` をサポートする。

### ファイル名

式中にファイルパスを直接書く。

例:

```text
a.json
./sboms/a.spdx.json
../old/app.cdx.json
```

空白を含むファイル名は v0.1 では非対応でよい。

字句解析:

- 空白と括弧でトークンを区切る。
- `and`、`or`、`minus`、`xor` は演算子として扱うため、これらと完全一致するファイル名は式中で使えない。
- クォートは v0.1 では非対応とする。
- ファイルパス中の `@`、`.`、`/`、`_`、`-` は通常文字として扱う。

### 評価

括弧を優先する。
同一階層では左結合で評価する。

```text
a.json and b.json minus c.json
```

は以下と同じ。

```text
(a.json and b.json) minus c.json
```

---

## query の意味

```bash
sbomcalc query -l1 "a.json and b.json"
```

意味:

```text
L1(a.json) ∩ L1(b.json)
```

```bash
sbomcalc query -l2 "a.json minus b.json"
```

意味:

```text
L2(a.json) - L2(b.json)
```

同一キーの扱い:

- 集合演算では `ComponentKey` だけを見る。
- 同じ入力内に同一キーのコンポーネントが複数ある場合、集合上は1要素として扱う。
- 結果の詳細情報を出力するときは、式中に現れたファイルを左から順に再走査し、結果キーに一致した最初の `ComponentRecord` を代表として採用する。
- `query -l1` の結果では `name` だけを出力し、`version`、`purl`、`supplier`、`licenses`、`hashes` は出力しない。

---

## diff の意味

```bash
sbomcalc diff old.json new.json
```

A = old.json
B = new.json

以下を出力する。

```text
added:
  Bに存在し、Aに存在しないname

removed:
  Aに存在し、Bに存在しないname

changed:
  AとBの両方に存在するnameだが、version集合が異なるもの

unchanged:
  AとBの両方に存在し、version集合も同じもの
```

diff はL1/L2の単純集合演算ではなく、nameとversion集合を同時に見る専用処理とする。

version集合:

- `name` ごとに、正規化後の `version` の集合を作る。
- 空の `version` は空文字 `""` として集合に含める。
- `added` と `removed` は `name` の存在有無だけで判定する。
- `changed` と `unchanged` は、両方に存在する `name` について `version` 集合を比較して判定する。

---

## changed の意味

```bash
sbomcalc changed old.json new.json
```

以下のみを出力する。

```text
nameが両方に存在し、version集合が異なるもの
```

例:

```text
openssl
  old: 1.1.1
  new: 3.0.0
```

---

## 出力指定

### 基本文法

```text
-o FORMAT
-o FORMAT=FILE
-o FORMAT@VERSION
-o FORMAT@VERSION=FILE
```

例:

```bash
-o table
-o json=result.json
-o spdx-json@2.3=result.spdx.json
-o cyclonedx-json@1.7=result.cdx.json
```

### 複数指定

複数の `-o` を許可する。

```bash
-o table -o json=result.json -o cyclonedx-json@1.7=result.cdx.json
```

### デフォルト

`-o` が指定されない場合:

```text
table を stdout に出力する
```

stdoutとファイル:

- `FILE` を省略した出力先は stdout とする。
- stdout に出力できる指定は1つまでとする。複数の `-o` で `FILE` を省略した場合はエラーにする。
- 同じ `FILE` が複数回指定された場合はエラーにする。
- `FILE` は上書き作成する。親ディレクトリは自動作成しない。
- 計算処理が成功した後に出力ファイルを書き込む。
- いずれかの出力に失敗した場合、プロセスはエラー終了する。

---

## 出力制約

### query

| level | table | txt | json | spdx-json | cyclonedx-json |
| ----- | ----: | --: | ---: | --------: | -------------: |
| L1    |   yes | yes |  yes |        no |             no |
| L2    |   yes | yes |  yes |       yes |            yes |

`query -l1` でSBOM形式出力が指定された場合はエラーにする。

### diff

対応:

* table
* json
* txt

非対応:

* spdx-json
* cyclonedx-json

### changed

対応:

* table
* json
* txt

非対応:

* spdx-json
* cyclonedx-json

---

## table出力例

### query -l2

```text
NAME        VERSION     PURL
openssl     1.1.1       pkg:deb/debian/openssl@1.1.1
curl        7.81.0      pkg:deb/debian/curl@7.81.0
```

### query -l1

```text
NAME
openssl
curl
zlib
```

### diff

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

### changed

```text
openssl
  old: 1.1.1
  new: 3.0.0
```

---

## JSON出力例

### query -l2

```json
{
  "level": "L2",
  "components": [
    {
      "name": "openssl",
      "version": "1.1.1",
      "purl": "pkg:deb/debian/openssl@1.1.1"
    }
  ]
}
```

### query -l1

```json
{
  "level": "L1",
  "components": [
    {
      "name": "openssl"
    }
  ]
}
```

### diff

```json
{
  "from": "old.json",
  "to": "new.json",
  "added": [
    {
      "name": "nginx",
      "versions": ["1.24.0"]
    }
  ],
  "removed": [
    {
      "name": "log4j",
      "versions": ["2.14.1"]
    }
  ],
  "changed": [
    {
      "name": "openssl",
      "old_versions": ["1.1.1"],
      "new_versions": ["3.0.0"]
    }
  ],
  "unchanged": [
    {
      "name": "curl",
      "versions": ["7.81.0"]
    }
  ]
}
```

---

## 出力順

結果は決定的な順序で出力する。

- `query -l1`: `name` の昇順。
- `query -l2`: `name` の昇順、同じ `name` では `version` の昇順。
- `diff`: `added`、`removed`、`changed`、`unchanged` の順に出力する。各分類内は `name` の昇順。
- version配列は昇順に並べる。
- 空の `version` は文字列としては `""` で扱い、並び順では通常の文字列比較に従う。

---

## 内部データモデル

```go
type Level int

const (
    Level1 Level = iota
    Level2
)

type ComponentKey struct {
    Name    string
    Version string
}

type ComponentRecord struct {
    Name     string
    Version  string
    PURL     string
    Supplier string
    Licenses []string
    Hashes   []Hash
    Source   string
}

type Hash struct {
    Algorithm string
    Value     string
}
```

L1では `ComponentKey.Version` は空文字にする。

---

## インデックス

巨大SBOM対応のため、SBOM全体のRaw JSONを保持しない。

```go
type SBOMIndex struct {
    NameSet map[string]struct{}
    L2Set   map[ComponentKey]struct{}

    ByName map[string][]ComponentRecord
    ByL2   map[ComponentKey][]ComponentRecord
}
```

ただし、メモリ削減のため、queryでは2パス処理を優先する。

---

## ストリーミング方針

stdinは非対応。
入力は必ずファイルパスとする。

### query処理

1. 式をパースする
2. 式中のファイルパスを列挙する
3. 各ファイルをstreaming scanしてキー集合のみ作る
4. ASTを評価して結果キー集合を作る
5. 結果キー集合に一致するcomponentだけを再度streaming scanで収集する
6. 指定された形式で出力する

### diff / changed処理

1. old/newの2ファイルをstreaming scanする
2. name -> version集合を作る
3. added/removed/changed/unchangedを計算する
4. 出力する

---

## JSONパーサ

Go標準の `encoding/json.Decoder` を使う。

入力ファイル全体を `os.ReadFile` しない。

実装方針:

- トップレベルオブジェクトを `json.Decoder.Token` で走査する。
- SPDX JSONの `packages[]` とCycloneDX JSONの `components[]` は、配列要素を1件ずつ構造体へデコードする。
- v0.1では未知フィールドを無視する。
- JSONとして不正な入力はエラーにする。

---

## SPDX JSON読み取り

主に以下を読む。

```text
packages[]
```

ComponentRecordへのマッピング例:

```text
name                     -> Name
versionInfo              -> Version
externalRefs purl        -> PURL
supplier                 -> Supplier
licenseConcluded         -> Licenses
checksums                -> Hashes
```

SPDXRef-DOCUMENTなど、実コンポーネントでないものは除外する。

SPDX JSONの詳細:

- `packages[]` が存在しない、または配列でない場合は未対応SBOM形式としてエラーにする。
- `name` が空、または正規化後に空になったパッケージは無視する。
- `SPDXID` が `SPDXRef-DOCUMENT` のパッケージは無視する。
- `versionInfo` がない場合は空文字として扱う。
- `externalRefs` のうち `referenceType` が `purl` の最初の `referenceLocator` を `PURL` とする。
- `supplier` が `NOASSERTION` の場合は空文字として扱う。
- `licenseConcluded` が空、`NOASSERTION`、`NONE` の場合、`Licenses` は空配列にする。
- `checksums` は `algorithm` と `checksumValue` を読み取る。

---

## CycloneDX JSON読み取り

主に以下を読む。

```text
components[]
```

ComponentRecordへのマッピング例:

```text
name      -> Name
version   -> Version
purl      -> PURL
supplier  -> Supplier
licenses  -> Licenses
hashes    -> Hashes
```

dependencies は v0.1 では読まない。

CycloneDX JSONの詳細:

- `components[]` が存在しない場合は空のコンポーネント集合として扱う。
- `components[]` が配列でない場合は未対応SBOM形式としてエラーにする。
- `name` が空、または正規化後に空になったコンポーネントは無視する。
- `version` がない場合は空文字として扱う。
- `supplier` は文字列の場合はその値を使い、オブジェクトの場合は `name` を使う。
- `licenses[]` は `license.id`、`license.name`、`expression` の順に取得できる値を使う。
- `hashes[]` は `alg` と `content` を読み取る。

---

## name正規化

v0.1では最小限とする。

```text
trim space
空文字は無視
```

大文字小文字変換はしない。

---

## version正規化

v0.1では最小限とする。

```text
trim space
空文字は許容
```

空versionは `""` として扱う。

---

## 集合演算

```go
type KeySet map[ComponentKey]struct{}
```

### and

```text
A ∩ B
```

### or

```text
A ∪ B
```

### minus

```text
A - B
```

### xor

```text
(A - B) ∪ (B - A)
```

---

## エラー条件

以下はエラーにする。

* `-l1` と `-l2` の同時指定
* query式が空
* 存在しないファイル
* 未対応SBOM形式
* 未対応出力形式
* `query -l1` で `spdx-json` または `cyclonedx-json` を指定
* `diff` / `changed` でSBOM形式出力を指定
* 括弧不整合
* 不正な演算子
* 空白を含む未クォートファイル名
* stdout出力指定が複数ある
* 同じ出力ファイルが複数回指定されている
* 出力ファイルを書き込めない

終了コード:

- 正常終了は `0`。
- 入力、引数、解析、出力のエラーは `1`。

---

## 推奨ディレクトリ構成

```text
cmd/
  sbomcalc/
    main.go

internal/
  cli/
    args.go
    output_spec.go

  expr/
    lexer.go
    parser.go
    ast.go
    eval.go

  model/
    component.go
    set.go
    result.go

  reader/
    detect.go
    spdx_json.go
    cyclonedx_json.go

  writer/
    table.go
    txt.go
    json.go
    spdx_json.go
    cyclonedx_json.go

  engine/
    query.go
    diff.go
    changed.go
```

---

## 実装優先順位

### Phase 1

* CLI骨格
* query
* `and/or/minus/xor`
* L1/L2
* table/json/txt出力
* SPDX JSON読み取り
* CycloneDX JSON読み取り

### Phase 2

* diff
* changed
* spdx-json@2.2 / 2.3 出力
* cyclonedx-json@1.5 / 1.6 / 1.7 出力

### Phase 3

* 大規模ファイルでのメモリ使用量確認
* 2パス処理の最適化
* テストデータ拡充

---

## テスト観点

### query

```text
a and b
a or b
a minus b
a xor b
(a and b) minus c
```

### level

```text
-l1ではversion違いでも同一nameとして扱う
-l2ではversion違いを別要素として扱う
```

### diff

```text
added
removed
changed
unchanged
```

### output

```text
-o未指定ならtable stdout
複数-oで複数ファイル出力
FORMAT@VERSIONを解釈できる
```

---

## 非対応事項 v0.1

以下は実装しない。

* stdin入力
* 依存関係グラフ演算
* SPDX relationships保持
* CycloneDX dependencies保持
* L1結果のSBOM形式出力
* 異種レベル演算
* purl完全一致 L3
* semver比較
* 脆弱性/CVE連携

---

## 実装前の確認事項

以下は仕様として実装可能な既定値を上に定義済みだが、必要なら実装前に変更する。

1. `-o table -o txt` のように stdout 出力が複数ある場合は、混在出力を避けるためエラーにする。
2. 同じ `name` と `version` が複数ある場合は、最初に読み取ったレコードを代表として出力する。
3. `query -l1` は `name` だけを出力し、詳細情報は捨てる。
4. SBOM形式出力は最小SBOMを新規生成し、入力SBOMのメタデータは引き継がない。
5. diff / changed はSBOM形式出力より先にPhase 2で実装する。
