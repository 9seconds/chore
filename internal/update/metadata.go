package update

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/9seconds/chore/internal/network"
)

const (
	urlRelease        = "https://api.github.com/repos/9seconds/chore/releases"
	fileNameMetadata  = "metadata.json"
	fileNameArtifacts = "artifacts.json"

	paginationPageSize = 100

	goReleaserTypeArchive   = 1
	goReleaserTypeSignature = 13
)

var ErrNoRelease = errors.New("no releases were found")

type jsonMetadata struct {
	Version     string `json:"version"`
	ProjectName string `json:"project_name"`
}

type jsonArtifact struct {
	Name string `json:"name"`
	OS   string `json:"goos"`
	Arch string `json:"goarch"`
	Type int    `json:"internal_type"`
}

type jsonRelease struct {
	CreatedAt  string `json:"created_at"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
	Assets     []struct {
		Name  string `json:"name"`
		URL   string `json:"browser_download_url"`
		State string `json:"state"`
	} `json:"assets"`
}

type Release struct {
	Version      string
	ArchiveURL   string
	SignatureURL string
}

func GetLatestRelease(ctx context.Context, withUnstable bool) (Release, error) { //nolint: cyclop
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		retval       Release
		version      string
		archiveURL   string
		signatureURL string
	)

	release, err := getLatestRelease(ctx, withUnstable)

	switch {
	case err != nil:
		return retval, fmt.Errorf("cannot collect releases from GitHub: %w", err)
	case release.CreatedAt == "":
		return retval, ErrNoRelease
	}

	assets := make(map[string]string)

	for _, asset := range release.Assets {
		if asset.State == "uploaded" {
			assets[asset.Name] = asset.URL
		}
	}

	errChan := make(chan error)
	waiters := &sync.WaitGroup{}

	waiters.Add(2) //nolint: gomnd

	go func() {
		waiters.Wait()
		close(errChan)
	}()

	go func() {
		defer waiters.Done()

		value, err := getVersion(ctx, assets)
		if err == nil {
			version = value
		} else {
			select {
			case errChan <- fmt.Errorf("cannot get the version: %w", err):
			case <-ctx.Done():
			}
		}
	}()

	go func() {
		defer waiters.Done()

		archive, signature, err := getArchiveURLs(ctx, assets)
		if err == nil {
			archiveURL = archive
			signatureURL = signature
		} else {
			select {
			case errChan <- fmt.Errorf("cannot get the url with archive: %w", err):
			case <-ctx.Done():
			}
		}
	}()

	if err := <-errChan; err != nil {
		return retval, err
	}

	retval.Version = version
	retval.ArchiveURL = archiveURL
	retval.SignatureURL = signatureURL

	return retval, nil
}

func getLatestRelease(ctx context.Context, stableOnly bool) (jsonRelease, error) {
	target := jsonRelease{}
	thisURL, _ := url.Parse(urlRelease)
	theseReleases := []jsonRelease{}

	query := thisURL.Query()
	query.Set("per_page", strconv.Itoa(paginationPageSize))

	for page := 1; ; page++ {
		theseReleases := theseReleases[:0]

		query.Set("page", strconv.Itoa(page))

		thisURL.RawQuery = query.Encode()

		if err := network.DoJSONRequest(ctx, thisURL.String(), &theseReleases); err != nil {
			return target, fmt.Errorf("cannot request %v: %w", thisURL, err)
		}

		if len(theseReleases) == 0 {
			return target, nil
		}

		for _, rel := range theseReleases {
			switch {
			case rel.Draft, stableOnly && rel.Prerelease:
			case target.CreatedAt == "", target.CreatedAt < rel.CreatedAt:
				target = rel
			}
		}
	}
}

func getVersion(ctx context.Context, assets map[string]string) (string, error) {
	value, exists := assets[fileNameMetadata]
	if !exists {
		return "", fmt.Errorf("cannot find %s within asset list", fileNameMetadata)
	}

	metadata := jsonMetadata{}

	if err := network.DoJSONRequest(ctx, value, &metadata); err != nil {
		return "", fmt.Errorf("cannot read metadata response: %w", err)
	}

	if metadata.ProjectName != "chore" {
		return "", fmt.Errorf("unexpectd project name %s", metadata.ProjectName)
	}

	return metadata.Version, nil
}

func getArchiveURLs(ctx context.Context, assets map[string]string) (string, string, error) {
	value, exists := assets[fileNameArtifacts]
	if !exists {
		return "", "", fmt.Errorf("cannot find %s within asset list", fileNameMetadata)
	}

	artifacts := []jsonArtifact{}

	if err := network.DoJSONRequest(ctx, value, &artifacts); err != nil {
		return "", "", fmt.Errorf("cannot read artifact response: %w", err)
	}

	archive, err := getArchiveFilename(artifacts)
	if err != nil {
		return "", "", fmt.Errorf("cannot find out correct archive: %w", err)
	}

	signature, err := getArchiveSignatureFilename(artifacts, archive)
	if err != nil {
		return "", "", fmt.Errorf("cannot find out correct signature: %w", err)
	}

	for _, name := range []string{archive, signature} {
		if _, ok := assets[name]; !ok {
			return "", "", fmt.Errorf("cannot find out URL for %s", archive)
		}
	}

	return assets[archive], assets[signature], nil
}

func getArchiveFilename(artifacts []jsonArtifact) (string, error) {
	name := ""

	for _, art := range artifacts {
		if art.Type == goReleaserTypeArchive && art.OS == runtime.GOOS && art.Arch == runtime.GOARCH {
			if name != "" {
				return "", fmt.Errorf(
					"found at least 2 archives for the same %s/%s: %s and %s",
					runtime.GOOS,
					runtime.GOARCH,
					name,
					art.Name)
			}

			name = art.Name
		}
	}

	if name == "" {
		return "", fmt.Errorf("cannot find out archive for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	return name, nil
}

func getArchiveSignatureFilename(artifacts []jsonArtifact, archiveFileName string) (string, error) {
	name := ""

	for _, art := range artifacts {
		if art.Type == goReleaserTypeSignature && strings.HasPrefix(art.Name, archiveFileName) {
			if name != "" {
				return "", fmt.Errorf(
					"found at least 2 signatures for the same %s: %s and %s",
					archiveFileName,
					name,
					art.Name)
			}

			name = art.Name
		}
	}

	if name == "" {
		return "", fmt.Errorf("cannot find out signature for %s", runtime.GOARCH)
	}

	return name, nil
}
