package utils

var (
	AllowedMimeTypes = []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/xml",
		"text/xml",
		"image/jpeg",
		"image/png",
	}

	DocxSequence = []byte{80, 75, 3, 4}
)

func In(s string, strs []string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}

func ComplySignature(sniff []byte, signature []byte) bool {
	for idx := range signature {
		if idx >= len(sniff) {
			return false
		} else if sniff[idx] != signature[idx] {
			return false
		}
	}
	return true
}
