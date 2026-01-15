package auth

import (
	"context"
	"errors"
	"github.com/Dzhodddi/EcommerceAPI/pkg/contextkeys"
)

var ErrUnauthorized = errors.New("unauthorized")

func GetUserIdInt(ctx context.Context) (int, error) {
	accountID, ok := ctx.Value(contextkeys.UserIDKey).(uint64)
	if !ok {
		return 0, ErrUnauthorized
	}
	return int(accountID), nil
}
