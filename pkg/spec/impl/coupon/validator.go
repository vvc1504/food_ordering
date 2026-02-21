package coupon

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"sync"

	"github.com/vvc1504/food_ordering/pkg/spec"
)

type validator struct {
	validCoupons map[string]struct{}
	mu           sync.RWMutex
}

// NewValidator creates and initializes a CouponValidator by processing the provided gzipped files.
func NewValidator(filePaths []string) (spec.CouponValidator, error) {
	v := &validator{
		validCoupons: make(map[string]struct{}),
	}

	if err := v.initialize(filePaths); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *validator) Validate(Code string) (Resp bool, Err spec.ErrorCouponValidatorStatus) {
	if len(Code) < 8 || len(Code) > 10 {
		return false, nil
	}

	v.mu.RLock()
	defer v.mu.RUnlock()

	_, ok := v.validCoupons[Code]
	return ok, nil
}

func (v *validator) initialize(filePaths []string) error {
	if len(filePaths) < 2 {
		return fmt.Errorf("at least two coupon base files are required")
	}

	// seenInFile1 tracks strings found in the first file.
	seenInFile1 := make(map[string]struct{})
	// seenInFile2 tracks strings found in the second file (not present in file 1).
	seenInFile2 := make(map[string]struct{})

	for i, path := range filePaths {
		err := func() error {
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer f.Close()

			gz, err := gzip.NewReader(f)
			if err != nil {
				return fmt.Errorf("failed to create gzip reader for %s: %w", path, err)
			}
			defer gz.Close()

			scanner := bufio.NewScanner(gz)
			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				word := scanner.Text()
				if len(word) >= 8 && len(word) <= 10 {
					switch i {
					case 0:
						seenInFile1[word] = struct{}{}
					case 1:
						if _, ok := seenInFile1[word]; ok {
							v.validCoupons[word] = struct{}{}
						} else {
							seenInFile2[word] = struct{}{}
						}
					default:
						// For 3rd and subsequent files
						if _, ok := seenInFile1[word]; ok {
							v.validCoupons[word] = struct{}{}
						} else if _, ok := seenInFile2[word]; ok {
							v.validCoupons[word] = struct{}{}
						}
					}
				}
			}
			return scanner.Err()
		}()

		if err != nil {
			return err
		}
	}

	// Cleanup intermediate maps to save memory
	// Note: We keep validCoupons which is what we need for validation.
	// seenInFile1 and seenInFile2 will be GC'd after this function returns.

	return nil
}
