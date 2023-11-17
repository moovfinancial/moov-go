package moov

import (
	"log"
	"testing"
)

func TestReadConfig(t *testing.T) {
	config, err := readConfig()
	log.Println(config)

	if err != nil {
		t.Fatal(err)
	}
}
