package viewer

import (
	"fmt"

	"github.com/ferux/gitPractice/rnd"
)

//View variables
func View() {
	fmt.Println("Version:", rnd.Version)
	fmt.Println("InitVersion:", rnd.InitVersion)
	rnd.Init("Alex", 19)
	fmt.Println("Version:", rnd.Version)
	fmt.Println("InitVersion:", rnd.InitVersion)
}
