package handler

import "github.com/nubesk/binn"

type responseItem struct {
	Msg string `json:"message"`
}

func bottleToResponseItem(b *binn.Bottle) *responseItem {
	return &responseItem{
		Msg: b.Msg,
	}
}

func bottlesToResponse(bs []*binn.Bottle) []*responseItem {
	is := make([]*responseItem, len(bs))
	for i, b := range bs {
		is[i] = &responseItem{
			Msg: b.Msg,
		}
	}
	return is
}
