package youtube

// implements source interface

import (
	"errors"
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
				Id:   item.Snippet.ChannelId,
				Name: item.Snippet.ChannelTitle,
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
