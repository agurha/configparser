package configparser_test

import (
	"../configparser"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {

	pwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Cannot get the current config directory " + err.Error())
	}

	tmp := pwd + "/default.config"
	defer os.Remove(tmp)

	c := configparser.ReturnString("World")

	if c != "WorldHello" {
		t.Errorf("Expected different valye for %s", c)
	}

}
