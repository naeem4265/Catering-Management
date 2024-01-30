package data

// User credintial
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}
