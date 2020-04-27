package watch

import (
	"log"
	"time"
	"torrentsWatcher/internal/api/db"
	"torrentsWatcher/internal/api/models"
	"torrentsWatcher/internal/api/parsing"
)

func Watch(period time.Duration) {
	for {
		go func() {
			var torrents []models.Torrent
			err := db.DB.Find(&torrents).Error
			if err != nil {
				log.Print("Couldn't get torrents for check")
			}
			for _, torrent := range torrents {
				updatedTorrent, err := parsing.GetTorrentInfo(torrent.PageUrl)

				if err != nil {
					log.Print("Couldn't get torrents for check")
					continue
				}
				if torrent.UploadedAt != updatedTorrent.UploadedAt {
					log.Printf("torrent '%s' (%s) was updated!", torrent.Title, torrent.PageUrl)
				}

				err = torrent.UpdateFrom(updatedTorrent)
				if err != nil {
					log.Printf("Couldn't save torrent: %v", updatedTorrent)
					continue
				}
			}
		}()

		time.Sleep(period)
	}
}