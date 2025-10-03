package common

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type ID = snowflake.ID

type NullID struct {
	Valid bool
	ID    ID
}

func UniqueID() snowflake.ID {
	return snowflake.New(time.Now().UTC())
}

func ParseID(id string) (ID, error) {
	return snowflake.Parse(id)
}
