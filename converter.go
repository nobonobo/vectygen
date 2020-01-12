package main

import (
	"fmt"
	"io"
	"log"
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
	eventTypes = map[string]string{
		"afterprint":               "event.AfterPrint",
		"animationend":             "event.AnimationEnd",
		"animationiteration":       "event.AnimationIteration",
		"animationstart":           "event.AnimationStart",
		"appinstalled":             "event.ApplicationInstalled",
		"audioprocess":             "event.AudioProcess",
		"audioend":                 "event.AudioEnd",
		"audiostart":               "event.AudioStart",
		"beforeprint":              "event.BeforePrint",
		"beforeunload":             "event.BeforeUnload",
		"blocked":                  "event.Blocked",
		"blur":                     "event.Blur",
		"boundary":                 "event.Boundary",
		"cached":                   "event.Cached",
		"canplay":                  "event.CanPlay",
		"canplaythrough":           "event.CanPlayThrough",
		"change":                   "event.Change",
		"chargingchange":           "event.ChargingChange",
		"chargingtimechange":       "event.ChargingTimeChange",
		"checking":                 "event.Checking",
		"click":                    "event.Click",
		"close":                    "event.Close",
		"complete":                 "event.Complete",
		"compassneedscalibration":  "event.compassneedscalibration",
		"compositionend":           "event.CompositionEnd",
		"compositionstart":         "event.CompositionStart",
		"compositionupdate":        "event.CompositionUpdate",
		"contextmenu":              "event.ContextMenu",
		"copy":                     "event.Copy",
		"cut":                      "event.Cut",
		"DOMContentLoaded":         "event.DOMContentLoaded",
		"devicechange":             "event.DeviceChange",
		"devicelight":              "event.DeviceLight",
		"devicemotion":             "event.DeviceMotion",
		"deviceorientation":        "event.DeviceOrientation",
		"deviceproximity":          "event.DeviceProximity",
		"dischargingtimechange":    "event.DischargingTimeChange",
		"dblclick":                 "event.DoubleClick",
		"downloading":              "event.Downloading",
		"drag":                     "event.Drag",
		"dragend":                  "event.DragEnd",
		"dragenter":                "event.DragEnter",
		"dragleave":                "event.DragLeave",
		"dragover":                 "event.DragOver",
		"dragstart":                "event.DragStart",
		"drop":                     "event.Drop",
		"durationchange":           "event.DurationChange",
		"emptied":                  "event.Emptied",
		"end":                      "event.End",
		"endEvent":                 "event.EndEvent",
		"ended":                    "event.Ended",
		"error":                    "event.Error",
		"focus":                    "event.focus",
		"focusin":                  "event.FocusIn",
		"focusout":                 "event.FocusOut",
		"fullscreenchange":         "event.FullScreenChange",
		"fullscreenerror":          "event.FullScreenError",
		"gamepadconnected":         "event.GamepadConnected",
		"gamepaddisconnected":      "event.GamepadDisconnected",
		"gotpointercapture":        "event.GotPointerCapture",
		"hashchange":               "event.HashChange",
		"input":                    "event.Input",
		"invalid":                  "event.Invalid",
		"keydown":                  "event.KeyDown",
		"keypress":                 "event.KeyPress",
		"keyup":                    "event.KeyUp",
		"languagechange":           "event.LanguageChange",
		"levelchange":              "event.LevelChange",
		"load":                     "event.Load",
		"loadend":                  "event.LoadEnd",
		"loadstart":                "event.LoadStart",
		"loadeddata":               "event.LoadedData",
		"loadedmetadata":           "event.LoadedMetadata",
		"lostpointercapture":       "event.LostPointerCapture",
		"mark":                     "event.Mark",
		"message":                  "event.Message",
		"messageerror":             "event.MessageError",
		"mousedown":                "event.MouseDown",
		"mouseenter":               "event.MouseEnter",
		"mouseleave":               "event.MouseLeave",
		"mousemove":                "event.MouseMove",
		"mouseout":                 "event.MouseOut",
		"mouseover":                "event.MouseOver",
		"mouseup":                  "event.MouseUp",
		"nomatch":                  "event.NoMatch",
		"noupdate":                 "event.NoUpdate",
		"notificationclick":        "event.NotificationClick",
		"obsolete":                 "event.Obsolete",
		"offline":                  "event.Offline",
		"online":                   "event.Online",
		"open":                     "event.Open",
		"orientationchange":        "event.OrientationChange",
		"pagehide":                 "event.PageHide",
		"pageshow":                 "event.PageShow",
		"paste":                    "event.Paste",
		"pause":                    "event.Pause",
		"play":                     "event.Play",
		"playing":                  "event.Playing",
		"pointercancel":            "event.PointerCancel",
		"pointerdown":              "event.PointerDown",
		"pointerenter":             "event.PointerEnter",
		"pointerleave":             "event.PointerLeave",
		"pointerlockchange":        "event.PointerLockChange",
		"pointerlockerror":         "event.PointerLockError",
		"pointermove":              "event.PointerMove",
		"pointerout":               "event.PointerOut",
		"pointerover":              "event.PointerOver",
		"pointerup":                "event.PointerUp",
		"popstate":                 "event.PopState",
		"progress":                 "event.Progress",
		"push":                     "event.Push",
		"pushsubscriptionchange":   "event.PushSubscriptionChange",
		"ratechange":               "event.RateChange",
		"readystatechange":         "event.ReadyStateChange",
		"repeatEvent":              "event.RepeatEvent",
		"reset":                    "event.Reset",
		"resize":                   "event.Resize",
		"resourcetimingbufferfull": "event.ResourceTimingBufferFull",
		"result":                   "event.Result",
		"resume":                   "event.Resume",
		"SVGAbort":                 "event.SVGAbort",
		"SVGError":                 "event.SVGError",
		"SVGLoad":                  "event.SVGLoad",
		"SVGResize":                "event.SVGResize",
		"SVGScroll":                "event.SVGScroll",
		"SVGUnload":                "event.SVGUnload",
		"SVGZoom":                  "event.SVGZoom",
		"scroll":                   "event.Scroll",
		"seeked":                   "event.Seeked",
		"seeking":                  "event.Seeking",
		"select":                   "event.Select",
		"selectstart":              "event.SelectStart",
		"selectionchange":          "event.SelectionChange",
		"show":                     "event.Show",
		"slotchange":               "event.SlotChange",
		"soundend":                 "event.SoundEnd",
		"soundstart":               "event.SoundStart",
		"speechend":                "event.SpeechEnd",
		"speechstart":              "event.SpeechStart",
		"stalled":                  "event.Stalled",
		"start":                    "event.Start",
		"storage":                  "event.Storage",
		"submit":                   "event.Submit",
		"success":                  "event.Success",
		"suspend":                  "event.Suspend",
		"timeupdate":               "event.TimeUpdate",
		"timeout":                  "event.Timeout",
		"touchcancel":              "event.TouchCancel",
		"touchend":                 "event.TouchEnd",
		//"touchenter":               "event.touchenter",
		//"touchleave":               "event.touchleave",
		"touchmove":        "event.TouchMove",
		"touchstart":       "event.TouchStart",
		"transitionend":    "event.TransitionEnd",
		"unload":           "event.Unload",
		"updateready":      "event.UpdateReady",
		"upgradeneeded":    "event.UpgradeNeeded",
		"userproximity":    "event.UserProximity",
		"versionchange":    "event.VersionChange",
		"visibilitychange": "event.VisibilityChange",
		"voiceschanged":    "event.VoicesChanged",
		"volumechange":     "event.VolumeChange",
		//"vrdisplayconnected":       "event.vrdisplayconnected",
		//"vrdisplaydisconnected":    "event.vrdisplaydisconnected",
		//"vrdisplaypresentchange":   "event.vrdisplaypresentchange",
		"waiting": "event.Waiting",
		"wheel":   "event.Wheel",
	}
)

