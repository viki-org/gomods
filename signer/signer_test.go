package signer

import (
	time "github.com/ssoroka/ttime"
	"github.com/viki-org/gspec"
	"strconv"
	"testing"
)

var (
	testAppId  = "123a"
	testSecret = "secret"
)

func TestCanaryGspec(t *testing.T) {
	spec := gspec.New(t)
	spec.Expect(true).ToEqual(true)
	spec.Expect(false).ToNotEqual(true)
}

func TestCreateNewSigner(t *testing.T) {
	spec := gspec.New(t)
	signer := New(testAppId, testSecret)
	spec.Expect(signer.AppID).ToEqual(testAppId)
	spec.Expect(signer.Secret).ToEqual(testSecret)
}

func TestSignURL(t *testing.T) {
	spec := gspec.New(t)

	now, err := time.Parse(time.RFC3339, "2015-07-31T00:00:00Z")
	if err != nil {
		panic("date time parse failed")
	}
	time.Freeze(now)
	defer time.Unfreeze()

	// Precomputed with node
	expectedSig := "3bcc03505a4e4f05ea8c9ba1979b788dcb35af94"
	expectedSignedURL := "https://api.viki.io/v4/movies.json?app=" + testAppId + "&sig=" + expectedSig + "&sort=views&t=" + strconv.FormatInt(now.Unix(), 10)

	host := "https://api.viki.io"
	path := "/v4/movies.json?sort=views"

	signer := New(testAppId, testSecret)
	spec.Expect(signer.Sign(host + path)).ToEqual(expectedSignedURL)
}

func TestSignPath(t *testing.T) {
	spec := gspec.New(t)

	now, err := time.Parse(time.RFC3339, "2015-07-31T00:00:00Z")
	if err != nil {
		panic("date time parse failed")
	}
	time.Freeze(now)
	defer time.Unfreeze()

	// Precomputed with node
	expectedSig := "3bcc03505a4e4f05ea8c9ba1979b788dcb35af94"
	expectedSignedPath := "/v4/movies.json?app=" + testAppId + "&sig=" + expectedSig + "&sort=views&t=" + strconv.FormatInt(now.Unix(), 10)

	path := "/v4/movies.json?sort=views"

	signer := New(testAppId, testSecret)
	spec.Expect(signer.Sign(path)).ToEqual(expectedSignedPath)
}
