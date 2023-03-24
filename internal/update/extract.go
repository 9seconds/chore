package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/9seconds/chore/internal/network"
	"github.com/ProtonMail/go-crypto/openpgp"
)

var ErrCannotFindBinary = errors.New("cannot find binary in archive")

func Extract(ctx context.Context, archiveURL, signatureURL string) (io.Reader, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	waiters := &sync.WaitGroup{}
	errChan := make(chan error)

	var (
		archiveStorage   *bytes.Reader
		signatureStorage *bytes.Reader
	)

	waiters.Add(2) //nolint: gomnd

	go func() {
		waiters.Wait()
		close(errChan)
	}()

	go func() {
		defer waiters.Done()

		if value, err := extractDownloadTo(ctx, archiveURL); err != nil {
			select {
			case errChan <- fmt.Errorf("cannot download archive: %w", err):
			case <-ctx.Done():
			}
		} else {
			archiveStorage = value
		}
	}()

	go func() {
		defer waiters.Done()

		if value, err := extractDownloadTo(ctx, signatureURL); err != nil {
			select {
			case errChan <- fmt.Errorf("cannot download signature: %w", err):
			case <-ctx.Done():
			}
		} else {
			signatureStorage = value
		}
	}()

	if err := <-errChan; err != nil {
		return nil, err
	}

	_, err := openpgp.CheckArmoredDetachedSignature(
		getKeyring(),
		archiveStorage,
		signatureStorage,
		nil)
	if err != nil {
		return nil, fmt.Errorf("cannot validate signature: %w", err)
	}

	archiveStorage.Seek(0, 0) //nolint: errcheck

	return extractBinary(archiveStorage)
}

func extractDownloadTo(ctx context.Context, url string) (*bytes.Reader, error) {
	resp, err := network.SendRequest(network.HTTPClient, network.NewRequest(ctx, url))
	if err != nil {
		return nil, fmt.Errorf("cannot get URL: %w", err)
	}

	defer network.CloseResponse(resp) //nolint: errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	return bytes.NewReader(data), nil
}

func extractBinary(archive io.Reader) (io.Reader, error) {
	decompressedReader, err := gzip.NewReader(archive)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize gzip reader: %w", err)
	}

	tarReader := tar.NewReader(decompressedReader)

	for {
		header, err := tarReader.Next()

		switch {
		case errors.Is(err, io.EOF):
			return nil, ErrCannotFindBinary
		case err != nil:
			return nil, fmt.Errorf("cannot read file header: %w", err)
		case header.Name == "chore":
			return tarReader, nil
		}
	}
}
