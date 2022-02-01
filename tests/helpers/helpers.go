package helpers

type TestCaseGetBearerToken struct {
	Name         string
	BearerString string
	Expected     string
}

type TestCaseValidateAccessToken struct {
	Name        string
	UserID      int
	TokenString string
	//Expected    *services.JwtCustomClaims
	Expected bool
	//ExpiredTime time.Time
}
