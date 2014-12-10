package rest

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/oleiade/reflections"
)

//more work to support http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html

func GetHeaderValue(key string, headers map[string][]string) string {
	values := headers[key]
	if values != nil && len(values) > 0 {
		return values[0]
	}
	return ""
}

// currently used primarily for testing, however, could be used to make request to the
// local process for cases where it made more sense to delpoy a typically remote resource locally
func LocalRequest(cr *ClientRequest) (*Request, *Response) {
	response := NewResponse()
	client := NewClient()
	if rawReq, err := client.NewHttpRequest(cr); err != nil {
		panic(err)
	} else {
		request := NewRequest(rawReq, NewRequestContext(), cr.Definition)
		request.Body = cr.Body
		request.Args = cr.args
		return request, response
	}
}

func FullApiWithHandlers(apiSpec interface{}) map[ServerResource]Handler {
	handlers := make(map[ServerResource]Handler)

	items, err := reflections.Items(apiSpec)
	if err != nil {
		panic(err)
	}

	for _, v := range items {
		if target, ok := v.(*ResourceSpec); ok {
			verbs, handler := target.ServeAll()
			if handler == nil {
				panic(fmt.Sprintf("no handler for endpoint %v ", target))
			}
			for _, verb := range verbs {
				handlers[verb] = handler
			}
		}
	}
	return handlers
}

// Bytes, is a helper method to reduce the number of lines to get a byte array out of the
// EndpointResponse
func Bytes(resp *EndpointResponse, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("no response from endpoint or code incorrectly called.")
	}

	if resp.HttpResponse == nil {
		return nil, errors.New("No HttpResponse Response")
	}
	if resp.HttpResponse.StatusCode != http.StatusOK {
		return nil, errors.New(resp.HttpResponse.Status)
	}

	if bytes, err := resp.Bytes(); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

// Encode - Returns a URL that is encoded in a way that is actually used in practice, preseves as many special characters
// as possible. The default encoding from GO makes no sense as it encodes even valid values such as + , :
func Encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := escape(k, encodePath) + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(escape(v, encodeQueryComponent))
		}
	}
	path := buf.String()
	u := &url.URL{
		Path: path,
	}
	return u.RequestURI()
}

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent:
			t[j] = '+'
			j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

type encoding int

const (
	encodePath encoding = 1 + iota
	encodeUserPassword
	encodeQueryComponent
	encodeFragment
)

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
// When 'all' is true the full range of reserved characters are matched.
func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case encodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last two as well. That leaves only ? to escape.
			return c == '?'

		case encodeUserPassword: // §3.2.2
			// The RFC allows ; : & = + $ , in userinfo, so we must escape only @ and /.
			// The parsing of userinfo treats : as special so we must escape that too.
			return c == '@' || c == '/' || c == ':'

		case encodeQueryComponent: // §3.4
			//screw RFC, in practice many api's require these characters including google cse
			return c == '?'

		case encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	// Everything else must be escaped.
	return true
}
