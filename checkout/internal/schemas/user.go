package schemas

import "errors"

var (
	ErrMissedUser = errors.New("user required")
)

type UserPayload struct {
	User int64 `json:"user"`
}

func (p UserPayload) Validate() error {
	if p.User == 0 {
		return ErrMissedUser
	}
	return nil
}
