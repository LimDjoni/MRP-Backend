package handler

import (
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/notificationuser"
	"ajebackend/model/rkab"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
)

type rkabHandler struct {
	rkabService             rkab.Service
	logService              logs.Service
	userIupopkService       useriupopk.Service
	historyService          history.Service
	notificationUserService notificationuser.Service
	v                       *validator.Validate
}
