package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims representa los claims del token
type JWTClaims struct {
	AppID      string `json:"app_id"`
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	ID         string `json:"id"`
	IsActive   bool   `json:"is_active"`
	IsVerify   bool   `json:"is_verify"`
	UserID     string `json:"user_id"`
	jwt.RegisteredClaims
}

func (c *JWTClaims) String() string {
	return fmt.Sprintf("AppID: %s, EntityID: %s, EntityType: %s, ID: %s, IsActive: %t, IsVerify: %t, UserID: %s",
		c.AppID, c.EntityID, c.EntityType, c.ID, c.IsActive, c.IsVerify, c.UserID)
}
