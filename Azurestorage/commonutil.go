package commonutil

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/log"
	"golang.org/x/net/html/charset"
)

//var log = logger.GetLogger("Azurestorage")
var activityLog = log.ChildLogger(log.RootLogger(), "azure-storage-activity")

// GetAzureStorageAPIpath forms azure storage api path
func GetAzureStorageAPIpath(accountName string, service string, operation string, paramMap map[string]string) (apiPath string) {
	azstorageAPIPath := "https://" + accountName + "." + service + ".core.windows.net"
	azstorageAPIRelPath := ""
	switch service {
	case "File":
		switch operation {
		case "List Shares":
			azstorageAPIRelPath = azstorageAPIRelPath + "/?comp=list"
			if !(paramMap["prefix"] == "" || len(paramMap["prefix"]) < 1) {
				azstorageAPIRelPath = azstorageAPIRelPath + "&prefix=" + paramMap["prefix"]
			}
			if !(paramMap["maxresults"] == "" || len(paramMap["maxresults"]) < 1) {
				azstorageAPIRelPath = azstorageAPIRelPath + "&maxresults=" + paramMap["maxresults"]
			}
			if !(paramMap["nextmarker"] == "" || len(paramMap["nextmarker"]) < 1) {
				azstorageAPIRelPath = azstorageAPIRelPath + "&marker=" + paramMap["nextmarker"]
			}
		case "List Directories and Files":
			azstorageAPIRelPath = azstorageAPIRelPath + "/" + paramMap["shareName"]
			if !(paramMap["directoryPath"] == "" || len(paramMap["directoryPath"]) < 1) {
				azstorageAPIRelPath = azstorageAPIRelPath + "/" + paramMap["directoryPath"]
			}
			azstorageAPIRelPath = azstorageAPIRelPath + "/?restype=directory&comp=list"
			if !(paramMap["prefix"] == "" || len(paramMap["prefix"]) < 1) {
				azstorageAPIRelPath = azstorageAPIRelPath + "&prefix=" + paramMap["prefix"]
			}
			if !(paramMap["maxresults"] == "" || len(paramMap["maxresults"]) < 1) {
				azstorageAPIRelPath = azstorageAPIRelPath + "&maxresults=" + paramMap["maxresults"]
			}
			if !(paramMap["nextmarker"] == "" || len(paramMap["nextmarker"]) < 1) {
				azstorageAPIRelPath = azstorageAPIRelPath + "&marker=" + paramMap["nextmarker"]
			}
		case "Get File", "Create File", "Delete File", "Write Content", "Clear Content":
			azstorageAPIRelPath = azstorageAPIRelPath + "/" + paramMap["shareName"]
			if !(paramMap["directoryPath"] == "" || len(paramMap["directoryPath"]) < 1) {
				azstorageAPIRelPath = azstorageAPIRelPath + "/" + paramMap["directoryPath"]
			}
			azstorageAPIRelPath = azstorageAPIRelPath + "/" + paramMap["fileName"]
		case "Create Share", "Delete Share":
			azstorageAPIRelPath = azstorageAPIRelPath + "/" + paramMap["shareName"] + "?restype=share"
		case "Create Directory", "Delete Directory":
			azstorageAPIRelPath = azstorageAPIRelPath + "/" + paramMap["shareName"]
			if !(paramMap["directoryPath"] == "" || len(paramMap["directoryPath"]) < 1) {
				azstorageAPIRelPath = azstorageAPIRelPath + "/" + paramMap["directoryPath"]
			}
			azstorageAPIRelPath = azstorageAPIRelPath + "/" + paramMap["directoryName"] + "?restype=directory"
		}
		break
	}
	azstorageAPIPath = azstorageAPIPath + azstorageAPIRelPath
	activityLog.Info("azstorageAPIPath ", azstorageAPIPath)
	return azstorageAPIPath

}

