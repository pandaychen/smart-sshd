package hash

import "crypto/md5"

func GetMd5Str(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return string(h.Sum(nil))
}
