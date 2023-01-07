package youtube

import (
	"errors"
	"fmt"
	"log"
	"what2cook/pkg/interfaces"
)

// implements creator interface
type Channel struct {
	Id     string
	Name   string
	Videos []interfaces.Content
}

func (c *Channel) GetContent(yt Youtube) ([]interfaces.Content, error) {
	// we need to get the upload playlist for the channel
	channelList := yt.service.Channels.List([]string{contentDetails})
	channelList.Id(c.Id)
	channelDetails, err := channelList.Do()
	handleError(err, "error getting upload playlist for channel")

	numItems := len(channelDetails.Items)
	if numItems == 0 {
		return nil, errors.New("the api does not work for this channel")
	} else if numItems != 1 {
		log.Fatalf("getChannelUploadID: searching for channel %s - Expected 1 Item, received %d", c.Name, numItems)
	}

	if channelDetails.Items[0].ContentDetails == nil || channelDetails.Items[0].ContentDetails.RelatedPlaylists == nil {
		log.Fatalf("incomplete response from api")
	}

	uploadsId := response.Items[0].ContentDetails.RelatedPlaylists.Uploads

	playlistItemsList := yt.service.PlaylistItems.List([]string{snippet, contentDetails})
	playlistItemsList.PlaylistId(uploadsId)
	playlistItemsList.MaxResults(25)
	response, err := playlistItemsList.Do()

	handleError(err, fmt.Sprintf("unable to get videos in playlist %s\n", forPlaylist))

	videos := make([]Video, 0)
	for _, item := range response.Items {
		videos = append(videos, Video{
			Description: item.Snippet.Description,
			Title:       item.Snippet.Title,
			Id:          item.Snippet.ResourceId.VideoId,
		})
	}
	return videos
}
