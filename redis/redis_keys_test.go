package redis

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(ThisWeekScoreKey("1"))
}
