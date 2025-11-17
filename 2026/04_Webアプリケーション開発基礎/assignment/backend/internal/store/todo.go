package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"backend/internal/models"
)

// TodoStore はTODOデータの永続化を担当する構造体
type TodoStore struct {
	sync.Mutex
	todos    []models.Todo
	nextID   int
	filePath string
}

// NewTodoStore は新しいTodoStoreインスタンスを作成
func NewTodoStore(file string) (*TodoStore, error) {
	s := &TodoStore{filePath: file, nextID: 1}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// load はファイルからTODOデータを読み込む
func (s *TodoStore) load() error {
	s.Lock()
	defer s.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o755); err != nil {
		return err
	}

	data, err := os.ReadFile(s.filePath)
	if errors.Is(err, os.ErrNotExist) || len(data) == 0 {
		s.todos = []models.Todo{}
		s.nextID = 1
		return nil
	}
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &s.todos); err != nil {
		return err
	}

	maxID := 0
	for _, t := range s.todos {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	s.nextID = maxID + 1

	return nil
}

// save はTODOデータをファイルに保存する
func (s *TodoStore) save() error {
	payload, err := json.MarshalIndent(s.todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, payload, 0o644)
}

// findIndexByID は指定されたIDのTODOのインデックスを検索する
func (s *TodoStore) findIndexByID(id int) int {
	return slices.IndexFunc(s.todos, func(t models.Todo) bool {
		return t.ID == id
	})
}

// List は全てのTODOを取得する
func (s *TodoStore) List() []models.Todo {
	s.Lock()
	defer s.Unlock()

	copied := make([]models.Todo, len(s.todos))
	copy(copied, s.todos)
	return copied
}

// Create は新しいTODOを作成する
func (s *TodoStore) Create(title string) (models.Todo, error) {
	if title == "" {
		return models.Todo{}, errors.New("title is required")
	}

	s.Lock()
	defer s.Unlock()

	todo := models.Todo{ID: s.nextID, Title: title}
	s.nextID++
	s.todos = append(s.todos, todo)

	if err := s.save(); err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}

// Update は既存のTODOを更新する
func (s *TodoStore) Update(id int, req models.UpdateTodoRequest) (models.Todo, error) {
	s.Lock()
	defer s.Unlock()

	idx := s.findIndexByID(id)
	if idx == -1 {
		return models.Todo{}, errors.New("todo not found")
	}

	if req.Title != nil {
		if *req.Title == "" {
			return models.Todo{}, errors.New("title is required")
		}
		s.todos[idx].Title = *req.Title
	}
	if req.Completed != nil {
		s.todos[idx].Completed = *req.Completed
	}

	if err := s.save(); err != nil {
		return models.Todo{}, err
	}

	return s.todos[idx], nil
}

// Delete は指定されたIDのTODOを削除する
func (s *TodoStore) Delete(id int) error {
	s.Lock()
	defer s.Unlock()

	idx := s.findIndexByID(id)
	if idx == -1 {
		return errors.New("todo not found")
	}

	s.todos = append(s.todos[:idx], s.todos[idx+1:]...)
	return s.save()
}
