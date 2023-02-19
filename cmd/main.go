package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	w2d_util "what2cook/pkg/util"
	"what2cook/pkg/youtube"
)

func main() {

	// as a part of initializing the application, we should start a webserver
	// the webserver will read the code from the oauth2 request
	yt := youtube.New("config/client_secret.json")

	err := yt.Init()

	w2d_util.HandleError(err, "Unable to initialize youtube source")

	fmt.Println("Initialized!")

	// aragusea, JKenjiLopezAlt, FoodNetwork, RainbowPlantLife
	// channelId, err := getChannelID(service, "RainbowPlantLife")

	// handleError(err, "unable to get upload playlistID")

	// fmt.Println(channelId)

	// fmt.Println(getPlaylistVideos(service, uploadPlaylistID))

	params := w2d_util.SearchParameters{
		youtube.TopicID: youtube.FoodTopic,
	}

	username := "Kenji Lopez"

	kla, err := yt.GetCreator(username, params)
	w2d_util.HandleError(err, "Unable to get youtube creator")

	// yt.fetchcontent(kla)???
	videos, err := yt.FetchContent(kla, 500)

	w2d_util.HandleError(err, "unable to fetch content for "+username)

	// for _, video := range videos {
	// 	fmt.Println(video.CreatedAt())
	// }
	jsonFile, err := json.Marshal(videos)
	if err != nil {
		log.Fatalln("Unable to marshal json file")
	}
	os.WriteFile("kenji.json", jsonFile, 0777)
}
