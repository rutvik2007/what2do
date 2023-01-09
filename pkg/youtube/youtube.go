package youtube

// implements source interface

import (
	"errors"
	"fmt"
	"log"
	"os"
	ifs "what2cook/pkg/interfaces"
	w2d_util "what2cook/pkg/util"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Youtube struct {
	service      *youtube.Service
	authFilePath string
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.

var Proteins = []string{"beef", "chicken", ""}

func (yt *Youtube) getChannel(forUsername, topicID string) (*Channel, error) {
	// to start off, I want to able to get video descriptions for the last 20 videos by
	// a youtuber

	// I am going to have to search for the username (not list)
	// once I get the channel ID, then I can get the playlistID

	call := yt.service.Search.List([]string{"snippet"})
	call.Q(forUsername)

	_, ok := topicIDs[topicID]
	if ok {
		call.TopicId(topicID)
	}
	response, err := call.Do()

	w2d_util.HandleError(err, "error during getChannelUploadID")

	for _, item := range response.Items {
		if item.Id != nil && item.Id.Kind == channelType {
			// successful
			if item.Snippet == nil {
				log.Fatalln("YT_GetCreator: Received channel without snippet")
			}
			return &Channel{
				id:     item.Snippet.ChannelId,
				name:   item.Snippet.ChannelTitle,
				videos: make([]ifs.Content, 0),
			}, nil
		}
	}
	// unsuccessful
	return nil, errors.New("search yielded no results")
}

func New(configFile string) ifs.Source {
	return &Youtube{
		authFilePath: configFile,
	}
}

func (yt *Youtube) Init() (err error) {
	ctx := context.Background()

	b, err := os.ReadFile(yt.authFilePath)
	if err != nil {
		return
		// log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		return
		// log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	token := getToken(ctx, config)

	service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if err != nil {
		return
	}

	yt.service = service

	return nil
}

func (yt *Youtube) GetCreator(username string, params w2d_util.SearchParameters) (ifs.Creator, error) {
	topic, _ := params.Get(TopicID)
	return yt.getChannel(username, topic)
}

func (yt *Youtube) FetchContent(creator ifs.Creator) ([]ifs.Content, error) {
	// we need to get the upload playlist for the channel
	c, ok := creator.(*Channel)
	if !ok {
		return nil, errors.New("invalid creator type")
	}
	channelList := yt.service.Channels.List([]string{contentDetails})

	channelList.Id(c.id)
	channelDetails, err := channelList.Do()
	w2d_util.HandleError(err, "error getting upload playlist for channel")

	numItems := len(channelDetails.Items)
	if numItems == 0 {
		return nil, errors.New("the api does not work for this channel")
	} else if numItems != 1 {
		log.Fatalf("getChannelUploadID: searching for channel %s - Expected 1 Item, received %d",
			c.name, numItems)
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

	// I need to decide whether videos should be a member of channel??
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
