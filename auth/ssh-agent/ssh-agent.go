package ssh

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	stdjwt "github.com/dgrijalva/jwt-go"
	kittjwt "github.com/go-kit/kit/auth/jwt"
	kittgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/jennyservices/jenny/auth"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"gopkg.in/square/go-jose.v2/jwt"
)

const jennySSHMetaKey = "ssh-authorization-bin"

var (
	UserExtractor auth.JWTUserExtractor = ue
	ClaimsFactory kittjwt.ClaimsFactory = func() stdjwt.Claims {
		return &Claims{}
	}

	SigningMethodSSH stdjwt.SigningMethod = &sshAgent{}
	GRPCJWTExtractor                      = kittgrpc.ServerRequestFunc(gRPCJWTExtractor)
)

type grpcSSHAgent struct {
	a       agent.ExtendedAgent
	subject string
}

func GRPCAuth(subject string) credentials.PerRPCCredentials {
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")

	if sshAuthSock == "" {
		fmt.Println("ssh agent is not running")
	}

	c, err := net.Dial("unix", sshAuthSock)
	if err != nil {
		log.Fatal(err)
	}

	return &grpcSSHAgent{agent.NewClient(c), subject}
}

type Claims jwt.Claims

func (c *Claims) Valid() error { return nil }

type SSHUser struct{ cl Claims }

func (u *SSHUser) UniqueID() []byte {
	h := sha512.New()
	io.WriteString(h, u.cl.Issuer)
	io.WriteString(h, u.cl.Subject)
	return h.Sum(nil)
}
func (u *SSHUser) Email() string                { return u.cl.Subject }
func (u *SSHUser) DisplayName() (string, error) { return "", nil }
func (u *SSHUser) Details() map[string]string   { return nil }

func (g *grpcSSHAgent) GetRequestMetadata(ctx context.Context, args ...string) (map[string]string, error) {
	creds := make(map[string]string)

	keys, err := g.a.Signers()
	if err != nil {
		return nil, err
	}
	if len(keys) < 1 {
		return nil, errors.New("there are no ssh keys added to ssh-agent")
	}
	signer := keys[0]

	claims := &stdjwt.StandardClaims{
		Issuer:  ssh.FingerprintSHA256(signer.PublicKey()),
		Subject: g.subject,
	}

	token := stdjwt.NewWithClaims(SigningMethodSSH, claims)
	ss, err := token.SignedString(signer)
	if err != nil {
		return creds, err
	}
	creds[jennySSHMetaKey] = ss
	return creds, nil
}

func (g *grpcSSHAgent) RequireTransportSecurity() bool { return false }

func gRPCJWTExtractor(ctx context.Context, md metadata.MD) context.Context {
	keys, ok := md[jennySSHMetaKey]
	if !ok {
		return ctx
	}
	if len(keys) < 1 {
		return ctx
	}
	return context.WithValue(ctx, kittjwt.JWTTokenContextKey, keys[0])
}

func ue(claims stdjwt.Claims) (auth.User, error) {
	cl, ok := claims.(*Claims)
	if !ok {
		return nil, errors.New("claims not in correct format")
	}
	return &SSHUser{*cl}, nil
}

type key struct {
	key         ssh.PublicKey
	fingerprint string
	comment     string
}

type sshKeyRing struct {
	keys map[string]key
}

func NewKeyRing(r io.Reader) (*sshKeyRing, error) {
	ks := sshKeyRing{
		keys: make(map[string]key),
	}
	payload, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	for {
		if len(payload) < 1 {
			break
		}
		pk, comment, _, rest, err := ssh.ParseAuthorizedKey(payload)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fp := ssh.FingerprintSHA256(pk)
		ks.keys[fp] = key{pk, fp, comment}
		payload = rest
	}
	return &ks, nil
}

var ErrKeyNotFound = errors.New("key not found")

func (s *sshKeyRing) GetKey(tok *stdjwt.Token) (interface{}, error) {
	claims, ok := tok.Claims.(*Claims)
	if !ok {
		return nil, ErrKeyNotFound
	}
	key, ok := s.keys[claims.Issuer]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return key, nil
}