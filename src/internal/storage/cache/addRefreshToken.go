package cache

import (
	"context"
	"fmt"
)

func (r *Storage) CacheRefreshToken(ctx context.Context, userID int, token string) error {
	key := fmt.Sprintf("%s%d", "refresh_token:", userID)
	return r.db.Set(ctx, key, token, 0).Err()
}
