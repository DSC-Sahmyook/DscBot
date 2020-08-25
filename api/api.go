package api

import (
	"github.com/DSC-Sahmyook/dscbot/secure"
	"github.com/adlio/trello"
)

var Client = trello.NewClient(secure.AppKey, secure.Token)
