package pansy

import (
	"encoding/hex"
	"fmt"
)

/*Xor
 * Xor 异或计算
 */
func Xor(payload []string) string {
	var (
		bbc byte
	)
	for _, v := range payload {
		bytes, _ := hex.DecodeString(v)
		for _, b := range bytes {
			bbc ^= b
		}
	}

	return fmt.Sprintf("%02X", bbc)
}
