package utils

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"lms_system/internal/domain/common"
	"lms_system/internal/domain/entity"
)

func GetUserFromContext(ctx context.Context) *entity.UserContext {
	userCtx, ok := ctx.Value(common.UserContextKey).(*entity.UserContext)
	if !ok {
		return nil
	}
	return userCtx
}

func ConvertKeycloakIDToUint(keycloakID string) uint {
	if keycloakID == "" {
		return 1 // Default fallback for testing
	}

	// Create MD5 hash of the Keycloak ID
	hash := md5.Sum([]byte(keycloakID))

	// Convert first 4 bytes to uint32, then to uint
	// This ensures we get a consistent mapping for the same ID
	id := binary.BigEndian.Uint32(hash[:4])

	// Ensure we don't return 0
	if id == 0 {
		id = 1
	}

	return uint(id)
}
