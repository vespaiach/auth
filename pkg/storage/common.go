package storage

type Direction int

const (
	BiDirection Direction = iota
	Ascendant
	Descendant
)

type Boolean struct {
	IsSet bool
	Bool  bool
}

const DefaultLimit = 10
