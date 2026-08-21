package encoder

import "errors"

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var ErrInvalidBase62 = errors.New("invalid base62 string")

func Encode(num int64) string {
	if num == 0 {
		return string(alphabet[0])
	}

	var result []byte
	for num > 0 {
		remainder := num % 62
		result = append(result, []byte(alphabet)[remainder])
		num /= 62
	}

	// Reverse result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

func Decode(value string) (int64, error) {
	var num int64

	for _, char := range value {
		var index int64

		switch {
		case char >= '0' && char <= '9':
			index = int64(char - '0')

		case char >= 'A' && char <= 'Z':
			index = int64(char-'A') + 10

		case char >= 'a' && char <= 'z':
			index = int64(char-'a') + 36

		default:
			return 0, ErrInvalidBase62
		}

		num = num*62 + index
	}

	return num, nil
}
