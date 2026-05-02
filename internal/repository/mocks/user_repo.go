package mocks

import (
	"context"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
)

type UserRepo struct {
	Users map[int64]*model.User
	NextID int64
}

func NewUserRepo() *UserRepo {
	return &UserRepo{Users: make(map[int64]*model.User), NextID: 1}
}

func (r *UserRepo) Create(_ context.Context, user *model.User) error {
	user.ID = r.NextID
	r.NextID++
	r.Users[user.ID] = user
	return nil
}

func (r *UserRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	for _, u := range r.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *UserRepo) GetByID(_ context.Context, id int64) (*model.User, error) {
	u, ok := r.Users[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}
