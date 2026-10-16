package gosis1

type URLParser interface {
	Parse(input URLInput) (*ParsedPathResult, error)
}
