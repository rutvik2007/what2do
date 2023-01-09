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
	"what2cook/pkg/interfaces"
	w2d_util "what2cook/pkg/util"

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

func (yt *Youtube) GetChannel(forUsername, topicID string) (*Channel, error) {
	// to start off, I want to able to get video descriptions for the last 20 videos by
	// a youtuber

	// I am going to have to search for the username (not list)
	// once I get the channel ID, then I can get the playlistID

	call := yt.service.Search.List([]string{"snippet"})
	call.Q(forUsername)
	call.TopicId(topicID)
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

func New(configFile string) interfaces.Source {
	return &Youtube{
		authFilePath: configFile,
	}
}

func (yt *Youtube) Init() (err error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(yt.authFilePath)
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
