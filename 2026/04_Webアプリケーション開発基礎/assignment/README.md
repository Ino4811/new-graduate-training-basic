# 課題: TODOアプリのアーキテクチャを考えよう

## このアプリについて

このリポジトリには、基本的なTODOアプリケーションが実装されています。

**主な機能**
- TODO項目の一覧表示
- 新しいTODO項目の追加
- TODO項目のタイトル編集
- 完了/未完了の切り替え
- TODO項目の削除

**技術構成**
- Backend: Go + Echo (REST API)
- Frontend: Next.js + TypeScript
- データ保存: JSONファイル

アプリは動作しますが、全てのコードが1つのファイルに集約されており、保守性・拡張性に課題があります。

## 課題の目的

現状の実装をレビューし、アーキテクチャパターンを適用してリファクタリングする課題です。

- 課題1: Backendのアーキテクチャ設計

- 課題2: Frontendのコンポーネント設計改善（余裕がある人向け）

---

# 環境構築

```bash
docker compose up --build
```

初回起動時、フロントエンドのビルドに数分かかります。

**アクセス**:
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080

# 現状コードの課題

リファクタリングを始める前に、現在のコードの問題点を理解しましょう。

## Backend (backend/main.go)

❌ **問題点**

1. **全ロジックが1ファイルに集約** (233行)
   - HTTPハンドラー、ビジネスロジック、データアクセスが混在
   - ファイルが長く、理解が困難

2. **責務が分離されていない**
   - `todoStore` がデータアクセスとビジネスロジックを兼務
   - ハンドラー関数内にバリデーションロジックが散在 (main.go:104, 132)

3. **テストが書きづらい**
   - 依存関係がハードコーディング (main.go:167)
   - インターフェースによる抽象化がない

4. **エラーハンドリングが一貫していない**
   - `echo.ErrNotFound` と `errors.New` が混在 (main.go:128, 154)
   - 各ハンドラーで異なるエラー処理 (main.go:182-227)

5. **拡張性が低い**
   - ファイルストレージからDB切り替えが困難
   - 新機能追加時の影響範囲が広い

## Frontend (frontend/src/app/page.tsx)

❌ **問題点**

1. **巨大なコンポーネント** (352行)
   - 表示ロジック、状態管理、API通信が全て混在
   - 1つのコンポーネントが多くの責務を持つ

2. **API通信ロジックが散在**
   - fetch呼び出しが複数箇所に重複 (page.tsx:34, 60, 80, 101, 135)
   - エンドポイントURLがハードコード

3. **状態管理が複雑**
   - 6つのstateが並列に存在 (page.tsx:12-17)
   - 状態の関連性が不明瞭

4. **エラーハンドリングの重複**
   - 同じパターンが繰り返される (page.tsx:40-42, 71-73, 92-94)

5. **再利用性が低い**
   - ロジックが特定のコンポーネントに密結合
   - 他のページで同じ機能を使えない

---

# リファクタリングの参考例

## Backend: レイヤードアーキテクチャの参考例

```
backend/
├── handler/         # HTTPハンドラー (Presentation層)
│   └── todo.go
├── service/         # ビジネスロジック (Domain層)
│   └── todo.go
├── repository/      # データアクセス (Infrastructure層)
│   └── todo.go
├── model/           # データ構造体
│   └── todo.go
└── main.go          # エントリーポイント
```

---

## Frontend: Atomic Designによるコンポーネント設計の参考例

```
frontend/src/
├── app/
│   └── page.tsx                    # Pages - メインページ
├── components/
│   ├── atoms/                      # Atoms - 最小単位のUI
│   │   ├── Button.tsx              # ボタン
│   │   ├── Input.tsx               # 入力フィールド
│   │   └── Text.tsx                # テキスト表示
│   ├── molecules/                  # Molecules - Atomsの組み合わせ
│   │   ├── TodoForm.tsx            # 入力フォーム (Input + Button)
│   │   └── TodoItem.tsx            # Todo項目 (Text + Buttons)
│   ├── organisms/                  # Organisms - 独立した機能ブロック
│   │   ├── TodoList.tsx            # Todo一覧
│   │   └── TodoSection.tsx         # セクション (タイトル + TodoList)
│   └── templates/                  # Templates - ページレイアウト
│       └── TodoTemplate.tsx        # Todoアプリのレイアウト
├── hooks/
│   ├── useTodos.ts                 # CRUD操作
│   └── useEditingState.ts          # 編集状態管理
└── api/
    └── todos.ts                    # API通信層
```

---

# 参考資料

## Backend

### レイヤードアーキテクチャ
- [クリーンアーキテクチャ](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [おすすめ](https://tech.every.tv/entry/2023/12/21/115242)

### Go設計パターン
- [Echo公式ドキュメント](https://echo.labstack.com/)
- [Goにおける依存性注入](https://go.dev/blog/wire)

## Frontend

### Atomic Design
- [Atomic Design公式解説](https://atomicdesign.bradfrost.com/)
- [Atomic Designを分かりやすく解説](https://design.dena.com/design/atomic-design-%E3%82%92%E5%88%86%E3%81%8B%E3%81%A3%E3%81%9F%E3%81%A4%E3%82%82%E3%82%8A%E3%81%AB%E3%81%AA%E3%82%8B)

### React設計
- [Reactコンポーネント設計原則](https://react.dev/learn/thinking-in-react)
- [カスタムフックのベストプラクティス](https://react.dev/learn/reusing-logic-with-custom-hooks)

### Next.js
- [Next.jsプロジェクト構成](https://nextjs.org/docs/app/getting-started/project-structure)
- [App Routerガイド](https://nextjs.org/docs/app)
