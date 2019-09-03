package util

import (
	"fmt"
	"os"
)

func Er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
