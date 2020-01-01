package main

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

var (
	elemNameMap = map[string]string{
		"a":          "elem.Anchor",
		"abbr":       "elem.Abbreviation",
		"address":    "elem.Address",
		"area":       "elem.Area",
		"article":    "elem.Article",
		"aside":      "elem.ASide",
		"audio":      "elem.Audio",
		"b":          "elem.Bold",
		"base":       "elem.Base",
		"bdi":        "elem.BidirectionalIsolation",
		"bdo":        "elem.BidirectionalOverride",
		"blockquote": "elem.BlockQuote",
		"body":       "elem.Body",
		"br":         "elem.Break",
		"button":     "elem.Button",
		"canvas":     "elem.Canvas",
		"caption":    "elem.Caption",
		"cite":       "elem.Citation",
		"code":       "elem.Code",
		"col":        "elem.Column",
		"colgroup":   "elem.ColumnGroup",
		"data":       "elem.Data",
		"datalist":   "elem.DataList",
		"dd":         "elem.Description",
		"del":        "elem.DeletedText",
		"details":    "elem.Details",
		"dfn":        "elem.Definition",
		"dialog":     "elem.Dialog",
		"div":        "elem.Div",
		"dl":         "elem.DescriptionList",
		"dt":         "elem.DefinitionTerm",
		"em":         "elem.Emphasis",
		"embed":      "elem.Embed",
		"fieldset":   "elem.FieldSet",
		"figcaption": "elem.FigureCaption",
		"figure":     "elem.Figure",
		"footer":     "elem.Footer",
		"form":       "elem.Form",
		"h1":         "elem.Heading1",
		"h2":         "elem.Heading2",
		"h3":         "elem.Heading3",
		"h4":         "elem.Heading4",
		"h5":         "elem.Heading5",
		"h6":         "elem.Heading6",
		"header":     "elem.Header",
		"hgroup":     "elem.HeadingsGroup",
		"hr":         "elem.HorizontalRule",
		"i":          "elem.Italic",
		"iframe":     "elem.InlineFrame",
		"img":        "elem.Image",
		"input":      "elem.Input",
		"ins":        "elem.InsertedText",
		"kbd":        "elem.KeyboardInput",
		"label":      "elem.Label",
		"legend":     "elem.Legend",
		"li":         "elem.ListItem",
		"link":       "elem.Link",
		"main":       "elem.Main",
		"map":        "elem.Map",
		"mark":       "elem.Mark",
		"menu":       "elem.Menu",
		"menuitem":   "elem.MenuItem",
		"meta":       "elem.Meta",
		"meter":      "elem.Meter",
		"nav":        "elem.Navigation",
		"noframes":   "elem.NoFrames",
		"noscript":   "elem.NoScript",
		"object":     "elem.Object",
		"ol":         "elem.OrderedList",
		"optgroup":   "elem.OptionsGroup",
		"option":     "elem.Option",
		"output":     "elem.Output",
		"p":          "elem.Paragraph",
		"param":      "elem.Parameter",
		"picture":    "elem.Picture",
		"pre":        "elem.Preformatted",
		"progress":   "elem.Progress",
		"q":          "elem.Quote",
		"rp":         "elem.RubyParenthesis",
		"rt":         "elem.RubyText",
		"rtc":        "elem.RubyTextContainer",
		"ruby":       "elem.Ruby",
		"s":          "elem.Strikethrough",
		"samp":       "elem.Sample",
		"script":     "elem.Script",
		"section":    "elem.Section",
		"select":     "elem.Select",
		"slot":       "elem.Slot",
		"small":      "elem.Small",
		"source":     "elem.Source",
		"span":       "elem.Span",
		"strong":     "elem.Strong",
		"style":      "elem.Style",
		"sub":        "elem.Subscript",
		"summary":    "elem.Summary",
		"sup":        "elem.Superscript",
		"table":      "elem.Table",
		"tbody":      "elem.TableBody",
		"td":         "elem.TableData",
		"template":   "elem.Template",
		"textarea":   "elem.TextArea",
		"tfoot":      "elem.TableFoot",
		"th":         "elem.TableHeader",
		"thead":      "elem.TableHead",
		"time":       "elem.Time",
		"tr":         "elem.TableRow",
		"track":      "elem.Track",
		"u":          "elem.Underline",
		"ul":         "elem.UnorderedList",
		"var":        "elem.Variable",
		"video":      "elem.Video",
		"wbr":        "elem.WordBreakOpportunity",
	}
	propMap = map[string]string{
		"alt":         "prop.Alt",
		"autofocus":   "prop.Autofocus",
		"checked":     "prop.Checked",
		"disabled":    "prop.Disabled",
		"for":         "prop.For",
		"href":        "prop.Href",
		"id":          "prop.ID",
		"name":        "prop.Name",
		"placeholder": "prop.Placeholder",
		"src":         "prop.Src",
		"type":        "prop.Type",
		"value":       "prop.Value",
	}
	propBool = map[string]struct{}{
		"autofocus": struct{}{},
		"checked":   struct{}{},
		"disabled":  struct{}{},
	}
	inputTypes = map[string]struct{}{
		"button":         struct{}{},
		"checkbox":       struct{}{},
		"color":          struct{}{},
		"date":           struct{}{},
		"datetime":       struct{}{},
		"datetime-local": struct{}{},
		"email":          struct{}{},
		"file":           struct{}{},
		"hidden":         struct{}{},
		"image":          struct{}{},
		"month":          struct{}{},
		"number":         struct{}{},
		"password":       struct{}{},
		"radio":          struct{}{},
		"range":          struct{}{},
		"min":            struct{}{},
		"max":            struct{}{},
		"value":          struct{}{},
		"step":           struct{}{},
		"reset":          struct{}{},
		"search":         struct{}{},
		"submit":         struct{}{},
		"tel":            struct{}{},
		"text":           struct{}{},
		"time":           struct{}{},
		"url":            struct{}{},
		"week":           struct{}{},
	}
)

