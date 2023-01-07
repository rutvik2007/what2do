config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
if err != nil {
	log.Fatalf("Unable to parse client secret file to config: %v", err)
}
token, err := config.Exchange(ctx, string(b))
handleError(err, "unable to do config.Exchange don't really know what that even means")
// Low priority TODO: replace deprecated function
service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))