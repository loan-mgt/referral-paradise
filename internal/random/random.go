package random

import (
    "math/rand"
)

func GetRandomNumberFromSeed(seed int64) int {
    rand.Seed(seed)
    return rand.Int()
}
