package auth_domain

type AccessToken struct {
	SS string `json:"access_token"`
}
type RefreshToken struct {
	SS  string `json:"refresh_token"`
	ID  string `json:"-"`
	UID uint   `json:"-"`
}

type TokenPair struct {
	AccessToken
	RefreshToken
}
