package sso

import (
	"context"
	"log"
	"www-api/internal/constants"

	"github.com/coreos/go-oidc/v3/oidc"
)

var claims struct {
	Role string
}

// VerifyToken validates and verify the raw token based on deployment
func VerifyToken(rawToken string, environement string) (bool, error) {
	var err error
	// by default creates a provider of prod type
	provider, err := oidc.NewProvider(context.Background(), constants.ProdAuthUrl)
	if err != nil {
		log.Fatalf("Could not setup oidc connect verification with prod: %v\n", err)
		return false, err
	}

	//from provider create a verifier
	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})

	//check if deployment type is of dev or rtqa
	if environement == constants.DevEnvironment {
		//overwrite provider in case of dev or rtqa deployment
		provider, err = oidc.NewProvider(context.Background(), constants.DevAuthUrl)
		if err != nil {
			log.Fatalf("Could not setup oidc connect verification with rtqa: %v\n", err)
			return false, err
		}
		//overwrite verifyer too
		verifier = provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	}

	//call openid_connect_verification to verify token
	if openid_connect_verification(verifier, rawToken) {
		return true, nil
	}

	return false, nil
}

// openid_connect_verification takes a verifyer and token to verify
func openid_connect_verification(verifier *oidc.IDTokenVerifier, token string) bool {
	//parse and verify Token payload.
	idToken, err := verifier.Verify(context.Background(), token)
	if err != nil {
		log.Printf("Authentication Error: Token failed verification: %v '%s'", err, token)
		return false
	}

	//fetch claims and check if role is backend
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("Authentication Error: Could not parse token claims")
		return false
	}

	if claims.Role == "Backend" {
		return true
	}

	log.Printf("Authentication Error: Invalid token, requires Role = 'Backend'")
	return false
}
