package common

import "gopkg.in/guregu/null.v4"

func PtrToNullString(s *string) null.String {
	if s == nil {
		return null.String{}
	}
	return null.NewString(*s, true)
}
