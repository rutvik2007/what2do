package main

import (
	"fmt"
	w2d_util "what2cook/pkg/util"
	"what2cook/pkg/youtube"
)

func main() {
	yt := youtube.New("client_secret.json")

	err := yt.Init()

	w2d_util.HandleError(err, "Unable to initialize youtube source")

	fmt.Println("Initialized!")

	// aragusea, JKenjiLopezAlt, FoodNetwork, RainbowPlantLife
	// channelId, err := getChannelID(service, "RainbowPlantLife")

	// handleError(err, "unable to get upload playlistID")

	// fmt.Println(channelId)

	// fmt.Println(getPlaylistVideos(service, uploadPlaylistID))

}
