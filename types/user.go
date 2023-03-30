package types

type User struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}
