package l402

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Caveat struct {
	// Condition serves as a way to identify a caveat and how to satisfy it.
	Condition string

	// Value is what will be used to satisfy a caveat. This can be as
	// flexible as needed, as long as it can be encoded into a string.
	Value string
}

var (
	// ErrInvalidCaveat is an error returned when we attempt to decode a
	// caveat with an invalid format.
	ErrInvalidCaveat = errors.New("caveat must be of the form " +
		"\"condition=value\"")
)

// DecodeCaveat decodes a caveat from its string representation.
func DecodeCaveat(s string) (Caveat, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return Caveat{}, ErrInvalidCaveat
	}
	return Caveat{Condition: parts[0], Value: parts[1]}, nil
}

func VerifyCaveats(caveats []Caveat, bitrate int) error {
	// make sure expiration and max_bandwidth exist and are valid
	validUntilFound := false
	maxBitrateFound := false
	for _, caveat := range caveats {
		if validUntilFound && maxBitrateFound {
			break
		}
		if caveat.Condition == "expiration" {
			validUntil, err := strconv.ParseInt(caveat.Value, 10, 64)
			if err != nil {
				return errors.New("failed to parse expiration caveat")
			}
			if time.Now().Unix() > validUntil {
				return errors.New("l402 has expired")
			} else {
				validUntilFound = true
			}
		}
		if caveat.Condition == "max_bandwidth" {
			maxBandwidth, err := strconv.ParseInt(caveat.Value, 10, 64)
			if err != nil {
				return errors.New("failed to parse max_bandwidth caveat")
			}
			fmt.Println("bitrate", bitrate)
			fmt.Println("Max bandwidth", maxBandwidth)
			if int64(bitrate) > maxBandwidth {
				return errors.New("bandwidth paid for exceeds requested bandwidth")
			} else {
				maxBitrateFound = true
			}
		}
	}
	return nil
}