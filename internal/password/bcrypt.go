package password

import "golang.org/x/crypto/bcrypt"

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

func Compare(h string, p string) error {
	bh := []byte(h)
	bp := []byte(p)

	if err := bcrypt.CompareHashAndPassword(bh, bp); err != nil {
		return err
	}
	return nil
}

func Generate(p string) (string, error) {
	bp := []byte(p)

	hp, err := bcrypt.GenerateFromPassword(bp, DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hp), err
}
