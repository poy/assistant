package userinput

import (
	"bufio"
	"errors"
	"os"
)

// ReadLine returns a line from the user.
func ReadLine() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		return scanner.Text(), nil
	}

	if err := scanner.Err(); err != nil {
		return "", scanner.Err()
	}
	return "", errors.New("unable to get input from user")
}
