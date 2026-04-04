package output

import (
  "errors"
  "fmt"
  "html"
  "net/http"
  "io"
  "io/ioutil"
  "os"
  "regexp"
  "sort"
  "strconv"
  "strings"
  "unicode"
  "net/url"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(^uint(0) >> 1)
const MinInt = -(MaxInt - 1)

var Contains = strings.Contains
var ContainsAny = strings.ContainsAny
var Fields = strings.Fields
var Find = strings.Index
var HasPrefix = strings.HasPrefix
var IsDigit = unicode.IsDigit
var Join = strings.Join
var Lower = strings.ToLower
var Prefix = strings.HasPrefix
var Replace = strings.ReplaceAll
var Sort = sort.SliceStable
var Split = strings.Split
var Str = fmt.Sprint
var Strln = fmt.Sprintln
var Strf = fmt.Sprintf
var Suffix = strings.HasSuffix
var Upper = strings.ToUpper

var Getenv = os.Getenv
var Setenv = os.Setenv
var Environ = os.Environ

const NL = "\n"

func QueryUnescape(s string) string { x, _ := url.QueryUnescape(s); return x }
func UnescapeString(s string) string { return html.UnescapeString(s) }

func TrimClean(s string) string { return Join(Fields(s), " ") }

type Stringer interface {
    String() string
}

type Builder strings.Builder

func (b *Builder)unwrap() *strings.Builder { return (*strings.Builder)(b) }
func (b *Builder)Len() int { return b.unwrap().Len() }
func (b *Builder)String() string { return b.unwrap().String() }
func (b *Builder)Write(p []byte) (int, error) { return b.unwrap().Write(p) }
func (b *Builder)Add(items ...interface{}) *Builder {
    fmt.Fprint(b.unwrap(), items...)
    return b
}
func (b *Builder)CSV(items ...interface{}) {
	for _, item := range items {
		b.Add(DQ(Str(item)), ",")
	}
}

func Q(items ...interface{}) string { return `"`+Str(items...)+`"` }

func Min(a, b int) int { if a <= b { return a }; return b }
func Max(a, b int) int { if a >= b { return a }; return b }

type Writer_t struct { w io.Writer }
func Writer(w io.Writer) Writer_t { return Writer_t{ w:w } }
func (x Writer_t)Add(items ...interface{}) (int, error) { return fmt.Fprint(x.w, items...) }
func (x Writer_t)Write(p []byte) (int, error) { return x.w.Write(p) }

var Stdout, Stderr = Writer(os.Stdout), Writer(os.Stderr)

func Out(items ...interface{}) { fmt.Println(Str(items...)) }

func Error(items ...interface{}) error { return errors.New(Str(items...)) }

//func IfInt(t bool, a int, b ...int) (C int) { if t { return a }; l := len(b); if l == 0 { return }; return b[l-1] }
//func IfString(t bool, a string, b ...string) (C string) { if t { return a }; l := len(b); if l == 0 { return }; return b[l-1] }
//func If(t bool, a interface{}, b interface{}) interface{} { if t { return a }; return b }
func If[T any](t bool, a T, b ...T) T { if t { return a }; if len(b) == 0 { var z T; return z; }; return b[0] }

func Atoi(s string) int { a, _ := strconv.Atoi(Trim(s)); return a }
func Trim(s string) string { return strings.TrimSpace(s) }
func Blank(s string) bool { return len(Trim(s)) == 0 }

func DQ(in string) string {
    // Replace all double-quote runes with preceding backslant & wrap in double-quotes
    var b Builder
    b.Add("\"")
    b.Add(strings.Replace(in, "\"", "\\\"", -1))
    b.Add("\"")
    return b.String()
}

func SQ(in string) string {
    // Replace all apostrophes runes with preceding backslant & wrap in apostrophes (aka single quotes)
    var b Builder
    b.Add("'")
    b.Add(Replace(in, "'", "\\'"))
    b.Add("'")
    return b.String()
}

func MatchWrap(pattern, value string) bool {
	ok, _ := regexp.MatchString(pattern, value)
	return ok
}

func RemoveFile(f string) { os.Remove(f) }

func SendFileToClient(w http.ResponseWriter, path, file, content string) (x string) {
    filepath := Str(path, "/", file)
    open, e := os.Open(filepath)
    defer open.Close()
    if e != nil { Log(e); return }
    fstat, e := open.Stat()
    if e != nil { Log(e); return }
    var size = Str(fstat.Size())
    w.Header().Set("Content-Disposition", "attachment; filename="+file)
    w.Header().Set("Content-Type", content)
    w.Header().Set("Content-Length", size);

    open.Seek(0, 0);
	io.Copy(w, open);
	return filepath
}

func GetFileFromClient(r *http.Request, savepath string) (fname string, e error) {
    r.ParseMultipartForm(1024*1024) // in practice, the largest files are < 64*1024
    file, handler, e := r.FormFile("uploadFile")
    defer file.Close()
    fname = savepath + "/" + handler.Filename
    saveFile, e := os.Create(fname)
    defer saveFile.Close()
    bytes, e := ioutil.ReadAll(file)
    x, e := saveFile.Write(bytes)
    Log("Ready to return", fname, x, e)
    return fname, nil
}

func IsUploadFile(req *http.Request) bool {
    return Prefix(req.Header.Get("Content-type"), "multipart/form-data;")
}

func GetRequestFile(req *http.Request, fileControlId, savepath string) (fname string, e error) {
    req.ParseMultipartForm(1024*1024) // in practice, the largest files are < 64*1024
    file, handler, e := req.FormFile(fileControlId)
    if e != nil { return fname, e }
    defer file.Close()
    fname = Str(savepath, "/", handler.Filename)
    saveFile, e := os.Create(fname)
    if e != nil { return fname, e }
    defer saveFile.Close()
    bytes, e := ioutil.ReadAll(file)
    _, e = saveFile.Write(bytes)
    return fname, e
}

func HasOnlyDigits(in string) bool {
    if len(in) == 0 { return false }
    for _, r := range in {
        if !IsDigit(r) { return false }
    }
    return true
}

func OnlyDigits(in string) int {
    var b Builder
    for _, r := range in { if IsDigit(r) { b.Add(string(r)) } }
    return Atoi(b.String())
}

func OnlyAlpha(in string) string {
    var b Builder
    for _, r := range in { if unicode.IsLetter(r) { b.Add(string(r)) } }
    return b.String()
}

func Left(s string, n int) string {
    if n == 0 { n = MaxInt }
    l := len(s); if l < n { n=l }; return s[:n]
}

func Right(s string, n int) string {
    if n == 0 { n = MaxInt }
    l := len(s); if l < n { n=l }; return s[l-n:]
}

func Mid(s string, b, n int) string {
    l := len(s)
    b--
	if n < 0 || b < 0 || l == 0 || b > l { return "" }
	return Left(Right(s, l - b), n)
}

func Bit(b bool) int { if b { return 1 }; return 0 }

func StrSpace(items ...interface{}) string {
    var b Builder
    max := len(items) - 1
    for k, item := range items {
        s := Str(item)
        if len(s) == 0 { continue }
        b.Add(s)
        if k < max { b.Add(" ") }
    }
    return b.String()
}

func HeadTail(s string) (h, t string) {
	s = strings.Trim(s," \n\t")
	if len(s) == 0 { return }
	parts := strings.Split(s," ")
	h = parts[0]
	t = strings.TrimLeft(s[len(h):], " \n\t")
	return
}

func InitCap(s string) string {
	if len(s) == 0 { return s }
	r := []rune(s)
	return Str(Upper(string(r[0])), string(r[1:]))
}

func In(candidate int, matches ...int) bool {
    for _, m := range matches {
        if candidate == m { return true }
    }
    return false
}

func True(v interface{}) bool { // arrays, pointers, and non-basic types always false
	switch v.(type) {
		case string: return v.(string) != ""
		case int: return v.(int) != 0
		case int8: return v.(int8) != 0
		case int16: return v.(int16) != 0
		case int32: return v.(int32) != 0
		case int64: return v.(int64) != 0
		case uint: return v.(uint) != 0
		case uint8: return v.(uint8) != 0
		case uint16: return v.(uint16) != 0
		case uint32: return v.(uint32) != 0
		case uint64: return v.(uint64) != 0
		case float32: return v.(float32) != 0
		case float64: return v.(float64) != 0
		case complex64: return v.(complex64) != complex64(0)
		case complex128: return v.(complex128) != complex128(0)
		case bool: return v.(bool)
	}
	return false
}

func Attrib(attrib string, values ...interface{}) string {
    // formats an html attribute, with leading space
	// with blank attrib returns ""
    // with blank values returns " <attrib>" (leading space -- aka boolean attribute)
	if len(attrib) == 0 { return "" }
	var list []string
	var s string
	for _, v := range values {
		if v == nil { continue }
		s = Str(v)
		if len(s) == 0 { continue }
		list = append(list, s)
	}
	L := len(list)
	if L == 0 { return Str(" ", attrib) }
	if L == 1 { return Str(" ", attrib,`="`, list[0], `"`) }
	return Str(" ", attrib,`="`, Join(list, " "), `"`)
}

func AttribIf(attrib string, value interface{}) string {
    // formats an html attribute, with leading space
    // only if attrib and value are both non-zero
	s := Str(value)
	if len(attrib) == 0 || len(s) == 0 { return "" }
	return Str(" ", attrib, `="`, s, `"`)
}

type UrlValues url.Values

func (f UrlValues)GetSafe(name, alt string) string {
	val, ok := f[name]
	if !ok || len(val) == 0 { return alt }
	return val[0]
}

func (f UrlValues)Get(name string) string {
	val, ok := f[name]
	if !ok || len(val) == 0 { return "" }
	return val[0]
}

func ReadGetFile(source string) (blank string) {
	resp, e := http.Get(source)
	if e != nil { return }
	defer resp.Body.Close()
	bytes, e := ioutil.ReadAll(resp.Body)
	if e != nil { return }
	return string(bytes)
}

func (f UrlValues)FirstValues() (m map[string]string) {
    m = make(map[string]string)
    for k, v := range f { m[k] = v[0] }
    return m
}

func AppendWriteFile(fname string) (*os.File, error) {
    const mode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
    f, e := os.OpenFile(fname, mode, 0644)
    return f, e
}

func LogStdout(items ...interface{}) { fmt.Fprintln(os.Stdout, items...) }
func LogStderr(items ...interface{}) { fmt.Fprintln(os.Stderr, items...) }
func LogString(items ...interface{}) string { return fmt.Sprintln(items...) }
func LogNull(items ...interface{}) {}
var logTarget = os.Stderr
var Log = func (items ...interface{}) { fmt.Fprintln(logTarget, items...) }
func LogTarget(f *os.File) { logTarget = f }