// Converter ...
type Converter struct {
	StdModules map[string]bool
	ExtModules map[string]bool
	Methods    map[string]string
	AppendCode []string
}

// New ...
func New() *Converter {
	return &Converter{
		StdModules: map[string]bool{},
		ExtModules: map[string]bool{
			"github.com/gopherjs/vecty": true,
		},
		Methods:    map[string]string{},
		AppendCode: []string{},
	}
}

func (c *Converter) tag(z *html.Tokenizer) string {
	b, _ := z.TagName()
	return strings.ToLower(string(b))
}

type attr struct {
	k string
	v string
}

func parseAttrs(z *html.Tokenizer) []attr {
	res := []attr{}
	for {
		key, val, more := z.TagAttr()
		k := string(key)
		v := string(val)
		res = append(res, attr{k: k, v: v})
		if !more {
			break
		}
	}
	return res
}

func (c *Converter) attrs(attrSlice []attr, indent int) string {
	tab0 := strings.Repeat("\t", indent)
	tab1 := strings.Repeat("\t", indent+1)
	tab2 := strings.Repeat("\t", indent+2)
	res := []string{}
	for _, attr := range attrSlice {
		k := attr.k
		v := attr.v
		if len(v) == 0 {
			v = "true"
		}
		if strings.HasPrefix(k, "@") {
			// event mapping
			name := k[1:]
			statement, ok := eventTypes[name]
			if !ok {
				log.Fatalln("unknown event:", name)
			}
			c.Methods[name] = v
			res = append(res, fmt.Sprintf("\n%s%s(c.%s),", tab1, statement, v))
			c.ExtModules["github.com/gopherjs/vecty/event"] = true
			continue
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
					res = append(res, fmt.Sprintf("\n%s%s: true,", tab2, s))
				}
				res = append(res, fmt.Sprintf("\n%s},", tab1))
			}
		} else if prop, ok := propMap[k]; ok {
			if _, ok := propBool[k]; ok {
				res = append(res, fmt.Sprintf("\n%s%s(%s),", tab1, prop, v))
			} else {
				res = append(res, fmt.Sprintf("\n%s%s(%q),", tab1, prop, v))
			}
			c.ExtModules["github.com/gopherjs/vecty/prop"] = true
		} else {
			if _, ok := propBool[k]; ok {
				res = append(res, fmt.Sprintf("\n%s%s(%s),", tab1, prop, v))
			} else {
				res = append(res, fmt.Sprintf("\n%svecty.Property(%q, %q),", tab1, k, v))
			}
		}
	}
	if len(res) == 0 {
		return ""
	}
	return fmt.Sprintf("\n%svecty.Markup(%s\n%s),", tab0, strings.Join(res, ""), tab0)
}