// DoCall is a private implementation helper for making HTTP calls
func DoCall(method string, operation string, path string, paramMap map[string]string, token string) (response *http.Response, err error) {
	//log.RootLogger().Info("Invoking azure storage backend")
	var req *http.Request
	headers := make(map[string]string)
	token = token[1:len(token)]
	if operation == "Create File" {
		requestBody := paramMap["fileContent"]
		data, _ := b64.StdEncoding.DecodeString(requestBody)
		contentLength := len(data)
		contentLengthinString := strconv.Itoa(contentLength)
		req, err = http.NewRequest(method, path+"?"+token, nil)
		req.Header.Add("x-ms-type", "file")
		req.Header.Add("x-ms-content-length", contentLengthinString)
		req.Header.Add("x-ms-version", "2014-02-14 ")
		req.Header.Add("x-ms-date", string(time.Now().Unix()))
		client := &http.Client{}
		response, err = client.Do(req)
		if err != nil {
			return response, fmt.Errorf("http invocation failed %s, status code %v", err.Error(), response.StatusCode)
		}
		if response.StatusCode != 201 {
			return response, nil
		}
		activityLog.Info("Status is ", response.Status)
		if !(requestBody == "" || contentLength < 1) {
			defer response.Body.Close()
			xmsrange := strconv.Itoa(contentLength - 1)
			xmsrange = "bytes=0-" + xmsrange
			response1, err := WriteClearFile(operation, requestBody, path, token, xmsrange)
			if err != nil {
				return response1, fmt.Errorf("http invocation failed %s, status code %v", err.Error(), response.StatusCode)
			}
			activityLog.Info(response1.StatusCode, response1.Status)
			return response1, nil
		}

	} else if operation == "Write Content" {
		requestBody := paramMap["fileContent"]
		startrange := paramMap["startRange"]
		endrange := paramMap["endRange"]
		err := ValidateRange(startrange, endrange)
		if err != nil {
			return nil, err
		}
		xmsrange := "bytes="
		if len(startrange) > 0 && len(endrange) > 0 {
			xmsrange = xmsrange + startrange + "-" + endrange

		} else {
			// appending content at end of file
			data, _ := b64.StdEncoding.DecodeString(requestBody)
			contentLength := len(data)
			response, err = Call("GET", path+"?"+token, headers)
			if response.StatusCode != 200 {
				return response, nil
			}
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			filecontentLength := len(string(body))
			startRange := strconv.Itoa(filecontentLength)
			endRange := strconv.Itoa(filecontentLength + contentLength - 1)
			xmsrange = xmsrange + startRange + "-" + endRange
			//resize the file
			headers = map[string]string{"x-ms-content-length": strconv.Itoa(filecontentLength + contentLength)}
			responseResize, err := Call("PUT", path+"?comp=properties"+"&"+token, headers)
			if err != nil {
				return responseResize, fmt.Errorf("http invocation failed %s, status code %v", err.Error(), responseResize.StatusCode)
			}
			defer responseResize.Body.Close()
		}
		response, err = WriteClearFile(operation, requestBody, path, token, xmsrange)
		if err != nil {
			return response, fmt.Errorf("http invocation failed %s, status code %v", err.Error(), response.StatusCode)
		}
		activityLog.Info(response.StatusCode, response.Status)

	} else if operation == "Clear Content" {
		startrange := paramMap["startRange"]
		endrange := paramMap["endRange"]
		err := ValidateRange(startrange, endrange)
		if err != nil {
			return nil, err
		}
		xmsrange := "bytes="
		if len(startrange) > 0 && len(endrange) > 0 {
			xmsrange = xmsrange + startrange + "-" + endrange
		} else {
			// appending content at end of file
			response, err = Call("GET", path+"?"+token, headers)
			if response.StatusCode != 200 {
				return response, nil
			}
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			filecontentLength := len(string(body))
			if filecontentLength <= 0 {
				return nil, fmt.Errorf("File is empty, Clearing Content is not possible on empty file")
			}
			endRange := strconv.Itoa(filecontentLength - 1)
			xmsrange = xmsrange + "0" + "-" + endRange
		}
		response, err = WriteClearFile(operation, "", path, token, xmsrange)
		if err != nil {
			return response, fmt.Errorf("http invocation failed %s, status code %v", err.Error(), response.StatusCode)
		}
		activityLog.Info(response.StatusCode, response.Status)
		//defer response.Body.Close()
		//resize the file -- setting content length to zero
		// headers = map[string]string{"x-ms-content-length": "0"}
		// response, err = Call("PUT", path+"?comp=properties"+"&"+token, headers)
		// if err != nil {
		// 	return response, fmt.Errorf("http invocation failed %s, status code %v", err.Error(), response.StatusCode)
		// }

	} else {
		apiPath := path + "&" + token
		if operation == "Get File" || operation == "Delete File" {
			apiPath = path + "?" + token
		}
		//	fmt.Println(apiPath)
		response, err = Call(method, apiPath, headers)
		activityLog.Info(" Status is ", response.Status)
		if err != nil {
			return response, fmt.Errorf("http invocation failed %s, status code %v", err.Error(), response.StatusCode)
		}
	}
	return response, nil
}

