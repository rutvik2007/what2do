package youtube

import (
	"errors"
	"fmt"
	"log"
	ifs "what2cook/pkg/interfaces"
	w2d_util "what2cook/pkg/util"
)

// implements creator interface
type Channel struct {
	Id     string
	Name   string
	Videos []ifs.Content
	source *Youtube
}

// func

func (c *Channel) FetchContent(src ifs.Source) ([]ifs.Content, error) {
	// we need to get the upload playlist for the channel
	yt := src.(*Youtube)
	channelList := yt.service.Channels.List([]string{contentDetails})

	channelList.Id(c.Id)
	channelDetails, err := channelList.Do()
	w2d_util.HandleError(err, "error getting upload playlist for channel")

	numItems := len(channelDetails.Items)
	if numItems == 0 {
		return nil, errors.New("the api does not work for this channel")
	} else if numItems != 1 {
		log.Fatalf("getChannelUploadID: searching for channel %s - Expected 1 Item, received %d", c.Name, numItems)
	}

	if channelDetails.Items[0].ContentDetails == nil || channelDetails.Items[0].ContentDetails.RelatedPlaylists == nil {
		log.Fatalf("incomplete response from api")
	}

	uploadsId := channelDetails.Items[0].ContentDetails.RelatedPlaylists.Uploads

	playlistItemsList := yt.service.PlaylistItems.List([]string{snippet, contentDetails})
	playlistItemsList.PlaylistId(uploadsId)
	playlistItemsList.MaxResults(50)

	response, err := playlistItemsList.Do()

	w2d_util.HandleError(err, fmt.Sprintf("unable to get videos in playlist %s\n", uploadsId))

	videos := make([]ifs.Content, 0)
	for _, item := range response.Items {
		videos = append(videos, &Video{
			description: item.Snippet.Description,
			title:       item.Snippet.Title,
			id:          item.Snippet.ResourceId.VideoId,
			contentType: ifs.YTVideoType,
		})
	}
	return videos, nil
}
