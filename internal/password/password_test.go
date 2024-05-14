package password

import "testing"

var testOptions = PasswordOptions{
	Time:    5,
	Memory:  1024 * 7,
	Threads: 1,
	SaltLen: 32,
	KeyLen:  64,
}

var testPasswords = []string{
	"",
	"password",
	"dfgjdsprogmn039g<unh0qw9}]@{£}€]@{€}]{@£€$]|€}bn102387hbpgspdah07behb/&=!öoasdkfpo2u4nht08372thspdajbfnp49732w4h",
	"dfgjdsprogmn039g<unh0qw9}]@{┬ú}Ôé¼]@{Ôé¼}]{@┬úÔé¼$]|Ôé¼}bn102387hbpgspdah07behb/&=!├Âoasdkfpo2u4nht08372thspdajbfnp49732w4hdfgjdsprogmn039g<unh0qw9}]@{┬ú}Ôé¼]@{Ôé¼}]{@┬ú",
}

func TestPassword(t *testing.T) {
	for _, password := range testPasswords {
		result := Hash(password, testOptions)
		ok := Verify(result, password)
		if !ok {
			t.Error("verify failed")
		}

		_, _, options, ok := Decode(result)
		if !ok {
			t.Error("decode failed")
		}

		ok = CompareOptions(testOptions, options)
		if !ok {
			t.Error("compare options failed")
		}
	}
}
