package db

import (
	"fmt"
	"log"
	"testing"
)

func TestNewClient(t *testing.T) {
	cli, err := NewClient("pgsql", "127.0.0.1", 5432, "yidongdeng", "dbuser", "123456")
	if err != nil {
		log.Printf("%s \n", err)
		return
	}

	ret, err := cli.Get("game_user", "__key__ = 'dsagag:1'", "userdata")
	if err != nil {
		log.Printf("%s \n", err)
		return
	}

	fmt.Println(ret)
}
