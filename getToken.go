package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

func getJwtToken(keyFilePath string) string {
	keyFile, err := os.ReadFile(keyFilePath)
	if err != nil {
		panic(err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Add(-60 * time.Second).Unix(),
		"iss": 158946,
		"exp": time.Now().Add(10 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: getToken <pathToPrivateKey>")
		os.Exit(1)
	}
	keyPath := os.Args[1]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getJwtToken(keyPath)},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// find installation id
	// https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#authenticating-as-an-installation

	// installations, _, err := client.Apps.ListInstallations(ctx, &github.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// var installId int64
	// for _, installation := range installations {
	// 	if *installation.Account.Login == "redhat-developer" {
	// 		installId = *installation.ID
	// 		fmt.Println(installId)
	// 	}
	// }

	var installId int64 = 21318258

	token, _, err := client.Apps.CreateInstallationToken(ctx, installId, &github.InstallationTokenOptions{})
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(token, "", "  ")
	fmt.Println(string(b))
}