// WriteClearFile write or clear bytes on a file
func WriteClearFile(operation string, requestBody string, path string, token string, xmsrange string) (responsewc *http.Response, err error) {
	var req *http.Request
	path = path + "?comp=range" + "&" + token
	if operation == "Clear Content" {
		req, _ = http.NewRequest("PUT", path, nil)
		req.Header.Add("x-ms-range", xmsrange)
		req.Header.Set("x-ms-write", "clear")
		req.Header.Add("x-ms-version", "2014-02-14 ")
		req.Header.Add("x-ms-date", string(time.Now().Unix()))
	} else {
		data, _ := b64.StdEncoding.DecodeString(requestBody)
		req, _ = http.NewRequest("PUT", path, bytes.NewReader(data))
		req.Header.Add("x-ms-range", xmsrange)
		req.Header.Set("x-ms-write", "update")
		req.Header.Add("x-ms-version", "2014-02-14 ")
		req.Header.Add("x-ms-date", string(time.Now().Unix()))
	}
	responsewc, err = http.DefaultClient.Do(req)
	if err != nil {
		return responsewc, fmt.Errorf("http invocation failed %s, status code %v", err.Error(), responsewc.StatusCode)
	}
	return responsewc, nil
}

// Call azure service
func Call(method string, path string, headersMap map[string]string) (responsec *http.Response, err error) {
	req, err := http.NewRequest(method, path, nil)
	req.Header.Add("x-ms-version", "2014-02-14 ")
	req.Header.Add("x-ms-date", string(time.Now().Unix()))
	//	fmt.Println(path)
	for k := range headersMap {
		//	fmt.Printf("key[%s] value[%s]\n", k, headersMap[k])
		req.Header.Add(k, headersMap[k])
	}
	client := &http.Client{}
	responsec, err = client.Do(req)
	if err != nil {
		return responsec, fmt.Errorf("http invocation failed %s, status code %v", err.Error(), responsec.StatusCode)
	}
	return responsec, nil
}

// ErrorHandeler handels the errors
func ErrorHandeler(resp *http.Response, jsonResponseData []byte) (erro error) {
	res1, _ := Convert(strings.NewReader(fmt.Sprintf("%s", jsonResponseData)))
	errorResp := res1.String()
	errorMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(errorResp), &errorMap)
	if err != nil {
		activityLog.Error(err)
		return err
	}
	errorElements := errorMap["Error"].(map[string]interface{})
	errorMessage := errorElements["Message"].(string)
	activityLog.Error("Error while performing activity operation. Message from backend :  ", errorElements["Message"])
	if errorElements["Code"].(string) == "InvalidHeaderValue" {
		return activity.NewError(errorElements["Code"].(string)+" : "+errorElements["HeaderName"].(string)+" : "+errorElements["HeaderValue"].(string), "azure storage-1007", errorElements)
	}
	return activity.NewError(errorMessage, "azure storage-1007", errorElements)

}

