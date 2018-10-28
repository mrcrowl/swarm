package web

import (
	"fmt"
	"net/http"
	"strings"
	"swarm/config"
	"swarm/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createWebServer(rootpath string) (*Server, *http.ServeMux) {
	fmt.Printf("Creating web server with root: %s\n", rootpath)
	config := config.NewServerConfig(9001, false, true)
	opts := CreateServerOptions(rootpath, config, nil, "app")
	server := CreateServer(opts)
	mux := http.NewServeMux()
	return server, mux
}

func TestURL(t *testing.T) {
	server, _ := createWebServer("")
	actual := server.URL()
	assert.Equal(t, "http://localhost:9001/app", actual)
}

func TestPort(t *testing.T) {
	server, _ := createWebServer("")
	actual := server.Port()
	assert.Equal(t, uint16(9001), actual)
}

func TestMakeServerAddress(t *testing.T) {
	actual := makeServerAddress(uint16(9001))
	assert.Equal(t, ":9001", actual)
}

func TestStaticFileServer(t *testing.T) {
	tempDir := testutil.CreateTempDir()
	defer testutil.RemoveTempDir(tempDir)
	styleDir := testutil.MakeSubdirectoryTree(tempDir, "one/two/three")
	testutil.WriteTextFile(styleDir, "styles.css", "body { background-color: hotpink; }")
	server, mux := createWebServer(tempDir)
	server.attachStaticFileServer(mux)

	// mock a request
	request, _ := http.NewRequest("GET", "/one/two/three/styles.css", nil)
	writer := newMockWriter()
	mux.ServeHTTP(writer, request)

	// check the response
	actual := writer.sb.String()
	expected := fmt.Sprintf("body { background-color: hotpink; }")
	assert.Equal(t, expected, actual)
	assert.Equal(t, "text/css; charset=utf-8", writer.ContentType())
}

func TestIndexInjectionListener(t *testing.T) {
	// configure files and server
	tempDir := testutil.CreateTempDir()
	defer testutil.RemoveTempDir(tempDir)
	appDir := testutil.MakeSubdirectoryTree(tempDir, "app")
	testutil.WriteTextFile(appDir, "index.html", "<body>HELLO WORLD</body>")
	testutil.WriteTextFile(appDir, "some-other.html", "<body>GOODBYE WORLD</body>")
	assetPath := testutil.MakeSubdirectoryTree(tempDir, assetsPhysicalPath)
	testutil.WriteTextFile(assetPath, hotReloadFilename, "WHATEVER")

	server, mux := createWebServer(tempDir)
	fileServer := http.FileServer(http.Dir(server.rootFilepath))
	server.attachIndexInjectionListener(mux, fileServer)

	cases := map[string]struct {
		url      string
		expected string
		mimetype string
	}{
		"index+inject": {
			url:      "/app/index.html",
			expected: fmt.Sprintf(`<body>HELLO WORLD<script type="module" src="%s/%s"></script></body>`, swarmVirtualPath, hotReloadFilename),
			mimetype: `text/html; charset=utf-8`,
		},
		"html-served-by-same-handler": {
			url:      "/app/some-other.html",
			expected: `<body>GOODBYE WORLD</body>`,
			mimetype: `text/html; charset=utf-8`,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// mock a request
			request, _ := http.NewRequest("GET", tc.url, nil)
			writer := newMockWriter()
			mux.ServeHTTP(writer, request)

			// check the response
			actual := writer.sb.String()
			expected := tc.expected
			assert.Equal(t, expected, actual)
			assert.Equal(t, tc.mimetype, writer.ContentType())
		})
	}
}

func TestHotReloadScripts(t *testing.T) {
	server, mux := createWebServer("c:\\")
	// configure files and server
	server.attachWebSocketListeners(mux, nil)

	for _, filename := range []string{hotReloadFilename, socketClientFilename} {
		// mock a request
		url := swarmify(filename)
		request, _ := http.NewRequest("GET", url, nil)
		writer := newMockWriter()
		mux.ServeHTTP(writer, request)

		// check the response
		actual := writer.sb.String()
		success := strings.HasPrefix(actual, "import ") || strings.HasPrefix(actual, "export ")
		assert.True(t, success, "%s doesn't start with import:\n--------------------------\n%s", url, actual)
		assert.Equal(t, "application/javascript", writer.ContentType())
	}
}

func TestSwarmify(t *testing.T) {
	actual := swarmify("bob")
	expected := fmt.Sprintf("%s/%s", swarmVirtualPath, "bob")
	assert.Equal(t, expected, actual)
}

func TestLoadAssetString(t *testing.T) {
	actual := loadAssetString("test-asset.js")
	assert.Equal(t, "alert('Hello world');", actual)
}

type MockWriter struct {
	sb      *strings.Builder
	headers map[string][]string
}

func newMockWriter() *MockWriter {
	return &MockWriter{
		sb:      &strings.Builder{},
		headers: make(map[string][]string),
	}
}

func (w *MockWriter) ContentType() string {
	if header, ok := w.headers["Content-Type"]; ok {
		return header[0]
	}
	return ""
}

func (w *MockWriter) Write(bytes []byte) (int, error) {
	w.sb.Write(bytes)
	return len(bytes), nil
}

func (w *MockWriter) Header() http.Header        { return w.headers }
func (w *MockWriter) WriteHeader(statusCode int) {}
