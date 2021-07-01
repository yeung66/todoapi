package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/yeung66/todoapi/internal/users"
	"github.com/yeung66/todoapi/pkg/jwt"

	"github.com/yeung66/todoapi/graph/generated"
	"github.com/yeung66/todoapi/graph/model"
	"github.com/yeung66/todoapi/internal/todos"
)

func (r *mutationResolver) CreateTodoItem(ctx context.Context, todo model.TodoItemInput) (*model.TodoItem, error) {
	var t todos.TodoItem
	t.Title = todo.Title
	t.Checked = todo.Checked
	t.CreatedTime = todo.CreatedTime
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
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TodoItems(ctx context.Context) ([]*model.TodoItem, error) {
	var ans []*model.TodoItem
	allTodos := todos.GetAllTodoItems()
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
	t := todos.GetTodoItem(id)
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
	allTodos := todos.GetTodoItemsByTimeRange(start, end)
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

func (r *queryResolver) Login(ctx context.Context, username string, password string) (*model.User, error) {
	user := users.GetUserByLogin(username, password)
	if user.Id == 0 {
		return &model.User{}, &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	return &model.User{
		ID:          user.Id,
		Username:    user.Username,
		CreatedTime: user.CreatedTime,
		Token:       &token,
	}, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
