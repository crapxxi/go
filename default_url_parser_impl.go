package gosis1

import (
	"fmt"
	"net/url"
	"strings"
)

type DefaultURLParser struct{}

func (p *DefaultURLParser) Parse(input URLInput) (*ParsedPathResult, error) {
	raw := strings.TrimSpace(input.RawURL)

	if raw == "" {
		return nil, fmt.Errorf("%w: URL string cannot be empty", ErrInvaidInput)
	}

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: Url malformed", ErrInvaidInput)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("%w: scheme must be http or https", ErrInvaidInput)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("%w: host cannot be empty", ErrInvaidInput)
	}

	rawSegments := strings.Split(u.Path, "/")
	segments := make([]string, 0, len(rawSegments))
	for _, segment := range rawSegments {
		if segment != "" {
			segments = append(segments, segment)
		}
	}

	return &ParsedPathResult{
		Scheme:       u.Scheme,
		Host:         u.Host,
		FullPath:     u.Path,
		Segments:     segments,
		SegmentCount: len(segments),
	}, nil
}