func tag(z *html.Tokenizer) string {
	b, _ := z.TagName()
	return strings.ToLower(string(b))
}

func attrs(z *html.Tokenizer, indent int) string {
	tab0 := strings.Repeat("\t", indent)
	tab1 := strings.Repeat("\t", indent+1)
	tab2 := strings.Repeat("\t", indent+2)
	res := []string{}
	for {
		key, val, more := z.TagAttr()
		k := string(key)
		v := string(val)
		if len(k) == 0 {
			if !more {
				break
			}
			continue
		}
		if len(v) == 0 {
			v = "true"
		}
		if k == "class" {
			classes := []string{}
			for _, s := range strings.Split(v, " ") {
				classes = append(classes, fmt.Sprintf("%q", s))
			}
			if len(classes) <= 4 {
				res = append(res, fmt.Sprintf("\n%svecty.Class(%s),", tab1, strings.Join(classes, ", ")))
			} else {
				res = append(res, fmt.Sprintf("\n%svecty.ClassMap{", tab1))
				for _, s := range classes {
					res = append(res, fmt.Sprintf("\n%s\t%s: true,", tab2, s))
				}
				res = append(res, fmt.Sprintf("\n%s},", tab1))
			}
		} else if prop, ok := propMap[k]; ok {
			if _, ok := propBool[k]; ok {
				res = append(res, fmt.Sprintf("\n%s%s(%s),", tab1, prop, v))
			} else {
				res = append(res, fmt.Sprintf("\n%s%s(%q),", tab1, prop, v))
			}
		} else {
			if _, ok := propBool[k]; ok {
				res = append(res, fmt.Sprintf("\n%s%s(%s),", tab1, prop, v))
			} else {
				res = append(res, fmt.Sprintf("\n%svecty.Property(%q, %q),", tab1, k, v))
			}
		}
		if !more {
			break
		}
	}
	if len(res) == 0 {
		return ""
	}
	return fmt.Sprintf("\n%svecty.Markup(%s\n%s),", tab0, strings.Join(res, ""), tab0)
}

func generate(w io.Writer, r io.Reader) (err error) {
	indent := 1
	z := html.NewTokenizer(r)
	for err == nil {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err = z.Err()
			if err == io.EOF {
				break
			}
			return
		case html.CommentToken:
		case html.DoctypeToken:
		case html.TextToken:
			tab := strings.Repeat("\t", indent)
			t := strings.TrimSpace(string(z.Text()))
			if len(t) > 0 {
				fmt.Fprintf(w, "\n%svecty.Text(%q),", tab, t)
			}
		case html.StartTagToken:
			if indent > 1 {
				fmt.Fprint(w, "\n")
			}
			tab := strings.Repeat("\t", indent)
			indent++
			a, t := attrs(z, indent), string(z.Text())
			e, ok := elemNameMap[tag(z)]
			if !ok {
				e = fmt.Sprintf("vecty.Tag(%q, ", tag(z))
			}
			fmt.Fprintf(w, "%s%s(%s", tab, e, a)
			if len(t) > 0 {
				fmt.Fprintf(w, "\n%svecty.Text(%q),", tab, t)
			}
		case html.SelfClosingTagToken:
			tab := strings.Repeat("\t", indent)
			indent++
			a := attrs(z, indent)
			e, ok := elemNameMap[tag(z)]
			if !ok {
				e = fmt.Sprintf("\n%svecty.Tag(%q, ", tab, tag(z))
			} else {
				e = e + "("
			}
			if len(a) > 0 {
				fmt.Fprintf(w, "\n%[1]s%s%s\n%[1]s),", tab, e, a)
			} else {
				fmt.Fprintf(w, "\n%[1]s%s),", tab, e)
			}
			indent--
		case html.EndTagToken:
			indent--
			tab := strings.Repeat("\t", indent)
			fmt.Fprintf(w, "\n%s)", tab)
			if indent > 1 {
				fmt.Fprint(w, ",")
			}
		}
	}
	return nil
}
