package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"
	"time"

	"github.com/yeung66/todoapi/graph/generated"
	"github.com/yeung66/todoapi/graph/model"
	"github.com/yeung66/todoapi/internal/auth"
	"github.com/yeung66/todoapi/internal/todos"
	"github.com/yeung66/todoapi/internal/users"
	"github.com/yeung66/todoapi/pkg/jwt"
)

func (r *mutationResolver) CreateTodoItem(ctx context.Context, todo model.TodoItemInput) (*model.TodoItem, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.TodoItem{}, &users.PermissionDeniedError{}
	}
	var t todos.TodoItem
	t.Title = todo.Title
	t.Checked = todo.Checked
	t.CreatedTime = todo.CreatedTime
	t.User = *user
	if todo.Content != nil {
		t.Content = *todo.Content
	}
	if todo.UpdatedTime != nil {
		t.UpdatedTime = *todo.UpdatedTime
	}

	err := t.Save()
	return &model.TodoItem{
		ID:          t.Id,
		Title:       t.Title,
		Content:     &t.Content,
		Checked:     t.Checked,
		CreatedTime: t.CreatedTime,
		UpdatedTime: &t.UpdatedTime,
	}, err
}

func (r *mutationResolver) UpdateTodoItem(ctx context.Context, id int, todo model.TodoItemInput) (*model.TodoItem, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.TodoItem{}, &users.PermissionDeniedError{}
	}
	if !todos.UserHasTodoItem(user.Id, id) {
		return &model.TodoItem{}, &users.HasNoTodoItemError{}
	}

	attrs := map[string]interface{}{
		"Title":       todo.Title,
		"Checked":     todo.Checked,
		"CreatedTime": todo.CreatedTime,
	}
	if todo.Content != nil {
		attrs["Content"] = *todo.Content
	}
	if todo.UpdatedTime != nil {
		attrs["UpdatedTime"] = *todo.UpdatedTime
	}

	t := &todos.TodoItem{
		Id: id,
	}
	err := t.Update(attrs)
	return &model.TodoItem{
		ID:          t.Id,
		Title:       t.Title,
		Content:     &t.Content,
		Checked:     t.Checked,
		CreatedTime: t.CreatedTime,
		UpdatedTime: &t.UpdatedTime,
	}, err
}

func (r *mutationResolver) DeleteTodoItem(ctx context.Context, id int) (*model.TodoItem, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.TodoItem{}, &users.PermissionDeniedError{}
	}
	if !todos.UserHasTodoItem(user.Id, id) {
		return &model.TodoItem{}, &users.HasNoTodoItemError{}
	}

	t := todos.DeleteTodoItem(id)
	return &model.TodoItem{
		ID:          t.Id,
		Title:       t.Title,
		Content:     &t.Content,
		Checked:     t.Checked,
		CreatedTime: t.CreatedTime,
		UpdatedTime: &t.UpdatedTime,
	}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, username string, password string) (*model.User, error) {
	var user users.User
	var err error
	var token string
	user.Username = username
	user.Password, _ = users.HashPassword(password)
	user.CreatedTime = strconv.FormatInt(time.Now().Unix(), 10)
	token, err = jwt.GenerateToken(username)
	err = user.Save()
	if err != nil {
		return &model.User{}, err
	}

	return &model.User{
		ID:          user.Id,
		Username:    username,
		Password:    nil,
		CreatedTime: user.CreatedTime,
		Token:       &token,
	}, err
}

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.User, error) {
	user, err := users.GetUserByLogin(username, password)
	if err != nil {
		return &model.User{}, err
	}
	token, err1 := jwt.GenerateToken(user.Username)
	return &model.User{
		ID:          user.Id,
		Username:    user.Username,
		CreatedTime: user.CreatedTime,
		Token:       &token,
	}, err1
}

func (r *queryResolver) TodoItems(ctx context.Context) ([]*model.TodoItem, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return []*model.TodoItem{}, &users.PermissionDeniedError{}
	}

	var ans []*model.TodoItem
	allTodos := todos.GetUserAllTodoItems(user.Id)
	for _, t := range allTodos {
		content := t.Content
		updatedTime := t.UpdatedTime
		ans = append(ans, &model.TodoItem{
			ID:          t.Id,
			Title:       t.Title,
			Content:     &content,
			Checked:     t.Checked,
			CreatedTime: t.CreatedTime,
			UpdatedTime: &updatedTime,
		})
	}
	return ans, nil
}

func (r *queryResolver) TodoItem(ctx context.Context, id int) (*model.TodoItem, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.TodoItem{}, &users.PermissionDeniedError{}
	}
	t := todos.GetTodoItem(id)
	if t.UserId != user.Id {
		return &model.TodoItem{}, &users.HasNoTodoItemError{}
	}
	return &model.TodoItem{
		ID:          t.Id,
		Title:       t.Title,
		Content:     &t.Content,
		Checked:     t.Checked,
		CreatedTime: t.CreatedTime,
		UpdatedTime: &t.UpdatedTime,
	}, nil
}

func (r *queryResolver) TodoItemsByTimeRange(ctx context.Context, start *string, end *string) ([]*model.TodoItem, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return []*model.TodoItem{}, &users.PermissionDeniedError{}
	}

	allTodos := todos.GetUserTodoItemsByTimeRange(start, end, user.Id)
	var ans []*model.TodoItem
	for _, t := range allTodos {
		content := t.Content
		updatedTime := t.UpdatedTime
		ans = append(ans, &model.TodoItem{
			ID:          t.Id,
			Title:       t.Title,
			Content:     &content,
			Checked:     t.Checked,
			CreatedTime: t.CreatedTime,
			UpdatedTime: &updatedTime,
		})
	}
	return ans, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
