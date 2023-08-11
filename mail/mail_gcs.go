package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/shaka0184/GoUtil/pkg/google/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"net/http"
)

func GetClientAtGCS(ctx context.Context, credentials, bn string) (*http.Client, error) {
	b, err := storage.GetByteSlice(ctx, bn, credentials)
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		return nil, err
	}
	return getClientAtGCS(ctx, config, bn)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClientAtGCS(ctx context.Context, config *oauth2.Config, bn string) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tf := "token.json"
	tok := &oauth2.Token{}

	tokFile, err := storage.GetByteSlice(ctx, bn, tf)
	if err != nil || len(tokFile) == 0 {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		err = saveTokenAtGCS(tf, bn, tok)
		if err != nil {
			return nil, err
		}
	} else {
		err = json.NewDecoder(bytes.NewReader(tokFile)).Decode(tok)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return config.Client(ctx, tok), nil
}

// Saves a token to a file path.
func saveTokenAtGCS(fName, bn string, token *oauth2.Token) error {
	v, err := json.Marshal(token)
	if err != nil {
		return errors.WithStack(err)
	}

	err = storage.UploadFile(bn, fName, bytes.NewReader(v))

	return nil
}
