package tests

import (
	"auth/config"
	"auth/services"
	"auth/tests/helpers"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTokenFromBearerString(t *testing.T) {

	//preparation
	cfg := &config.Config{
		AccessSecret:           "access",
		RefreshSecret:          "refresh",
		AccessLifetimeMinutes:  1,
		RefreshLifetimeMinutes: 1,
	}
	tokenService := services.NewTokenService(cfg)

	testCases := []helpers.TestCaseGetBearerToken{
		{
			Name:         "Get token successful",
			BearerString: "Bearer 56789xczasoigfacfNKOIAasdJcasm",
			Expected:     "56789xczasoigfacfNKOIAasdJcasm",
		},
		{
			Name:         "Get token from incorrect string",
			BearerString: "Bearer56789xczasoigfacfNKOIAasdJcasm",
			Expected:     "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := tokenService.GetTokenFromBearerString(testCase.BearerString)
			assert.Equal(t, testCase.Expected, actual)
		})
	}
	//
	//bearerString := "Bearer 56789xczasoigfacfNKOIAasdJcasm"
	////run
	//actual := tokenService.GetTokenFromBearerString(bearerString)
	//
	////check
	//expected := "56789xczasoigfacfNKOIAasdJcasm"
	////if want != got{
	////  t.Fail()
	////}
	//// exist assert
	//
	//assert.Equal(t, expected, actual)
}

//
//func TestGetTokenFromIncorrectBearerString(t *testing.T) {
//
//  //preparation
//  cfg := &config.Config{
//    AccessSecret:           "access",
//    RefreshSecret:          "refresh",
//    AccessLifetimeMinutes:  1,
//    RefreshLifetimeMinutes: 1,
//  }
//  tokenService := services.NewTokenService(cfg)
//
//
//
//
//  bearerString := "Bearer56789xczasoigfacfNKOIAasdJcasm"
//  //run
//  actual := tokenService.GetTokenFromBearerString(bearerString)
//
//  //check
//  expected := ""
//  //if want != got{
//  //  t.Fail()
//  //}
//  // exist assert
//
//  assert.Equal(t, expected, actual)
//}

func TestGenerateAccessToken(t *testing.T) {
	// preparation
	cfg := &config.Config{
		AccessSecret:           "access",
		RefreshSecret:          "refresh",
		AccessLifetimeMinutes:  1,
		RefreshLifetimeMinutes: 1,
	}
	tokenService := services.NewTokenService(cfg)
	userID := 1

	//run
	tokenString, err := tokenService.GenerateAccessToken(userID)

	//check
	assert.NoError(t, err)

	token, err := jwt.ParseWithClaims(tokenString, &services.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.AccessSecret), nil
	})
	assert.NoError(t, err)
	claims, ok := token.Claims.(*services.JwtCustomClaims)
	assert.True(t, ok)
	assert.True(t, token.Valid)
	actual := claims.ID
	assert.Equal(t, userID, actual)

	actualExpireTime := time.Unix(claims.ExpiresAt, 0)
	assert.WithinDuration(
		t,
		time.Now().Add(time.Minute*time.Duration(cfg.AccessLifetimeMinutes)),
		actualExpireTime,
		time.Second,
	)
}

//func TestValidateAccessToken(t *testing.T) {
//  // preparation
//  cfg := &config.Config{
//    AccessSecret:           "access",
//    RefreshSecret:          "refresh",
//    AccessLifetimeMinutes:  1,
//    RefreshLifetimeMinutes: 1,
//  }
//  tokenService := services.NewTokenService(cfg)
//  userID := 1
//  validTokenString, _ := tokenService.GenerateAccessToken(userID)
//  //run
//
//  actualClaims, err := tokenService.ValidateAccessToken(validTokenString)
//  assert.NoError(t, err)
//  assert.Equal(t, userID, actualClaims.ID)
//
//}
func TestValidateAccessToken(t *testing.T) {
	// preparation
	cfg := &config.Config{
		AccessSecret:           "access",
		RefreshSecret:          "refresh",
		AccessLifetimeMinutes:  1,
		RefreshLifetimeMinutes: 2,
	}
	tokenService := services.NewTokenService(cfg)
	userID := 1
	validTokenString, _ := tokenService.GenerateAccessToken(userID)
	validTokenRefString, _ := tokenService.GenerateRefreshToken(userID)
	expiredTokenString, _ := tokenService.GenerateAccessToken(userID)

	testCases := []helpers.TestCaseValidateAccessToken{
		{
			Name:        "Valid string",
			UserID:      userID,
			TokenString: validTokenString,
			//Expected: &services.JwtCustomClaims{
			//  ID: userID,
			//  StandardClaims: jwt.StandardClaims{
			//    ExpiresAt: time.Now().Add(time.Minute * time.Duration(cfg.AccessLifetimeMinutes)).Unix(),
			//  },
			Expected: true,
		},

		{
			Name:        "invalid token",
			UserID:      userID,
			TokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			Expected:    false,
		},
		{
			Name:        "valid token signed with refresh secret",
			UserID:      userID,
			TokenString: validTokenRefString,
			Expected:    false,
		},
		{
			Name:        "expired token",
			UserID:      userID,
			TokenString: expiredTokenString,
			Expected:    false,
		},
	}
	//{
	//  Name:        "Invalid string",
	//  UserID:      userID,
	//  TokenString: "invalidTokenString",
	//  Expected: &services.JwtCustomClaims{
	//    ID: userID,
	//    StandardClaims: jwt.StandardClaims{
	//      ExpiresAt: time.Now().Add(time.Minute * time.Duration(cfg.AccessLifetimeMinutes)).Unix(),
	//    },
	//  },
	//},
	//{
	//Name:
	//  "Valid string",
	//    UserID:      userID,
	//  TokenString: validTokenRefString,
	//  Expected: &services.JwtCustomClaims{
	//  ID: userID,
	//  StandardClaims: jwt.StandardClaims{
	//    ExpiresAt: time.Now().Add(time.Minute * time.Duration(cfg.AccessLifetimeMinutes)).Unix(),
	//  },
	//},

	//run
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actualClaims, err := tokenService.ValidateAccessToken(testCase.TokenString)
			if testCase.Expected {
				assert.NoError(t, err)
				assert.Equal(t, userID, actualClaims.ID)
			} else {
				if err != nil {
					assert.Error(t, err)
				}
				actualExpireTime := time.Now().Add(-time.Duration(cfg.AccessLifetimeMinutes) * time.Minute)
				fmt.Println(actualExpireTime)
				assert.WithinDuration(
					t,
					time.Now().Add(-time.Minute*time.Duration(cfg.AccessLifetimeMinutes)),
					actualExpireTime,
					time.Second,
				)
			}

		})
	}
	//actualClaims, err := tokenService.ValidateAccessToken(validTokenString)
	//assert.NoError(t, err)
	//assert.Equal(t, userID, actualClaims.ID)

}
