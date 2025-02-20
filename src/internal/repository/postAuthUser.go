package repository

import (
	"MerchStore/src/internal/auth"
	"context"
	"database/sql"
	"errors"
	"time"
)

func (r Repository) PostAuthUser(ctx context.Context, username, password string) (*string, error) {
	user, err := r.Storage.GetUserByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		hashPass, err1 := auth.HashPassword(password)
		if err1 != nil {
			return nil, err1
		}
		if err1 = r.Storage.AddUser(ctx, username, hashPass, "user"); err1 != nil {
			return nil, err1
		}
		user, err = r.Storage.GetUserByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user retrieval failed after creation")
	}

	userID := user.UserID
	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, ErrMsgWrongPass
	}

	tokenS, err := r.Storage.GetRefreshToken(ctx, userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		tokenS = nil
	}

	if tokenS != nil && time.Now().Before(tokenS.ExpiresAt) {
		return &tokenS.Token, nil
	}

	expires := time.Now().Add(7 * 24 * time.Hour)
	token, err := auth.GenerateRefreshToken(userID, user.Username, expires)
	if err != nil {
		return nil, ErrMsgTokenGenFailed
	}

	if tokenS == nil {
		if err = r.Storage.SaveRefreshToken(ctx, userID, token, expires); err != nil {
			return nil, ErrMsgTokenSaveFailed
		}
	} else {
		if err = r.Storage.UpdateRefreshToken(ctx, userID, token, expires); err != nil {
			return nil, err
		}
	}

	return &token, nil
}
