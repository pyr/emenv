package emenv

import (
	"bufio"
	"fmt"
	"os"
)

func Confirm() bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Proceed? [Y/n]: ")
	r, _, err := reader.ReadRune()

	if err != nil {
		return false
	}

	return (r == 'Y' || r == 'y' || r == '\n')
}
