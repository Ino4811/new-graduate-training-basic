package models

// Todo はTODOアイテムを表す構造体
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// CreateTodoRequest はTODO作成時のリクエスト構造体
type CreateTodoRequest struct {
	Title string `json:"title"`
}

// UpdateTodoRequest はTODO更新時のリクエスト構造体
type UpdateTodoRequest struct {
	Title     *string `json:"title"`     // 未指定とゼロ値の区別のためにポインタを使用
	Completed *bool   `json:"completed"` // 未指定とゼロ値の区別のためにポインタを使用
}
