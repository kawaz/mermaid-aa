# mermaid-aa 理想設計

3つの異なる言語実装（Rust, Go, TypeScript）の比較分析から導出した設計仕様。

## アーキテクチャ

```
mermaid-aa/
├── cli          # CLIエントリポイント
├── parser/      # Mermaid構文解析
│   ├── lexer    # トークン化
│   └── ast      # AST定義
├── renderer/    # ASCII Art生成
│   ├── canvas   # 描画キャンバス
│   ├── layout   # レイアウト計算
│   └── style    # 文字セット定義
├── width/       # 文字幅計算
└── types        # 共通型定義
```

### 設計原則

1. **Single Responsibility**: 各モジュールは単一責務
2. **Open/Closed**: 新しい図タイプをコア修正なしで追加可能
3. **Dependency Inversion**: パーサーとレンダラーは抽象に依存

## 機能仕様

### 対応図タイプ

| 優先度 | 図タイプ        | 説明           |
| ------ | --------------- | -------------- |
| 必須   | flowchart/graph | フローチャート |
| 推奨   | sequence        | シーケンス図   |
| 将来   | classDiagram    | クラス図       |

### ノード形状 (8種類)

| 形状          | Mermaid記法 | ASCII例    |
| ------------- | ----------- | ---------- |
| Rectangle     | `A[text]`   | `+------+` |
| Rounded       | `A(text)`   | `(------)` |
| Stadium       | `A([text])` | `([----])` |
| Diamond       | `A{text}`   | `<------>` |
| Hexagon       | `A{{text}}` | `{{----}}` |
| Circle        | `A((text))` | `((--))`   |
| Parallelogram | `A[/text/]` | `/------/` |
| Trapezoid     | `A[/text\]` | `/------\` |

### エッジスタイル

| スタイル     | 記法            | ASCII例   |
| ------------ | --------------- | --------- |
| Solid arrow  | `-->`           | `──>`     |
| Dotted arrow | `-.->`          | `- ->`    |
| Thick arrow  | `==>`           | `══>`     |
| Open         | `---`           | `───`     |
| With label   | `--\|text\|-->` | `─text─>` |

### 文字セット (5種類)

```
ascii:         +--+  |  -  +  >
unicode:       ┌──┐  │  ─  ┼  →
unicode-round: ╭──╮  │  ─  ┼  →
unicode-bold:  ┏━━┓  ┃  ━  ╋  ➤
unicode-double:╔══╗  ║  ═  ╬  ⇒
```

### 方向

- `TB` / `TD`: Top to Bottom (デフォルト)
- `BT`: Bottom to Top
- `LR`: Left to Right
- `RL`: Right to Left

## CLI仕様

```bash
mermaid-aa [OPTIONS] [INPUT]

Arguments:
  [INPUT]  Mermaid text or file path

Options:
  -f, --file <FILE>           Read from file
  -c, --charset <CHARSET>     Character set [default: unicode]
                              [ascii|unicode|unicode-round|unicode-bold|unicode-double]
  -a, --ambiguous-width <N>   Width for ambiguous chars [default: 1]
                              [1: half-width, 2: full-width]
  -d, --direction <DIR>       Override direction [TB|TD|BT|LR|RL]
  -h, --help                  Print help
  -V, --version               Print version

Input methods (priority order):
  1. -f/--file option
  2. [INPUT] argument
  3. stdin (pipe or interactive)

Examples:
  mermaid-aa 'graph TD; A-->B'
  mermaid-aa -f diagram.mmd
  echo 'graph LR; A-->B' | mermaid-aa
  mermaid-aa -c ascii -a 2 'graph TD; A[日本語]-->B'
```

### 環境変数

```bash
MERMAID_AA_CHARSET=unicode        # デフォルト文字セット
MERMAID_AA_AMBIGUOUS_WIDTH=1      # Ambiguous width設定
```

## Ambiguous Width対応

### 問題

Unicode East Asian Width の Ambiguous カテゴリの文字:

- 罫線文字 (`─`, `│`, `┌` など U+2500-257F)
- 一部の記号 (`※`, `→` など)

これらは環境によって幅1または幅2で表示される。

### 解決策

```
--ambiguous-width=1  # 欧米ターミナル向け（デフォルト）
--ambiguous-width=2  # CJKターミナル向け
```

### 実装

```
width(char) =
  if is_fullwidth(char):     2
  elif is_ambiguous(char):   ambiguous_width_setting
  elif is_halfwidth(char):   1
  else:                      1
```

## エラーハンドリング

### 構文エラー

```
Error: Parse error at line 3, column 15
  |
3 |   A --> B -->
  |               ^ Expected node identifier
```

### 警告

```
Warning: Unknown node shape 'A<text>' at line 2, treating as rectangle
```

### エラーコード

| コード | 説明                   |
| ------ | ---------------------- |
| 0      | 成功                   |
| 1      | 構文エラー             |
| 2      | ファイル読み取りエラー |
| 3      | 無効なオプション       |

## テスト戦略

### 単体テスト

- パーサー: 各構文要素のパース
- レンダラー: 各形状の描画
- Width: 文字幅計算

### スナップショットテスト

```
tests/snapshots/
├── flowchart_basic.txt
├── flowchart_lr.txt
├── sequence_basic.txt
└── unicode_labels.txt
```

### プロパティベーステスト

- ランダムな有効入力でクラッシュしない
- 出力は常に矩形（各行の幅が一致）

## パフォーマンス要件

- 正規表現はコンパイル済みを再利用
- Canvas操作は O(1) アクセス
- レイアウト計算は O(V + E)
- 循環グラフを正しく処理（無限ループ回避）

## 言語別実装ノート

### Go

- `go-runewidth` パッケージ使用
- パッケージ分割: `parser/`, `renderer/`, `width/`

### Rust

- `unicode-width` クレート使用
- `clap` でCLI
- モジュール分割: `parser.rs`, `renderer.rs`, `width.rs`

### TypeScript (Deno)

- 自前の幅計算または `@std/cli`
- strict mode必須
- ES Modules

### MoonBit

- 自前の幅計算実装
- WASM出力でブラウザ版も可能
- CLI用にNode.jsラッパー or Wasmtime
