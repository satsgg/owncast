package controllers

import (
	"net/http"
	"sort"

	"github.com/owncast/owncast/core/data"
)

type variantsSort struct {
	Name               string
	Index              int
	VideoBitrate       int
	IsVideoPassthrough bool
	Price              int
	Bandwidth          int
}

type variantsResponse struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
	Price int    `json:"price"`
	Bandwidth int `json:"bandwidth"`
}

// GetVideoStreamOutputVariants will return the video variants available.
func GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	outputVariants := data.GetStreamOutputVariants()

	streamSortVariants := make([]variantsSort, len(outputVariants))
	for i, variant := range outputVariants {
		variantSort := variantsSort{
			Index:              i,
			Name:               variant.GetName(),
			IsVideoPassthrough: variant.IsVideoPassthrough,
			VideoBitrate:       variant.VideoBitrate,
			Price:	            variant.Price,
			Bandwidth:           variant.Bandwidth,
		}
		streamSortVariants[i] = variantSort
	}

	sort.Slice(streamSortVariants, func(i, j int) bool {
		if streamSortVariants[i].IsVideoPassthrough && !streamSortVariants[j].IsVideoPassthrough {
			return true
		}

		if !streamSortVariants[i].IsVideoPassthrough && streamSortVariants[j].IsVideoPassthrough {
			return false
		}

		return streamSortVariants[i].VideoBitrate > streamSortVariants[j].VideoBitrate
	})

	response := make([]variantsResponse, len(streamSortVariants))
	for i, variant := range streamSortVariants {
		variantResponse := variantsResponse{
			Index: variant.Index,
			Name:  variant.Name,
			Price: variant.Price,
			Bandwidth: variant.Bandwidth,
		}
		response[i] = variantResponse
	}

	WriteResponse(w, response)
}
