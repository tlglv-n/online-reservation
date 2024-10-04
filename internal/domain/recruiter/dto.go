package recruiter

import (
	"errors"
	"net/http"
)

type Request struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Phone    int    `json:"phone"`
}

func (s *Request) Bind(r *http.Request) error {
	if s.FullName == "" {
		return errors.New("fullname: cannot be blank")
	}

	if s.Email == "" {
		return errors.New("pseudonym: cannot be blank")
	}

	if s.Phone == 0 {
		return errors.New("phone: cannot be blank")
	}

	return nil
}

type Response struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Phone    int    `json:"phone"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:       data.ID,
		FullName: *data.FullName,
		Email:    *data.Email,
		Phone:    *data.Phone,
	}
	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
