package youtube

// implements source interface

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Youtube struct {
	service      *youtube.Service
	authFilePath string
}

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

const ChannelItem = "youtube#channel"

// topics
const foodTopic = "/m/02wbm"

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.

var proteins = []string{"beef", "chicken", ""}

func getToken(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return tok
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func (yt *Youtube) GetChannel(forUsername, topicID string) (*youtube.Channel, error) {
	// to start off, I want to able to get video descriptions for the last 20 videos by
	// a youtuber

	// I am going to have to search for the username (not list)
	// once I get the channel ID, then I can get the playlistID

	call := yt.service.Search.List([]string{"snippet"})
	call.Q(forUsername)
	call.TopicId(topicID)
	response, err := call.Do()
	handleError(err, "error during getChannelUploadID")

	for _, item := range response.Items {
		if item.Id != nil && item.Id.Kind == ChannelItem {
			// successful
			if item.Snippet == nil {
				log.Fatalln("YT_GetCreator: Received channel without snippet")
			}
			return &types.Creator{
				Id:   item.Snippet.ChannelId,
				Name: item.Snippet.ChannelTitle,
			}, nil
		}
	}
	// unsuccessful
	return nil, errors.New("search yielded no results")
}

func GetChannelVideos(service *youtube.Service, channelId string) ([]youtube.Video, error) {
	call := service.Channels.List([]string{"contentDetails"})
	call.Id(channelId)
	response, err := call.Do()
	handleError(err, "error during getChannelUploadID")
	// fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
	// 	"and it has %d views.",
	// 	response.Items[0].Id,
	// 	response.Items[0].Snippet.Title))
	numItems := len(response.Items)
	if numItems == 0 {
		return nil, errors.New("the api does not work for this channel")
	}
	if numItems != 1 {
		log.Fatalf("getChannelUploadID: searching for channel %s - Expected 1 Item, received %d", forUsername, numItems)
	}

	call := service.PlaylistItems.List([]string{"snippet", "contentDetails"})
	call.PlaylistId(forPlaylist)
	call.MaxResults(25)
	response, err := call.Do()

	handleError(err, fmt.Sprintf("unable to get videos in playlist %s\n", forPlaylist))

	videos := make([]what2cook.Video, 0)
	for _, item := range response.Items {
		videos = append(videos, what2cook.Video{
			Description: item.Snippet.Description,
			Title:       item.Snippet.Title,
			VideoId:     item.Snippet.ResourceId.VideoId,
		})
	}
	return videos
}

func Init() {

}

func main() {

	// TODO: change the
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	token := getToken(ctx, config)

	service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	handleError(err, "Error creating YouTube client"+service.BasePath)

	// aragusea, JKenjiLopezAlt, FoodNetwork, RainbowPlantLife
	channelId, err := getChannelID(service, "RainbowPlantLife")

	handleError(err, "unable to get upload playlistID")

	fmt.Println(channelId)

	// fmt.Println(getPlaylistVideos(service, uploadPlaylistID))

}
