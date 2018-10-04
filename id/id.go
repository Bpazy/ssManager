package id

import (
	"github.com/bwmarrin/snowflake"
)

var node, err = snowflake.NewNode(1)

func Next() string {
	if err != nil {
		panic(err)
	}

	return node.Generate().String()
}
