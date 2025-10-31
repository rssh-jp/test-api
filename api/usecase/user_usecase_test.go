package usecase

import (
	"context"
	"testing"

	"github.com/rssh-jp/test-api/api/domain"
)

// Mock repository for testing
type mockUserRepository struct {
	users []domain.User
}

func (m *mockUserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	return m.users, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = int64(len(m.users) + 1)
	m.users = append(m.users, *user)
	return nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	for i, u := range m.users {
		if u.ID == user.ID {
			m.users[i] = *user
			return nil
		}
	}
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id int64) error {
	for i, u := range m.users {
		if u.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestGetAllUsers(t *testing.T) {
	mockRepo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Name: "Test User", Email: "test@example.com"},
		},
	}
	usecase := NewUserUsecase(mockRepo)

	ctx := context.Background()
	users, err := usecase.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(users))
	}

	if users[0].Name != "Test User" {
		t.Errorf("Expected user name 'Test User', got '%s'", users[0].Name)
	}
}

func TestCreateUser(t *testing.T) {
	mockRepo := &mockUserRepository{
		users: []domain.User{},
	}
	usecase := NewUserUsecase(mockRepo)

	ctx := context.Background()
	age := int32(25)
	user, err := usecase.CreateUser(ctx, "New User", "new@example.com", &age)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Name != "New User" {
		t.Errorf("Expected user name 'New User', got '%s'", user.Name)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}
