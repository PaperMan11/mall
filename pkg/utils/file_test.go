package utils

import "testing"

func TestDeleteFiles(t *testing.T) {
	err := DeletFiles("./*.txt")
	t.Log(err)
}
