package users_service

import (
	"context"
	"fmt"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"get user: %w",
			err,
		)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf(
			"apply user pathc: %w",
			err,
		)
	}

	patchedUser, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"patch user: %w",
			err,
		)
	}

	return patchedUser, nil
}
