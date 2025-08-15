package main

import (
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		to       string
		expected string
		offset   int64
		limit    int64
		wantErr  bool
		errValue error
	}{
		{
			name:     "empty paths",
			from:     "",
			to:       "",
			expected: "",
			offset:   0,
			limit:    0,
			wantErr:  true,
			errValue: ErrIllegalArgument,
		},
		{
			name:     "successful copy offset=0 limit=0",
			from:     "testdata/input.txt",
			to:       "testdata/output.txt",
			expected: "testdata/out_offset0_limit0.txt",
			offset:   0,
			limit:    0,
			wantErr:  false,
		},
		{
			name:     "successful copy offset=0 limit=10",
			from:     "testdata/input.txt",
			to:       "testdata/output.txt",
			expected: "testdata/out_offset0_limit10.txt",
			offset:   0,
			limit:    10,
			wantErr:  false,
		},
		{
			name:     "successful copy offset=0 limit=1000",
			from:     "testdata/input.txt",
			to:       "testdata/output.txt",
			expected: "testdata/out_offset0_limit1000.txt",
			offset:   0,
			limit:    1000,
			wantErr:  false,
		},
		{
			name:     "successful copy offset=0 limit=10000",
			from:     "testdata/input.txt",
			to:       "testdata/output.txt",
			expected: "testdata/out_offset0_limit10000.txt",
			offset:   0,
			limit:    10000,
			wantErr:  false,
		},
		{
			name:     "successful copy offset=100 limit=1000",
			from:     "testdata/input.txt",
			to:       "testdata/output.txt",
			expected: "testdata/out_offset100_limit1000.txt",
			offset:   100,
			limit:    1000,
			wantErr:  false,
		},
		{
			name:     "successful copy offset=6000 limit=1000",
			from:     "testdata/input.txt",
			to:       "testdata/output.txt",
			expected: "testdata/out_offset6000_limit1000.txt",
			offset:   6000,
			limit:    1000,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.from, tt.to, tt.offset, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != tt.errValue {
				t.Errorf("Expected error %v, got %v", tt.errValue, err)
			}

			// Проверка содержимого файлов при успешном копировании.
			if !tt.wantErr {
				result, _ := os.ReadFile(tt.to)
				expected, _ := os.ReadFile(tt.expected)
				if string(result) != string(expected) {
					t.Errorf("the files don't match:  \nresult = \n\"%s\" \nexpected = \n\"%s\"", result, expected)
				}
				os.Remove(tt.to)
			}
		})
	}
}
