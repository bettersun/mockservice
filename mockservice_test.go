package mockservice

import (
	"log"
	"testing"
)

func Test_001(t *testing.T) {
	s := EscapseSlash("/bettersun/hello/")
	log.Println(s)
}

func Test002(t *testing.T) {
}
