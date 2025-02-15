package cache

import (
	"context"
	"fmt"
)

func (r *Storage) GetCachedRefreshToken(ctx context.Context, userID int) (string, error) {
	key := fmt.Sprintf("%s%d", "refresh_token:", userID)
	return r.db.Get(ctx, key).Result()
}
