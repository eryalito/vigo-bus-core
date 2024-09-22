package api

// ProviderType is an enum that represents the possible values for the identity provider
type ProviderType string

const (
	// ProviderTypeTelegram represents the Telegram identity provider
	ProviderTypeTelegram ProviderType = "telegram"
)

// Identity is a struct that holds the information of a user, including their auth provider type and ID
type Identity struct {
	// ID is the unique identifier of the identity
	ID int `json:"id"`

	// UUID is the unique identifier of the identity, usually provided by the auth provider
	UUID string `json:"uuid"`

	// Provider is the type of the identity provider
	Provider ProviderType `json:"provider"`

	// FavoriteStops is a list of the user's favorite stops
	FavoriteStops []Stop `json:"favorite_stops"`
}
