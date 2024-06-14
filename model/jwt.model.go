package model

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	jwt.MapClaims
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
	Role     string `json:"role"`
}

type Whoami struct {
	SessionID string `json:"sessionID"`
}
