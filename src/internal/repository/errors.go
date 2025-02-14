package repository

import "errors"

var ErrMsgOutOfStock = errors.New("the product is out of stock")
var ErrMsgNotEnoughCoins = errors.New("the user does not have enough coins")
var ErrMsgUserNotExist = errors.New("user not exist")
var ErrMsgProductNotExist = errors.New("product not exist")
