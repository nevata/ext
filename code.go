package ext

import (
	"fmt"
	"math/rand"
	"time"
)

//GenerateCode 生成唯一ID
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%d%d%d%d%d%d",
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
	)
	return code
}