// ValidateRange validates min and max range
func ValidateRange(startrange string, endrange string) (erro error) {
	if (len(startrange) < 1 && len(endrange) > 0) || (len(endrange) < 1 && len(startrange) > 0) {
		return fmt.Errorf("Either provide both range or keep both empty")
	} else if len(startrange) > 0 && len(endrange) > 0 {
		// if strings.Contains(startrange, ".") || strings.Contains(endrange, ".") {
		// 	return fmt.Errorf("range can not have float/decimal values")
		// }
		min, _ := strconv.Atoi(startrange)
		max, _ := strconv.Atoi(endrange)
		if min > max {
			return fmt.Errorf("start range can't be bigger than end range")
		}
	}
	return nil
}

// Convert converts the given XML document to JSON
func Convert(r io.Reader) (*bytes.Buffer, error) {
	// Decode XML document
	root := &Node{}
	err := NewDecoder(r).Decode(root)
	if err != nil {
		return nil, err
	}

	// Then encode it in JSON
	buf := new(bytes.Buffer)
	err = NewEncoder(buf).Encode(root)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

const (
	attrPrefix    = "-"
	contentPrefix = "#"
)

// A Decoder reads and decodes XML objects from an input stream.
type Decoder struct {
	r               io.Reader
	err             error
	attributePrefix string
	contentPrefix   string
}

type element struct {
	parent *element
	n      *Node
	label  string
}

// SetAttributePrefix sets attributePrefix on a decoder
func (dec *Decoder) SetAttributePrefix(prefix string) {
	dec.attributePrefix = prefix
}

// SetContentPrefix sets contentPrefix on a decoder
func (dec *Decoder) SetContentPrefix(prefix string) {
	dec.contentPrefix = prefix
}

// DecodeWithCustomPrefixes decodes using given custom prefixes
func (dec *Decoder) DecodeWithCustomPrefixes(root *Node, contentPrefix string, attributePrefix string) error {
	dec.contentPrefix = contentPrefix
	dec.attributePrefix = attributePrefix
	return dec.Decode(root)
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the next JSON-encoded value from its
// input and stores it in the value pointed to by v.
func (dec *Decoder) Decode(root *Node) error {

	if dec.contentPrefix == "" {
		dec.contentPrefix = contentPrefix
	}
	if dec.attributePrefix == "" {
		dec.attributePrefix = attrPrefix
	}

	xmlDec := xml.NewDecoder(dec.r)

	// That will convert the charset if the provided XML is non-UTF-8
	xmlDec.CharsetReader = charset.NewReaderLabel

	// Create first element from the root node
	elem := &element{
		parent: nil,
		n:      root,
	}

	for {
		t, _ := xmlDec.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			// Build new a new current element and link it to its parent
			elem = &element{
				parent: elem,
				n:      &Node{},
				label:  se.Name.Local,
			}

			// Extract attributes as children
			for _, a := range se.Attr {
				elem.n.AddChild(dec.attributePrefix+a.Name.Local, &Node{Data: a.Value})
			}
		case xml.CharData:
			// Extract XML data (if any)
			elem.n.Data = trimNonGraphic(string(xml.CharData(se)))
		case xml.EndElement:
			// And add it to its parent list
			if elem.parent != nil {
				elem.parent.n.AddChild(elem.label, elem.n)
			}

			// Then change the current element to its parent
			elem = elem.parent
		}
	}

	return nil
}

// trimNonGraphic returns a slice of the string s, with all leading and trailing
// non graphic characters and spaces removed.
//
// Graphic characters include letters, marks, numbers, punctuation, symbols,
// and spaces, from categories L, M, N, P, S, Zs.
// Spacing characters are set by category Z and property Pattern_White_Space.
func trimNonGraphic(s string) string {
	if s == "" {
		return s
	}

	var first *int
	var last int
	for i, r := range []rune(s) {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) {
			continue
		}

		if first == nil {
			f := i // copy i
			first = &f
			last = i
		} else {
			last = i
		}
	}

	// If first is nil, it means there are no graphic characters
	if first == nil {
		return ""
	}

	return string([]rune(s)[*first : last+1])
}

