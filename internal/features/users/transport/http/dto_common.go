package users_transport_http

import "github.com/MisdGhost17/golang-todoapp/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id"             example:"7"`
	Version     int     `json:"version"        example:"3"`
	FullName    string  `json:"full_name"      example:"Ivan Petrov"`
	PhoneNumber *string `json:"phone_number"   example:"+79875554433"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, userDomain := range users {
		usersDTO[i] = userDTOFromDomain(userDomain)
	}

	return usersDTO
}
