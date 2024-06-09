package pansy

import "slices"

var (
	DEFAULT_CABINET = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

func GetCabinetName(id uint8) string {
	var (
		bytes = []byte(DEFAULT_CABINET[0])
	)

	for i := range bytes {
		bytes[i] += id
	}

	return string(bytes)
}

func GetCabinetId(name string) int {
	if !slices.Contains(DEFAULT_CABINET, name) {
		return 0
	}

	return slices.Index(DEFAULT_CABINET, name)
}

func ListCabinetNames(end int) []string {
	if end > len(DEFAULT_CABINET) || end == 0 {
		return nil
	}

	return DEFAULT_CABINET[:end]
}
