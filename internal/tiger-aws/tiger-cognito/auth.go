package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

const (
	// Cognito configuration
	UserPoolID = "us-east-1_qoSTRtzGT"
	ClientID   = "1df6gkm0eaiq4a33i6mu5gp1ji"
)

// CognitoAuth handles authentication with AWS Cognito
type CognitoAuth struct {
	client        *cognitoidentityprovider.Client
	region        string
	userPoolID    string
	clientID      string
	session       string
	challengeName string
}

// CognitoAuthResponse contains the response from Cognito authentication
type CognitoAuthResponse struct {
	AccessToken         string
	IdToken             string
	RefreshToken        string
	NewPasswordRequired bool
}

// NewCognitoAuth creates a new CognitoAuth instance
func NewCognitoAuth(cfg aws.Config, region, userPoolID, clientID string) *CognitoAuth {
	client := cognitoidentityprovider.NewFromConfig(cfg)
	return &CognitoAuth{
		client:     client,
		region:     region,
		userPoolID: userPoolID,
		clientID:   clientID,
	}
}

// SignIn authenticates a user with Cognito
func (c *CognitoAuth) SignIn(username, password string) (*CognitoAuthResponse, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: &c.clientID,
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}

	result, err := c.client.InitiateAuth(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	if result.ChallengeName == types.ChallengeNameTypeNewPasswordRequired {
		c.session = *result.Session
		c.challengeName = string(result.ChallengeName)

		return &CognitoAuthResponse{
			NewPasswordRequired: true,
		}, nil
	}

	return &CognitoAuthResponse{
		AccessToken:  *result.AuthenticationResult.AccessToken,
		IdToken:      *result.AuthenticationResult.IdToken,
		RefreshToken: *result.AuthenticationResult.RefreshToken,
	}, nil
}

// CompleteNewPassword completes the new password challenge
func (c *CognitoAuth) CompleteNewPassword(newPassword string) (*CognitoAuthResponse, error) {
	input := &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
		ClientId:      &c.clientID,
		Session:       &c.session,
		ChallengeResponses: map[string]string{
			"NEW_PASSWORD": newPassword,
		},
	}

	result, err := c.client.RespondToAuthChallenge(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to complete new password challenge: %w", err)
	}

	return &CognitoAuthResponse{
		AccessToken:  *result.AuthenticationResult.AccessToken,
		IdToken:      *result.AuthenticationResult.IdToken,
		RefreshToken: *result.AuthenticationResult.RefreshToken,
	}, nil
}
