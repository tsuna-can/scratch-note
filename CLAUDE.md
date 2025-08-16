# ターミナルメモツール実装

## プロジェクト概要
Go言語でターミナル上で動作するメモツール `scratch-note` を実装する。時間をベースとしたファイル名でマークダウンファイルを生成し、指定されたエディタで開く。

## 基本仕様

### コマンド仕様
```bash
scratch-note                    # 現在時刻でファイル生成&エディタ起動
scratch-note "タイトル"         # タイトル付きでファイル生成
scratch-note --config          # 設定ファイルの編集
```

### ファイル命名規則
- 基本形式: `2025-08-16_143045.md` (YYYY-MM-DD_HHMMSS.md)
- タイトル付き: `2025-08-16_143045_タイトル.md`

### 設定ファイル
- 場所: `~/.config/scratch-note/config.yaml`
- 形式:
```yaml
scratch-note_dir: "~/scratch-notes"    # メモファイル保存先ディレクトリ
editor: "nvim"         # 使用するエディタ
```

## 動作仕様

### 基本動作
1. 現在時刻からファイル名を生成
2. 指定されたディレクトリに空のマークダウンファイルを作成
3. 設定されたエディタでファイルを開く
4. エディタ起動後、ツールは終了

### エラーハンドリング
- **保存先ディレクトリが存在しない場合**: エラーメッセージを表示して終了
- **設定ファイルが存在しない場合**: 設定ファイル作成を促すメッセージを表示
- **エディタが見つからない場合**: エラーメッセージを表示して終了
- **設定ファイル形式が不正な場合**: エラーメッセージを表示して終了

### デフォルト値
- エディタが設定されていない場合: `vi` を使用

## 実装要件

### 技術スタック
- **言語**: Go
- **設定ファイル解析**: `gopkg.in/yaml.v3`
- **パス操作**: `filepath`, `path`
- **ホームディレクトリ取得**: `os.UserHomeDir()`

### プロジェクト構成
```
scratch-note/
├── main.go
├── config/
│   └── config.go
├── utils/
│   └── file.go
├── go.mod
└── go.sum
```

### 主要機能

#### 1. 設定管理 (config/config.go)
- YAML設定ファイルの読み書き
- デフォルト設定の提供
- 設定ファイル存在確認
- 設定ファイル作成プロンプト

#### 2. ファイル操作 (utils/file.go)  
- タイムスタンプベースのファイル名生成
- ディレクトリ存在確認
- 空ファイル作成
- パス展開（~/ の処理）

#### 3. メイン処理 (main.go)
- コマンドライン引数の解析
- エディタ起動
- エラーハンドリング

### コマンドライン引数処理
- 引数なし: 基本メモファイル作成
- 1つの文字列引数: タイトル付きメモファイル作成
- `--config` フラグ: 設定ファイル編集
- その他: 使用方法を表示

### エディタ起動
- `os/exec` パッケージを使用
- エディタプロセスを起動後、即座に終了
- エディタが見つからない場合の適切なエラー処理

### クロスプラットフォーム対応
- Windows, Linux, macOS で動作
- パス区切り文字の適切な処理
- ホームディレクトリの適切な取得

## エラーメッセージ例
```
Error: scratch-note directory does not exist: /path/to/scratch-note/dir
Error: Config file not found. Run 'scratch-note --config' to create one.
Error: Editor 'nvim' not found in PATH
Error: Invalid config file format
```

## 成功メッセージ例
```
Created scratch-note: /path/to/scratch-notes/2025-08-16_143045.md
Created scratch-note: /path/to/scratch-notes/2025-08-16_143045_shopping-list.md
```

## TDD実装方針

### 開発手法
- **テスト駆動開発 (TDD)** によるRed → Green → Refactorサイクル
- Go標準の`testing`パッケージを使用
- 外部依存関係はモックを使用してテスト
- `t.TempDir()`等の標準機能でテスト環境を構築

### 実装順序
1. **utils/file.go** - ファイル名生成機能（純粋関数でテストしやすい）
2. **config/config.go** - 設定ファイル管理機能
3. **main.go** - メイン処理とエディタ起動機能

### テスト戦略
- **単体テスト**: 各パッケージの機能を個別にテスト
- **モック**: ファイルシステム操作、エディタ起動等の外部依存をモック化
- **統合テスト**: 実際のファイル操作を含むエンドツーエンドテスト
- **テストデータ**: 一時ディレクトリと設定ファイルを動的生成

### テストファイル構成
```
scratch-note/
├── main.go
├── main_test.go
├── config/
│   ├── config.go
│   └── config_test.go
├── utils/
│   ├── file.go
│   └── file_test.go
└── integration_test.go
```

## ビルド指示
- `go mod init scratch-note`
- クロスコンパイル用のMakefileも作成
- Linux, Windows, macOS向けのバイナリ生成コマンド例を含める
