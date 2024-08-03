package controllers

import (
	"context"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/owncast/owncast/core/lnd"
	"github.com/owncast/owncast/router/middleware"
)

type bolt11Response struct {
	Status  string `json:"status"`
	Preimage string    `json:"preimage"`
}

func CheckPaymentStatus(w http.ResponseWriter, r *http.Request) {
	r.URL.Query()
	queryParams := r.URL.Query()
	hash, hashExists := queryParams["h"]
	middleware.EnableCors(w)

	if !hashExists {
		BadRequestHandler(w, errors.New("missing payment hash"))
		return
	}
	invoice, err := lnd.Client.LookupInvoice(context.Background(), &lnrpc.PaymentHash{
		RHashStr: hash[0],
	})
	// invoice, err := lnd.Client.AddInvoice(context.Background(), &lnrpc.Invoice{
	// 	Memo: "test",
	// 	Value: 100000,
	// })
	if err != nil {
		println("Failed to lookup invoice")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if invoice.State == lnrpc.Invoice_SETTLED {
		WriteResponse(w, bolt11Response{
			Status: invoice.State.String(),
			// Preimage: base64.StdEncoding.EncodeToString(invoice.RPreimage),
			Preimage: hex.EncodeToString(invoice.RPreimage),

		})
	} else {
		WriteResponse(w, bolt11Response{
			Status: invoice.State.String(),
			Preimage: "",
		})
	}
}