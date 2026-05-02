package model

import (
	"sort"
	"strings"
)

type Level int

const (
	Level1 Level = iota
	Level2
)

func (l Level) String() string {
	if l == Level1 {
		return "L1"
	}
	return "L2"
}

type ComponentKey struct {
	Name    string
	Version string
}

type Hash struct {
	Algorithm string `json:"algorithm,omitempty"`
	Value     string `json:"value,omitempty"`
}

type ComponentRecord struct {
	Name     string   `json:"name"`
	Version  string   `json:"version,omitempty"`
	PURL     string   `json:"purl,omitempty"`
	Supplier string   `json:"supplier,omitempty"`
	Licenses []string `json:"licenses,omitempty"`
	Hashes   []Hash   `json:"hashes,omitempty"`
	Source   string   `json:"-"`
}

func NormalizeName(name string) string {
	return strings.TrimSpace(name)
}

func NormalizeVersion(version string) string {
	return strings.TrimSpace(version)
}

func KeyFor(record ComponentRecord, level Level) ComponentKey {
	key := ComponentKey{Name: NormalizeName(record.Name)}
	if level == Level2 {
		key.Version = NormalizeVersion(record.Version)
	}
	return key
}

func SortRecords(records []ComponentRecord, level Level) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].Name != records[j].Name {
			return records[i].Name < records[j].Name
		}
		if level == Level1 {
			return false
		}
		return records[i].Version < records[j].Version
	})
}
