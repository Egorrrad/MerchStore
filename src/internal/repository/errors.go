package repository

import "errors"

var (
	ErrMsgOutOfStock      = errors.New("the product is out of stock")
	ErrMsgNotEnoughCoins  = errors.New("the user does not have enough coins")
	ErrMsgUserNotExist    = errors.New("user not exist")
	ErrMsgProductNotExist = errors.New("product not exist")
	ErrMsgTokenGenFailed  = errors.New("token generation failed")
	ErrMsgTokenSaveFailed = errors.New("failed to save token")
	ErrMsgSentToSelf      = errors.New("user can't coin to self")
	ErrMsgInvalidAmount   = errors.New("invalid amount")
)
