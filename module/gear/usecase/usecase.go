package usecase

import (
	"fmt"

	"github.com/goldenfealla/gear-manager/module/gear/repository"
)

func Test() string {
	s := repository.Test()

	return fmt.Sprintf("usecase and %v", s)
}
