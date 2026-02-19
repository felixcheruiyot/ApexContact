// Package livekit provides JWT access token generation for LiveKit rooms.
// Tokens are generated manually using the existing JWT library to avoid
// pulling in the heavy LiveKit server SDK.
package livekit

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type videoGrant struct {
	RoomCreate     bool   `json:"roomCreate,omitempty"`
	RoomJoin       bool   `json:"roomJoin,omitempty"`
	RoomAdmin      bool   `json:"roomAdmin,omitempty"`
	Room           string `json:"room,omitempty"`
	CanPublish     *bool  `json:"canPublish,omitempty"`
	CanSubscribe   *bool  `json:"canSubscribe,omitempty"`
	CanPublishData *bool  `json:"canPublishData,omitempty"`
}

type liveKitClaims struct {
	Video *videoGrant `json:"video,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed LiveKit access token.
// canPublish=true lets the participant send audio/video.
// isAdmin=true grants room management permissions.
func GenerateToken(apiKey, secret, roomName, identity string, canPublish, isAdmin bool) (string, error) {
	trueVal := true
	falseVal := false

	canPub := &falseVal
	canSub := &trueVal
	canData := &trueVal
	if canPublish {
		canPub = &trueVal
	}

	grant := &videoGrant{
		RoomJoin:       true,
		Room:           roomName,
		CanPublish:     canPub,
		CanSubscribe:   canSub,
		CanPublishData: canData,
	}
	if isAdmin {
		grant.RoomAdmin = true
		grant.RoomCreate = true
	}

	now := time.Now()
	claims := liveKitClaims{
		Video: grant,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    apiKey,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(6 * time.Hour)),
			Subject:   identity,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
