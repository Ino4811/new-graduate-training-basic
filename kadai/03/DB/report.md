# 課題 3

## 1 SQL Table の正規化

### 1－1 問題点の指摘

現在のテーブルには user_id と user_name、recipe_id と recipe_title で重複したカラムが存在しており冗長性がある。冗長性のあるテーブルはデータ変更時に一貫したフィールドの変更を必要とするという保守上の問題が存在する。

### 1-2 テーブルの正規化

users Table

| カラム名    | 型      |
| ----------- | ------- |
| **user_id** | VARCHAR |
| user_name   | VARCHAR |

---

recipes Table

| カラム名      | 型      |
| ------------- | ------- |
| **recipe_id** | VARCHAR |
| recipe_title  | VARCHAR |

---

favorites Table

| カラム名          | 型      |
| ----------------- | ------- |
| **favorite_id**   | SERIAL  |
| user_id           | VARCHAR |
| recipe_id         | VARCHAR |
| registration_date | DATE    |

### 1-3 SQL クエリの作成

```
SELECT r.recipe_title
FROM favorites f
JOIN users u ON f.user_id = u.user_id
JOIN recipes r ON f.recipe_id = r.recipe_id
WHERE u.user_name = '田中 圭';
```

## 2 巨大ログテーブルのパフォーマンス改善

### 2-1 パフォーマンスが遅い原因の推測

WHERE 句で recipe_id と access_timestamp に対する絞り込み操作部分で、効率的な探索が出来ていないと考えられる。
また、出力時のソートにおいても同様と考えられる。

### 2-2 インデックスの設計

2－1 で示した範囲検索や完全一致を行う句が存在する事を踏まえて、recipe_id と access_timestamp に対してインデックスを追加するべきであると考える。
