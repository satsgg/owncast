package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/macaroon.v2"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/l402"
	"github.com/owncast/owncast/core/lnd"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
)

// HandleHLSRequest will manage all requests to HLS content.
func HandleHLSRequest(w http.ResponseWriter, r *http.Request) {
	// Sanity check to limit requests to HLS file types.
	if filepath.Ext(r.URL.Path) != ".m3u8" && filepath.Ext(r.URL.Path) != ".ts" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	requestedPath := r.URL.Path
	relativePath := strings.Replace(requestedPath, "/hls/", "", 1)
	fullPath := filepath.Join(config.HLSStoragePath, relativePath)

	// If using external storage then only allow requests for the
	// master playlist at stream.m3u8, no variants or segments.
	if data.GetS3Config().Enabled && relativePath != "stream.m3u8" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Handle playlists
	if path.Ext(r.URL.Path) == ".m3u8" {
		// Playlists should never be cached.
		middleware.DisableCache(w)

		// Force the correct content type
		w.Header().Set("Content-Type", "application/x-mpegURL")

		// Use this as an opportunity to mark this viewer as active.
		viewer := models.GenerateViewerFromRequest(r)
		core.SetViewerActive(&viewer)

		queryParams := r.URL.Query()

		// if it's a variant playlist and the request contains a duration query param,
		// create and return a macaroon
		if fullPath != "data/hls/stream.m3u8" {
			d, dExists := queryParams["t"]
			if dExists && len(d) != 0 {
				// d, err := strconv.Atoi(d[0])
				duration, err := strconv.ParseInt(d[0], 10, 64)
				if err != nil {
					fmt.Println("Invalid d query param", d[0])
					return
				}

				// get the variant and check its price
				dir := path.Dir(fullPath)	
				indexStr := path.Base(dir)

				index, err := strconv.Atoi(indexStr)
				if err != nil {
					return 
				}

				variants := data.GetStreamOutputVariants()
				if index < 0 || index >= len(variants) {
					return 
				}
				variantPrice := variants[index].Price
				invoiceAmt := int64(variantPrice) * duration

				invoice, err := lnd.Client.AddInvoice(context.Background(), &lnrpc.Invoice{
					Memo: "l402-hls: " + variants[index].Name,
					ValueMsat: invoiceAmt,
				})
				if err != nil {
					fmt.Println("Failed to create invoice", err)
					return
				}
				println("invoice", invoice.String())

				id, err := l402.CreateUniqueIdentifier([32]byte(invoice.RHash))
				if err != nil {
					return 
				}

				mac, err := macaroon.New([]byte("rootKey123456789"), id, "l402.sats.gg", macaroon.LatestVersion)
				if err != nil {
					println("err", err)
					return
				}

				expiration := time.Now().Unix() + duration
				fmt.Printf("valid_until=%d", expiration)
				rawCaveat := []byte(fmt.Sprintf("valid_until=%d", expiration))
				if err := mac.AddFirstPartyCaveat(rawCaveat); err != nil {
					fmt.Println("Error adding expires caveat", err)
					return
				}

				rawCaveat = []byte(fmt.Sprintf("max_bitrate=%d", variants[index].VideoBitrate))
				if err := mac.AddFirstPartyCaveat(rawCaveat); err != nil {
					fmt.Println("Error adding max bitrate caveat", err)
					return
				}

				fmt.Println("Requested max bitrate", variants[index].VideoBitrate)
				println("Requested variant", variants[index].Name)

				macBytes, err := mac.MarshalBinary()
				if err != nil {
					println("Error serializing L402", err)
					return
				}

				str := fmt.Sprintf("macaroon=\"%s\", invoice=\"%s\"",
				base64.StdEncoding.EncodeToString(macBytes), invoice.PaymentRequest)

        w.Header().Set("WWW-Authenticate", "L402 " + str)
        w.WriteHeader(http.StatusPaymentRequired)
				return
			}
		}
	} else {
		// .ts segment files
		queryParams := r.URL.Query()
		l402Params, l402Exists := queryParams["l402"]
		// TODO: check variant to see if it has price > 0
		if !l402Exists {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		} 
		l402String := l402Params[0]
		// println("l402", l402String)
		parts := strings.SplitN(l402String, " ", 2)
    if len(parts) != 2 {
        fmt.Println("Unexpected input format")
        return
    }
		remaining := parts[1]

		splitParts := strings.Split(remaining, ":")
    if len(splitParts) != 2 {
        fmt.Println("Unexpected remaining string format")
        return
    }
    
		macBase64Str := splitParts[0]
		preimageHexStr := splitParts[1]

		macBytes, err := base64.StdEncoding.DecodeString(macBase64Str)
		if err != nil {
			fmt.Println("failed to decode base64 macaroon")
			return
		}
		mac := &macaroon.Macaroon{}
		err = mac.UnmarshalBinary(macBytes)
		if err != nil {
			fmt.Println("failed to unmarshal macaroon")
			return
		}

		// verify preimage hasehs to the payment hash in the macaroon's identifier
		preimage, err := lntypes.MakePreimageFromStr(preimageHexStr)
		if err != nil {
			fmt.Println("failed to decode preimage")
			return
		}	

		paymentHash, err := l402.DecodePaymentHash(bytes.NewReader(mac.Id()))
		if err != nil {
			fmt.Println("Failed to decode payment hash")
			return
		}
		if preimage.Hash() !=  *paymentHash {
			fmt.Println("invalid preimage for payment hash")
			return
		}


		rawCaveats, err := mac.VerifySignature([]byte("rootKey123456789"), nil)
		if err != nil {
			fmt.Println("failed to verify signature")
			return
		}

		// get the variant and its bitrate
		dir := path.Dir(fullPath)	
		indexStr := path.Base(dir)

		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return 
		}

		variants := data.GetStreamOutputVariants()
		if index < 0 || index >= len(variants) {
			return 
		}

		// With the L402 verified, we'll now inspect its caveats to ensure the
		// target service is authorized.
		caveats := make([]l402.Caveat, 0, len(rawCaveats))
		for _, rawCaveat := range rawCaveats {
			// L402s can contain third-party caveats that we're not aware
			// of, so just skip those.
			caveat, err := l402.DecodeCaveat(rawCaveat)
			if err != nil {
				continue
			}
			caveats = append(caveats, caveat)
			// println("caveat", rawCaveat)
		}

		err = l402.VerifyCaveats(caveats, variants[index].VideoBitrate)
		if err != nil {
			fmt.Println("Failed to verify caveats", err)
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		cacheTime := utils.GetCacheDurationSecondsForPath(relativePath)
		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(cacheTime))

	}

	middleware.EnableCors(w)
	http.ServeFile(w, r, fullPath)
}
