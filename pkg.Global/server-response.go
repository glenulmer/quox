package globals

import (
    . "klpm/lib/output"
    . "klpm/lib/wrapdb"
    "net/http"
)

/*
Public functions:
    OuterHTML(target string, html string) ClientMessage_i
    InnerHTML(target string, html string) ClientMessage_i
    Remove(target string, html string) ClientMessage_i
    RewriteTBody(tbodyID string, html string) ClientMessage_i
    RewriteRow(rowID string, html string) ClientMessage_i
    RemoveRow(rowID string) ClientMessage_i
    RewriteMsg(style style_t, body, title string) RewriteMsg_t
    MessageList() *messageList

MessageList methods:
    Append(...ClientMessage_i)
    Send(http.ResponseWriter)

Style constants:
    Info    
    Success 
    Warning
    Danger
*/

type style_t int
const Info, Success, Warning, Danger style_t = 0, 1, 2, 3

func (s style_t)String() string { 
    return [...]string{"is-info", "is-success", "is-warning", "is-danger"}[s]
}

type ClientMessage_i interface {
    jsonString() string
    String() string
}

type method_t int
const OuterHTML, InnerHTML, Remove method_t = 0, 1, 2
func (m method_t)String() string {
    return [...]string{"outerHTML", "innerHTML", "remove"}[m]
}

func RewriteHTML(method method_t, target string, htmls ...any) Rewrite_t {
/*    if len(htmlList) == 1 { html = htmlList[0] }
    switch htmls := htmls.(type) {
    case string: out = Str(htmls)
    case Stringer: out = htmls[0].String()
    }
*/
    return Rewrite_t{ method: method, target: `#`+target, html: Str(htmls...) }
}

type Rewrite_t struct {
    method method_t
    target string
    html   string
}

func (m Rewrite_t) jsonString() string {
    var result Builder
    result.Add(
        `{`,
        ` "kind": "rewrite", `,
        ` "target": "`, EscapeSelector(m.target), `", `,
        ` "content": "`, EscapeHTML(m.html), `", `,
        ` "method": "`, m.method, `"`,
        ` }`,
    )
    return result.String()
}

func (m Rewrite_t) String() string {
    var result Builder
    result.Add("HTML ", m.method.String(), ` `, m.target, ` = `, m.html)
    return result.String()
}

func Note(style style_t, body string) RewriteMsg_t {
    return RewriteMsg_t{ style:style.String(), body: body }
}

type RewriteMsg_t struct {
    style string
    body  string
}

func (m RewriteMsg_t) jsonString() string {
    var result Builder
    result.Add(`{ `,
        `"kind": "queue", `,
        `"body": "`, EscapeHTML(m.body), `", `,
        `"style": "`, string(m.style), `"`,
        `}`,
    )
    return result.String()
}

func (m RewriteMsg_t) String() string {
    var result Builder
    result.Add(" [", string(m.style), "]: ", m.body)
    return result.String()
}

type messageList struct {
    messages []ClientMessage_i
}

func MessageList() *messageList {
    return &messageList{make([]ClientMessage_i, 0)}
}

func (ml *messageList) Append(msgs ...ClientMessage_i) *messageList{
    ml.messages = append(ml.messages, msgs...)
    return ml
}

func (ml *messageList) jsonString() string {
    if len(ml.messages) == 0 {
        return `[]`
    }
    
    msgs := make([]string, len(ml.messages))
    for i, m := range ml.messages {
        msgs[i] = m.jsonString()
    }
    var result Builder
    result.Add(`[`, Join(msgs, ","), `]`)
    return result.String()
}

func (ml *messageList) Send(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(ml.jsonString()))
}

type SPResult_t struct {
    Errs, NRows, ID int
    Note string
}

func SendResponse(w http.ResponseWriter, list ...any) {
    messages := MessageList()

    for _, item := range list {
        switch v := item.(type) {
        case *Row_t:
            var res SPResult_t
            pack := v.Scan(&res.Errs, &res.NRows, &res.ID, &res.Note)
            q := pack.Query()

            switch {
            case pack.HasError():
                messages.Append(Note(Danger, pack.Message()+` `+q))
            case res.Errs > 0:
                messages.Append(Note(Warning, res.Note + ` ` + q))
            default:
                messages.Append(Note(Success, res.Note))
            }

        case ClientMessage_i:
            messages.Append(v)
        }
    }

    messages.Send(w)
}
