package ssh

import (
	"crypto/rand"
	"errors"

	stdjwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ssh"
)

func init() {
	stdjwt.RegisterSigningMethod("ssh", func() stdjwt.SigningMethod { return &sshAgent{} })
}

type sshAgent struct{}

func (a *sshAgent) Alg() string { return "ssh" }
func (a *sshAgent) Sign(signingString string, key interface{}) (string, error) {
	signer, ok := key.(ssh.Signer)
	if !ok {
		return "", errors.New("key not and ssh public key")
	}
	signature, err := signer.Sign(rand.Reader, []byte(signingString))
	if err != nil {
		return "", err
	}
	return stdjwt.EncodeSegment(signature.Blob), nil
}
func (a *sshAgent) Verify(signingString, signature string, k interface{}) error {
	var err error
	var sig []byte
	if sig, err = stdjwt.DecodeSegment(signature); err != nil {
		return err
	}

	key := k.(key).key

	err = key.Verify([]byte(signingString), &ssh.Signature{
		Format: key.Type(),
		Blob:   sig,
	})
	if err != nil {
		panic(err)
	}
	return nil
}