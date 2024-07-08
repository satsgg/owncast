package storageproviders

import (
	"bufio"
	"os"

	"github.com/grafov/m3u8"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

// Reads the local master playlist and writes the bandwidths to the stream output variants
func WriteVariantBandwidths(localFilePath string) error {
	f, err := os.Open(localFilePath) // nolint
	if err != nil {
		log.Fatalln(err)
	}

	p := m3u8.NewMasterPlaylist()
	if err := p.DecodeFrom(bufio.NewReader(f), false); err != nil {
		log.Warnln(err)
	}
	variants := data.GetStreamOutputVariants();

	for i, item := range p.Variants {
		variants[i].Bandwidth = int(item.Bandwidth)
	}

	data.SetStreamOutputVariants(variants)

	return nil
}