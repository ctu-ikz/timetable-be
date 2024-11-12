package helpers

import "github.com/bwmarrin/snowflake"

func GenerateSnowflakeID() (int64, error) {
	node, err := snowflake.NewNode(69)
	if err != nil {
		return 0, err
	}

	id := node.Generate()

	return id.Int64(), nil
}
