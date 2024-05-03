package utilities

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// ReadK8SecretMount : reads a single secret given the mount path for the secret
// Handy function when reading secrets and deployment
func ReadK8SecretMount(fp string) ([]string, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	byt, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if bytes.HasSuffix(byt, []byte("\n")) {
		byt, _ = bytes.CutSuffix(byt, []byte("\n")) //often file read in will have this as a suffix
	}
	/* There could be multiple secrets in the same file separated by white space */
	return strings.Split(string(byt), " "), nil
}