// An Encoder writes JSON objects to an output stream.
type Encoder struct {
	w               io.Writer
	err             error
	contentPrefix   string
	attributePrefix string
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// SetAttributePrefix sets attributePrefix on an Encoder
func (enc *Encoder) SetAttributePrefix(prefix string) {
	enc.attributePrefix = prefix
}

// SetContentPrefix sets contentPrefix on an Encoder
func (enc *Encoder) SetContentPrefix(prefix string) {
	enc.contentPrefix = prefix
}

// EncodeWithCustomPrefixes encodes with given custom prefixes
func (enc *Encoder) EncodeWithCustomPrefixes(root *Node, contentPrefix string, attributePrefix string) error {
	enc.contentPrefix = contentPrefix
	enc.attributePrefix = attributePrefix
	return enc.Encode(root)
}

// Encode writes the JSON encoding of v to the stream
func (enc *Encoder) Encode(root *Node) error {
	if enc.err != nil {
		return enc.err
	}
	if root == nil {
		return nil
	}
	if enc.contentPrefix == "" {
		enc.contentPrefix = contentPrefix
	}
	if enc.attributePrefix == "" {
		enc.attributePrefix = attrPrefix
	}

	enc.err = enc.format(root, 0)

	// Terminate each value with a newline.
	// This makes the output look a little nicer
	// when debugging, and some kind of space
	// is required if the encoded value was a number,
	// so that the reader knows there aren't more
	// digits coming.
	enc.write("\n")

	return enc.err
}

func (enc *Encoder) format(n *Node, lvl int) error {
	if n.IsComplex() {
		enc.write("{")

		// Add data as an additional attibute (if any)
		if len(n.Data) > 0 {
			enc.write("\"")
			enc.write(enc.contentPrefix)
			enc.write("content")
			enc.write("\": ")
			enc.write(sanitiseString(n.Data))
			enc.write(", ")
		}

		i := 0
		tot := len(n.Children)
		for label, children := range n.Children {
			enc.write("\"")
			enc.write(label)
			enc.write("\": ")

			if len(children) > 1 {
				// Array
				enc.write("[")
				for j, c := range children {
					enc.format(c, lvl+1)

					if j < len(children)-1 {
						enc.write(", ")
					}
				}
				enc.write("]")
			} else {
				// Map
				enc.format(children[0], lvl+1)
			}

			if i < tot-1 {
				enc.write(", ")
			}
			i++
		}

		enc.write("}")
	} else {
		// TODO : Extract data type
		enc.write(sanitiseString(n.Data))
	}

	return nil
}

func (enc *Encoder) write(s string) {
	enc.w.Write([]byte(s))
}

// https://golang.org/src/encoding/json/encode.go?s=5584:5627#L788
var hex = "0123456789abcdef"

func sanitiseString(s string) string {
	var buf bytes.Buffer

	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \n and \r,
				// as well as <, > and &. The latter are escaped because they
				// can lead to security holes when user-controlled strings
				// are rendered into JSON and served to some browsers.
				buf.WriteString(`\u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\u202`)
			buf.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
	return buf.String()
}

// Node is a data element on a tree
type Node struct {
	Children map[string]Nodes
	Data     string
}

// Nodes is a list of nodes
type Nodes []*Node

// AddChild appends a node to the list of children
func (n *Node) AddChild(s string, c *Node) {
	// Lazy lazy
	if n.Children == nil {
		n.Children = map[string]Nodes{}
	}

	n.Children[s] = append(n.Children[s], c)
}

// IsComplex returns whether it is a complex type (has children)
func (n *Node) IsComplex() bool {
	return len(n.Children) > 0
}

// GetBody returns a new bytes buffer from byte array of content
func GetBody(content interface{}) (io.Reader, error) {
	var reqBody io.Reader
	switch content.(type) {
	case string:
		reqBody = bytes.NewBuffer([]byte(content.(string)))
	default:
		b, err := json.Marshal(content)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(b)
	}
	return reqBody, nil
}
