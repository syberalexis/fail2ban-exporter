package core

type JailsState struct {
	CurrentlyFailed map[string]uint
	TotalFailed     map[string]uint
	CurrentlyBanned map[string]uint
	TotalBanned     map[string]uint
	CountriesBanned map[string]map[string]uint
}