func (c *Converter) generate(w io.Writer, r io.Reader) (err error) {
	indent := 1
	isGo := false
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
			if isGo {
				c.AppendCode = append(c.AppendCode, string(z.Text()))
				continue
			}
			tab := strings.Repeat("\t", indent)
			t := strings.TrimSpace(string(z.Text()))
			if len(t) > 0 {
				fmt.Fprintf(w, "\n%svecty.Text(%q),", tab, t)
			}
		case html.StartTagToken:
			tab := strings.Repeat("\t", indent)
			indent++
			attrSlice := parseAttrs(z)
			tag := c.tag(z)
			e, ok := elemNameMap[tag]
			if !ok {
				e = fmt.Sprintf("vecty.Tag(%q, ", tag)
			}
			log.Println(tag, attrSlice)
			if tag == "script" {
				for _, attr := range attrSlice {
					if attr.k == "type" && attr.v == "application/x-go" {
						isGo = true
					}
				}
				if isGo {
					continue
				}
			}
			a, t := c.attrs(attrSlice, indent), string(z.Text())
			if indent > 2 {
				fmt.Fprint(w, "\n")
			}
			fmt.Fprintf(w, "%s%s(%s", tab, e, a)
			if len(t) > 0 {
				fmt.Fprintf(w, "\n%svecty.Text(%q),", tab, t)
			}
		case html.SelfClosingTagToken:
			tab := strings.Repeat("\t", indent)
			indent++
			tag := c.tag(z)
			a := c.attrs(parseAttrs(z), indent)
			e, ok := elemNameMap[tag]
			if !ok {
				e = fmt.Sprintf("\n%svecty.Tag(%q, ", tab, tag)
			} else {
				e = e + "("
				c.ExtModules["github.com/gopherjs/vecty/elem"] = true
			}
			if len(a) > 0 {
				fmt.Fprintf(w, "\n%[1]s%s%s\n%[1]s),", tab, e, a)
			} else {
				fmt.Fprintf(w, "\n%[1]s%s),", tab, e)
			}
			indent--
		case html.EndTagToken:
			indent--
			if isGo {
				isGo = false
				continue
			}
			tab := strings.Repeat("\t", indent)
			fmt.Fprintf(w, "\n%s)", tab)
			if indent > 1 {
				fmt.Fprint(w, ",")
			}
		}
	}
	return nil
}

// Do ...
func (c *Converter) Do(output io.Writer, input io.Reader, pkg string) error {
	if err := c.generate(output, input); err != nil && err != io.EOF {
		return err
	}
	return nil
}
