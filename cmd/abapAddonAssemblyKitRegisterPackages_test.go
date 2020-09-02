package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SAP/jenkins-library/pkg/abaputils"
	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func mockReader(path string) ([]byte, error) {
	var file []byte
	if path == "exists" {
		return file, nil
	}
	return file, errors.New("error reading the file")
}

// ********************* Test uploadSarFiles *******************
func TestUploadSarFiles(t *testing.T) {
	t.Run("test uploadSarFiles", func(t *testing.T) {
		repositories, conn := setupRepos("exists", planned, clMockRegisterPackages{})
		err := uploadSarFiles(repositories, conn, mockReader)
		assert.NoError(t, err)
	})
}

func TestUploadSarFilesInvalidInput(t *testing.T) {
	t.Run("test uploadSarFiles with missing file path", func(t *testing.T) {
		repositories, conn := setupRepos("", planned, clMockRegisterPackages{})
		err := uploadSarFiles(repositories, conn, mockReader)
		assert.Error(t, err)
	})
}

func TestUploadSarFilesNoFile(t *testing.T) {
	t.Run("test uploadSarFiles with missing file", func(t *testing.T) {
		repositories, conn := setupRepos("does_not_exist", planned, clMockRegisterPackages{})
		err := uploadSarFiles(repositories, conn, mockReader)
		assert.Error(t, err)
	})
}

func TestUploadSarFilesErrorUploading(t *testing.T) {
	t.Run("test uploadSarFiles with error during upload", func(t *testing.T) {
		c := clMockRegisterPackages{
			err: errors.New("Failure"),
		}
		repositories, conn := setupRepos("exists", planned, c)
		err := uploadSarFiles(repositories, conn, mockReader)
		assert.Error(t, err)
	})
}

// ********************* Test registerPackages *******************
func TestRegisterPackages(t *testing.T) {
	t.Run("test registerPackages", func(t *testing.T) {
		repositories, conn := setupRepos("Filepath", planned, clMockRegisterPackages{})
		repos, err := registerPackages(repositories, conn)
		assert.NoError(t, err)
		assert.Equal(t, string(locked), repos[0].Status)
	})
}

func TestRegisterPackagesReleased(t *testing.T) {
	t.Run("test registerPackages", func(t *testing.T) {
		repositories, conn := setupRepos("Filepath", released, clMockRegisterPackages{})
		repos, err := registerPackages(repositories, conn)
		assert.NoError(t, err)
		assert.Equal(t, string(released), repos[0].Status)
	})
}

func TestRegisterPackagesError(t *testing.T) {
	t.Run("test registerPackages with error", func(t *testing.T) {
		c := clMockRegisterPackages{
			err: errors.New("Failure"),
		}
		repositories, conn := setupRepos("Filepath", planned, c)
		repos, err := registerPackages(repositories, conn)
		assert.Error(t, err)
		assert.Equal(t, string(planned), repos[0].Status)
	})
}

// ********************* Test Setup *******************
func setupRepos(filePath string, status packageStatus, cl clMockRegisterPackages) ([]abaputils.Repository, connector) {
	repositories := []abaputils.Repository{
		{
			Name:           "/DRNMSPC/COMP01",
			VersionYAML:    "1.0.0",
			PackageName:    "SAPK-001AAINDRNMSPC",
			Status:         string(status),
			SarXMLFilePath: filePath,
		},
	}
	conn := new(connector)
	conn.Client = &cl
	conn.Header = make(map[string][]string)
	return repositories, *conn
}

// ********************* Mocking *******************

type clMockRegisterPackages struct {
	err       error
	errorbody string
}

func (c *clMockRegisterPackages) SetOptions(opts piperhttp.ClientOptions) {}

func (c *clMockRegisterPackages) SendRequest(method string, url string, bdy io.Reader, hdr http.Header, cookies []*http.Cookie) (*http.Response, error) {
	switch method {
	case "HEAD":
		return c.sendRequestHead()
	case "PUT":
		return c.sendRequestPut()
	case "POST":
		return c.sendRequestPost()
	}
	return nil, nil
}

func (c *clMockRegisterPackages) sendRequestHead() (*http.Response, error) {
	var body []byte
	header := http.Header{}
	header.Set("X-CSRF-Token", "myToken")
	body = []byte("")
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}, nil
}

func (c *clMockRegisterPackages) sendRequestPut() (*http.Response, error) {
	var body []byte
	if c.err != nil {
		body = []byte(c.errorbody)
		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, c.err
	}
	body = []byte("")
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}, nil
}

func (c *clMockRegisterPackages) sendRequestPost() (*http.Response, error) {
	var body []byte
	if c.err != nil {
		body = []byte(c.errorbody)
		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, c.err
	}
	body = []byte(responseRegisterPackagesPost)
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}, nil
}

// ********************* Testdata *******************

var responseRegisterPackagesPost = `{
    "d": {
        "__metadata": {
            "id": "https://W7Q.DMZWDF.SAP.CORP:443/odata/aas_ocs_package/OcsPackageSet('SAPK-001AAINDRNMSPC')",
            "uri": "https://W7Q.DMZWDF.SAP.CORP:443/odata/aas_ocs_package/OcsPackageSet('SAPK-001AAINDRNMSPC')",
            "type": "SSDA.AAS_ODATA_PACKAGE_SRV.OcsPackage"
        },
        "Name": "SAPK-001AAINDRNMSPC",
        "Type": "AOI",
        "Component": "/DRNMSPC/COMP01",
        "Release": "0001",
        "Level": "0000",
        "Status": "L",
        "Operation": "",
        "Namespace": "/DRNMSPC/",
        "Vendorid": "0000203069"
    }
}`