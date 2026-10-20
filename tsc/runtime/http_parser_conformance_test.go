package tsruntime

import (
	"testing"
	"unsafe"
)

type httpParserCapture struct {
	headerArgs []TsValue
	completed  int
}

func newHTTPParserForTest(t *testing.T, methods []string, capture *httpParserCapture) (*TsLoop, TsValue, *TsNodeHTTPParserState) {
	t.Helper()
	loop := TsNewLoop()
	TsNodeHTTPCommonSingleton.States = make(map[uintptr]*TsNodeHTTPParserState)
	TsNodeHTTPCommonSingleton.AllMethods = methods

	object := TsNewObject()
	parser := TsObjectValue(object)
	state := &TsNodeHTTPParserState{Initialized: true, ParserType: 1, Connections: TsU}
	TsNodeHTTPCommonSingleton.States[uintptr(unsafe.Pointer(object))] = state
	TsNodeHTTPCommonSet(loop, parser, "_headers", TsArrayValue(&TsArray{}))
	TsNodeHTTPCommonSet(loop, parser, "_url", TsNodeHTTPCommonKey(""))
	TsNodeHTTPCommonSet(loop, parser, "maxHeaderPairs", TsNumberValue(2000))
	object.Set("2", TsFunctionValue(TsFunc(func(_ *TsLoop, args ...TsValue) TsValue {
		capture.headerArgs = append([]TsValue(nil), args...)
		return TsBooleanValue(true)
	})))
	object.Set("4", TsFunctionValue(TsFunc(func(_ *TsLoop, _ ...TsValue) TsValue {
		capture.completed++
		return TsU
	})))
	t.Cleanup(func() { delete(TsNodeHTTPCommonSingleton.States, uintptr(unsafe.Pointer(object))) })
	return loop, parser, state
}

func TestHTTP09GETSupportsFragmentedRequestLine(t *testing.T) {
	capture := &httpParserCapture{}
	loop, parser, _ := newHTTPParserForTest(t, []string{"GET", "POST"}, capture)
	first := []byte("GET /health\r")
	if got := TsNodeHTTPCommonParserExecute(loop, parser, TsNodeBuffer(first)); got.Kind != TsNumberKind || int(TsNumber(got)) != len(first) {
		t.Fatalf("first fragment consumed %v; want %d", got, len(first))
	}
	if len(capture.headerArgs) != 0 {
		t.Fatal("incomplete request line triggered headers-complete")
	}
	second := []byte("\n")
	if got := TsNodeHTTPCommonParserExecute(loop, parser, TsNodeBuffer(second)); got.Kind != TsNumberKind || int(TsNumber(got)) != len(second) {
		t.Fatalf("second fragment consumed %v; want %d", got, len(second))
	}
	if len(capture.headerArgs) != 9 {
		t.Fatalf("headers-complete received %d arguments; want 9", len(capture.headerArgs))
	}
	if got := int(TsNumber(capture.headerArgs[0])); got != 0 {
		t.Errorf("HTTP major version = %d; want 0", got)
	}
	if got := int(TsNumber(capture.headerArgs[1])); got != 9 {
		t.Errorf("HTTP minor version = %d; want 9", got)
	}
	if got := TsText(capture.headerArgs[4]); got != "/health" {
		t.Errorf("request target = %q; want /health", got)
	}
	if TsTruthy(capture.headerArgs[8]) {
		t.Error("HTTP/0.9 request unexpectedly permits keep-alive")
	}
	if capture.completed != 1 {
		t.Errorf("message-complete count = %d; want 1", capture.completed)
	}
}

func TestHTTP09RejectsNonGETAndHeaders(t *testing.T) {
	for _, request := range []struct {
		name string
		wire string
		code string
	}{
		{name: "non-GET", wire: "POST /submit\r\n", code: "HPE_INVALID_METHOD"},
		{name: "headers", wire: "GET / HTTP/0.9\r\nX-Test: value\r\n\r\n", code: "HPE_INVALID_HEADER_TOKEN"},
	} {
		t.Run(request.name, func(t *testing.T) {
			capture := &httpParserCapture{}
			loop, parser, _ := newHTTPParserForTest(t, []string{"GET", "POST"}, capture)
			got := TsNodeHTTPCommonParserExecute(loop, parser, TsNodeBuffer([]byte(request.wire)))
			if got.Kind != TsErrorKind {
				t.Fatalf("request returned %v; want parser error", got)
			}
			if code := TsText(TsGet(loop, got, TsNodeHTTPCommonKey("code"))); code != request.code {
				t.Errorf("parser error code = %q; want %s", code, request.code)
			}
			if len(capture.headerArgs) != 0 {
				t.Fatal("invalid HTTP/0.9 request reached headers-complete")
			}
		})
	}
}

func TestHTTPParserRejectsAmbiguousContentLengthAndTransferEncoding(t *testing.T) {
	for _, wire := range []string{
		"POST /upload HTTP/1.1\r\nContent-Length: 1\r\nTransfer-Encoding: chunked\r\n\r\n0\r\n\r\n",
		"POST /upload HTTP/1.1\r\nContent-Length: 1\r\nContent-Length: 1\r\n\r\nx",
	} {
		capture := &httpParserCapture{}
		loop, parser, _ := newHTTPParserForTest(t, []string{"GET", "POST"}, capture)
		got := TsNodeHTTPCommonParserExecute(loop, parser, TsNodeBuffer([]byte(wire)))
		if got.Kind != TsErrorKind {
			t.Errorf("ambiguous request framing returned %v; want parser error", got)
		}
		if len(capture.headerArgs) != 0 {
			t.Error("ambiguous request framing reached headers-complete")
		}
	}
}

func TestHTTPParserRejectsDuplicateChunkedTransferCoding(t *testing.T) {
	capture := &httpParserCapture{}
	loop, parser, _ := newHTTPParserForTest(t, []string{"GET", "POST"}, capture)
	wire := "POST /upload HTTP/1.1\r\nTransfer-Encoding: chunked, chunked\r\n\r\n0\r\n\r\n"
	got := TsNodeHTTPCommonParserExecute(loop, parser, TsNodeBuffer([]byte(wire)))
	if got.Kind != TsErrorKind {
		t.Fatalf("duplicate chunked coding returned %v; want parser error", got)
	}
	if code := TsText(TsGet(loop, got, TsNodeHTTPCommonKey("code"))); code != "HPE_INVALID_TRANSFER_ENCODING" {
		t.Errorf("parser error code = %q; want HPE_INVALID_TRANSFER_ENCODING", code)
	}
	if len(capture.headerArgs) != 0 {
		t.Fatal("invalid transfer coding reached headers-complete")
	}
}
