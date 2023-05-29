package service

import (
	"fmt"
	"path"
	"testing"
)

func TestPathJoin(t *testing.T) {
	fmt.Println(path.Join(".", "/1234", "/24"))

}
