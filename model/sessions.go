package model

import "time"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	User      string
	Token     string
	Role      string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionToken string

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
