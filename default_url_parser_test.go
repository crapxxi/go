package gosis1

import "testing"

func TestURLParser(t *testing.T) {
	parser := &DefaultURLParser{}

	tests := []struct {
		name         string
		input        URLInput
		wantErr      bool
		wantScheme   string
		wantHost     string
		wantCount    int
		wantSegments []string
	}{
		{
			name:         "Valid URL with multiple segments",
			input:        URLInput{RawURL: "https://smth.com/URL/TO/PARSE"},
			wantErr:      false,
			wantScheme:   "https",
			wantHost:     "smth.com",
			wantCount:    3,
			wantSegments: []string{"URL", "TO", "PARSE"},
		},
		{
			name:    "Invalid URL",
			input:   URLInput{RawURL: "asdf"},
			wantErr: true,
		},
		{
			name:    "Empty URL",
			input:   URLInput{RawURL: ""},
			wantErr: true,
		},
		{
			name:    "Invalid Scheme",
			input:   URLInput{RawURL: "htt://smth.com"},
			wantErr: true,
		},
		{
			name:       "Empty host",
			input:      URLInput{RawURL: "https:///URL/TO/PARSE"},
			wantScheme: "https",
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := parser.Parse(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("Parse() error = %v, want = %v", err, test.wantErr)
			}

			if !test.wantErr {
				if res.Scheme != test.wantScheme {
					t.Errorf("Scheme = %v, want = %v", res.Scheme, test.wantScheme)
				}

				if res.Host != test.wantHost {
					t.Errorf("Host = %v, want = %v", res.Host, test.wantHost)
				}

				if res.SegmentCount != test.wantCount {
					t.Errorf("SegmentCount = %v, want = %v", res.SegmentCount, test.wantCount)
				}

				for i, seg := range res.Segments {
					if seg != test.wantSegments[i] {
						t.Errorf("Segment[%d] = %v, want = %v", i, seg, test.wantSegments[i])
					}
				}
			}
		})
	}
}
