package wes_client

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	sigv4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

// NewSigningHttpClient returns an *http.Client that will sign all requests with AWS V4 Signing.
func NewSigningHttpClient(config aws.Config) (*http.Client, error) {
	if config.Region == "" {
		return nil, fmt.Errorf("aws region is not set")
	}

	client := http.DefaultClient
	client.Transport = &sigV4Transport{
		wrapped: http.DefaultTransport,
		signer:  sigv4.NewSigner(),
		config:  config,
	}

	return client, nil
}

// sigV4Transport is a RoundTripper that will sign requests with AWS V4 Signing
type sigV4Transport struct {
	wrapped http.RoundTripper
	signer  *sigv4.Signer
	config  aws.Config
}

// RoundTrip uses the underlying RoundTripper transport, but signs request first with AWS V4 Signing
func (tr *sigV4Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	signedReq, err := tr.signRequest(*req)
	if err != nil {
		return nil, err
	}

	if tr.wrapped != nil {
		resp, err := tr.wrapped.RoundTrip(signedReq)
		return resp, err
	}

	return nil, nil
}

func (tr *sigV4Transport) signRequest(req http.Request) (*http.Request, error) {
	if !strings.Contains(req.Host, "execute-api") {
		return &req, nil
	}
	ctx := context.Background()
	payloadHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Header.Del("content-length")
		bodyHash := sha256.Sum256(body)
		payloadHash = hex.EncodeToString(bodyHash[:])
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	creds, err := tr.config.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}
	err = tr.signer.SignHTTP(ctx, creds, &req, payloadHash, "execute-api", tr.config.Region, time.Now())
	if err != nil {
		return nil, err
	}

	return &req, nil
}
