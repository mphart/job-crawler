package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupPreferences struct {
	Keywords      []string `json:"keywords"`
	Locations     []string `json:"locations"`
	DesiredTitles []string `json:"desiredTitles"`
	MinComp       int      `json:"minComp"`
	EmailOptIn    bool     `json:"emailOptIn"`
	DarkMode      bool     `json:"darkMode"`
}

type SignupRequest struct {
	Email                 string             `json:"email"`
	Name                  string             `json:"name"`
	Password              string             `json:"password"`
	NotificationFrequency string             `json:"notificationFrequency"`
	Preferences           *SignupPreferences `json:"preferences,omitempty"`
	ResumeFileName        string             `json:"resumeFileName,omitempty"`
	ResumeContentBase64   string             `json:"resumeContentBase64,omitempty"`
}
