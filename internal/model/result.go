package model

type QueryResult struct {
	Level      Level
	Components []ComponentRecord
}

type VersionedName struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

type ChangedName struct {
	Name       string   `json:"name"`
	OldVersion []string `json:"old_versions"`
	NewVersion []string `json:"new_versions"`
}

type DiffResult struct {
	From      string
	To        string
	Added     []VersionedName
	Removed   []VersionedName
	Changed   []ChangedName
	Unchanged []VersionedName
}
