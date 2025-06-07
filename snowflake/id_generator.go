package snowflake

import "github.com/bwmarrin/snowflake"

var node, _ = snowflake.NewNode(1)

func Id() int64 {
	return node.Generate().Int64()
}

func IdString() string {
	return node.Generate().String()
}
