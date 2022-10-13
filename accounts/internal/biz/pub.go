package biz

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/gookit/goutil/strutil"
	json "github.com/json-iterator/go"
	nanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type PubRepo interface {
	FetchRsaObj(string) (*RsaObj, error)
	StoreRsaObj(string, time.Duration, *RsaObj) error
	UniqueIdExists(string) bool
	StoreUniqueId(string, time.Duration) error
}

type PubUsecase struct {
	repo   PubRepo
	logger *zap.Logger
}

func NewRsaUsecase(repo PubRepo, logger *zap.Logger) *PubUsecase {
	SetRsaKeyPrefix("rsa")
	return &PubUsecase{repo: repo, logger: logger}
}

func (g *PubUsecase) FetchObj(requestId string, opts ...Option) (*RsaObj, error) {
	if len(requestId) == 0 {
		return nil, errors.New("empty requestId for fetch rsa key")
	}

	requestId = g.solveId(requestId)
	obj, err := g.repo.FetchRsaObj(requestId)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		rg := &rsaGenerator{rsaGlobal.bits, true, time.Duration(-1)}
		for _, opt := range opts {
			opt(rg)
		}

		if !rg.renew {
			return nil, status.Errorf(codes.NotFound, "can't find rsaObj")
		}

		rg.init()
		obj, err = rg.genRsaObj()
		if err != nil {
			return nil, err
		}

		err = g.repo.StoreRsaObj(requestId, rg.expiration, obj)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func (g *PubUsecase) FetchUniqueId(ttl time.Duration) (string, error) {
	uniqueId, err := nanoid.New()
	if err != nil {
		return "", err
	}
	err = g.repo.StoreUniqueId(uniqueId, ttl)
	return uniqueId, err
}

func (g *PubUsecase) ValidateUniqueId(uniqueId string) error {
	if strutil.IsBlank(uniqueId) || !g.repo.UniqueIdExists(uniqueId) {
		return status.Error(codes.ResourceExhausted, "Timed out, please refresh the page and try again!")
	}
	return nil
}

func (g *PubUsecase) solveId(requestId string) string {
	if len(rsaGlobal.prefix) > 0 {
		requestId = strings.TrimPrefix(requestId, ":")
	}
	requestId = rsaGlobal.prefix + requestId
	return requestId
}

// RsaObj is a rsa keys instance
type RsaObj struct {
	requestId string
	Private   []byte
	Public    []byte
}

func (o *RsaObj) MarshalBinary() (data []byte, err error) {
	return json.Marshal(o)
}

func (o *RsaObj) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

// Encrypt returns to encode data and if an error
func (o *RsaObj) Encrypt(data []byte) ([]byte, error) {
	block, _ := pem.Decode(o.Public)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
}

// Decrypt returns to decode data and if an error
func (o *RsaObj) Decrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(o.Private)
	if block == nil {
		return nil, errors.New("private key error")
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, private, ciphertext)
}

// rsaGlobalOption is using default option
type rsaGlobalOption struct {
	prefix     string
	bits       int
	expiration time.Duration
}

func SetRsaBits(bits int) {
	rsaGlobal.bits = bits
}

func SetRsaExpiration(expiration time.Duration) {
	rsaGlobal.expiration = expiration
}

// SetRsaKeyPrefix settings the redis cached rsa key prefix
func SetRsaKeyPrefix(prefix string) {
	if len(prefix) > 0 && !strings.HasSuffix(prefix, ":") {
		prefix += ":"
	}
	rsaGlobal.prefix = prefix
}

type Option func(*rsaGenerator)

func WithRsaBits(bits int) Option {
	return func(g *rsaGenerator) {
		g.bits = bits
	}
}

// WithRsaExpiration settings the global rsa expiration
// Zero expiration means the key has no expiration time.
func WithRsaExpiration(expiration time.Duration) Option {
	return func(g *rsaGenerator) {
		g.expiration = expiration
	}
}

// WithRsaNoGen when the requestId is not exist don't create
func WithRsaNoGen(g *rsaGenerator) {
	g.renew = false
}

var rsaGlobal = &rsaGlobalOption{bits: 1024, expiration: time.Minute * 10}

// rsaGenerator create a rsa key
type rsaGenerator struct {
	bits       int
	renew      bool
	expiration time.Duration
}

func (g *rsaGenerator) init() {
	if g.bits <= 0 {
		g.bits = rsaGlobal.bits
	}
	if g.expiration < 0 {
		g.expiration = rsaGlobal.expiration
	}
}

// genRsaObj generated rsa private and public key
func (g *rsaGenerator) genRsaObj() (*RsaObj, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, g.bits)
	if err != nil {
		return nil, err
	}

	// private key
	privateData := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateData,
	}
	var prv bytes.Buffer
	err = pem.Encode(&prv, privateBlock)
	if err != nil {
		return nil, err
	}

	// public key
	publicKey := &privateKey.PublicKey
	publicData := x509.MarshalPKCS1PublicKey(publicKey)
	publicBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicData,
	}
	var pub bytes.Buffer
	err = pem.Encode(&pub, publicBlock)
	if err != nil {
		return nil, err
	}

	return &RsaObj{
		Private: prv.Bytes(),
		Public:  pub.Bytes(),
	}, nil
}
