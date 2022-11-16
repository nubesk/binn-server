package handler

import "github.com/nubesk/binn"

type request struct {
	Msg string `json:"message"`
}

func requestToBottle(body *request) *binn.Bottle {
	return &binn.Bottle{
		Msg: body.Msg,
	}
}
