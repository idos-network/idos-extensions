package ethereum

import (
	"testing"
)

func Test_isValidPublicKey(t *testing.T) {
	tests := []struct {
		publicKey string
		want      bool
	}{
		{
			"",
			false,
		},
		{
			"some random string",
			false,
		},
		{
			// Address; not a public key
			"0x9E660ba85118b722147BBaf04ED697C95549dF03",
			false,
		},
		{
			"0408cf359417716c8c4dd03ab0c3b243b383599cb05c1b276b326c92a8f4b2b4acdcbdd98e9443f8bfc370b40e80f677142dab8cffd348a22fdf4b68ab61c7d78f",
			true,
		},
		{
			"0x0408cf359417716c8c4dd03ab0c3b243b383599cb05c1b276b326c92a8f4b2b4acdcbdd98e9443f8bfc370b40e80f677142dab8cffd348a22fdf4b68ab61c7d78f",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.publicKey, func(t *testing.T) {
			if got := isValidPublicKey(tt.publicKey); got != tt.want {
				t.Errorf("isValidPublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
