package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type CognitoAuth struct {
	client     *cognitoidentityprovider.Client
	region     string
	userPoolID string
	clientID   string
	// State for handling challenges
	session  string
	username string
}

func NewCognitoAuth(cfg aws.Config, region string, userPoolID string, clientID string) *CognitoAuth {
	return &CognitoAuth{
		client:     cognitoidentityprovider.NewFromConfig(cfg),
		region:     region,
		userPoolID: userPoolID,
		clientID:   clientID,
	}
}

func (c *CognitoAuth) SignIn(username, password string) (*CognitoAuthResponse, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: &c.clientID,
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}

	result, err := c.client.InitiateAuth(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result == nil || result.AuthenticationResult == nil {
		if result.ChallengeName == "NEW_PASSWORD_REQUIRED" {
			// Store session and username for later use
			if result.Session != nil {
				c.session = *result.Session
			} else {
				return nil, fmt.Errorf("no session returned for password change")
			}
			c.username = username
			return &CognitoAuthResponse{
				NewPasswordRequired: true,
			}, nil
		}
		return nil, fmt.Errorf("authentication failed")
	}

	var accessToken, idToken, tokenType, refreshToken string
	if result.AuthenticationResult.AccessToken != nil {
		accessToken = *result.AuthenticationResult.AccessToken
	}
	if result.AuthenticationResult.IdToken != nil {
		idToken = *result.AuthenticationResult.IdToken
	}
	if result.AuthenticationResult.TokenType != nil {
		tokenType = *result.AuthenticationResult.TokenType
	}
	if result.AuthenticationResult.RefreshToken != nil {
		refreshToken = *result.AuthenticationResult.RefreshToken
	}

	return &CognitoAuthResponse{
		AccessToken:  accessToken,
		IdToken:      idToken,
		TokenType:    tokenType,
		ExpiresIn:    result.AuthenticationResult.ExpiresIn,
		RefreshToken: refreshToken,
	}, nil
}

func (c *CognitoAuth) CompleteNewPassword(newPassword string) (*CognitoAuthResponse, error) {
	if c.session == "" || c.username == "" {
		return nil, fmt.Errorf("no active session, please sign in first")
	}

	input := &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: "NEW_PASSWORD_REQUIRED",
		ClientId:      &c.clientID,
		ChallengeResponses: map[string]string{
			"USERNAME":     c.username,
			"NEW_PASSWORD": newPassword,
		},
		Session: &c.session,
	}

	result, err := c.client.RespondToAuthChallenge(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result == nil || result.AuthenticationResult == nil {
		return nil, fmt.Errorf("failed to complete password change")
	}

	var accessToken, idToken, tokenType, refreshToken string
	if result.AuthenticationResult.AccessToken != nil {
		accessToken = *result.AuthenticationResult.AccessToken
	}
	if result.AuthenticationResult.IdToken != nil {
		idToken = *result.AuthenticationResult.IdToken
	}
	if result.AuthenticationResult.TokenType != nil {
		tokenType = *result.AuthenticationResult.TokenType
	}
	if result.AuthenticationResult.RefreshToken != nil {
		refreshToken = *result.AuthenticationResult.RefreshToken
	}

	// Clear the session and username
	c.session = ""
	c.username = ""

	return &CognitoAuthResponse{
		AccessToken:  accessToken,
		IdToken:      idToken,
		TokenType:    tokenType,
		ExpiresIn:    result.AuthenticationResult.ExpiresIn,
		RefreshToken: refreshToken,
	}, nil
}

type CognitoAuthResponse struct {
	AccessToken         string
	IdToken             string
	TokenType           string
	ExpiresIn           int32
	RefreshToken        string
	NewPasswordRequired bool
}
