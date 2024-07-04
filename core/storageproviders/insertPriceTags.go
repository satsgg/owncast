package storageproviders

import (
	"bufio"
	"os"

	"github.com/grafov/m3u8"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/playlist"

	log "github.com/sirupsen/logrus"
)

// insertPriceTags will take a local playlist and rewrite it to have #EXT-X-PRICE tags per variant
func insertPriceTags(localFilePath string) error {
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
		item.Price = uint32(variants[i].Price)
	}

	newPlaylist := p.String()

	return playlist.WritePlaylist(newPlaylist, localFilePath)
}