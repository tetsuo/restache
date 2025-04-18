package restache

import (
	"golang.org/x/net/html/atom"
)

// voidElements only have a start tag; ends tags are not specified.
var voidElements = map[atom.Atom]struct{}{
	atom.Area: {}, atom.Br: {}, atom.Embed: {}, atom.Img: {},
	atom.Input: {}, atom.Wbr: {}, atom.Col: {}, atom.Hr: {},
	atom.Link: {}, atom.Track: {}, atom.Source: {},
}

var commonElements = map[atom.Atom]struct{}{
	// HTML
	atom.A: {}, atom.Abbr: {}, atom.Address: {}, atom.Area: {},
	atom.Article: {}, atom.Aside: {}, atom.Audio: {}, atom.B: {},
	atom.Base: {}, atom.Bdi: {}, atom.Bdo: {}, atom.Blockquote: {},
	atom.Body: {}, atom.Br: {}, atom.Button: {}, atom.Canvas: {},
	atom.Caption: {}, atom.Cite: {}, atom.Code: {}, atom.Col: {},
	atom.Colgroup: {}, atom.Command: {}, atom.Data: {}, atom.Datalist: {},
	atom.Dd: {}, atom.Del: {}, atom.Details: {}, atom.Dfn: {},
	atom.Dialog: {}, atom.Div: {}, atom.Dl: {}, atom.Dt: {},
	atom.Em: {}, atom.Embed: {}, atom.Fieldset: {}, atom.Figcaption: {},
	atom.Figure: {}, atom.Footer: {}, atom.Form: {}, atom.H1: {},
	atom.H2: {}, atom.H3: {}, atom.H4: {}, atom.H5: {},
	atom.H6: {}, atom.Head: {}, atom.Header: {}, atom.Hgroup: {},
	atom.Hr: {}, atom.Html: {}, atom.I: {}, atom.Iframe: {},
	atom.Img: {}, atom.Input: {}, atom.Ins: {}, atom.Kbd: {},
	atom.Keygen: {}, atom.Label: {}, atom.Legend: {}, atom.Li: {},
	atom.Link: {}, atom.Main: {}, atom.Map: {}, atom.Mark: {},
	atom.Menu: {}, atom.Menuitem: {}, atom.Meta: {}, atom.Meter: {},
	atom.Nav: {}, atom.Noscript: {}, atom.Object: {}, atom.Ol: {},
	atom.Optgroup: {}, atom.Option: {}, atom.Output: {}, atom.P: {},
	atom.Param: {}, atom.Picture: {}, atom.Pre: {}, atom.Progress: {},
	atom.Q: {}, atom.Rp: {}, atom.Rt: {}, atom.Ruby: {},
	atom.S: {}, atom.Samp: {}, atom.Script: {}, atom.Section: {},
	atom.Select: {}, atom.Slot: {}, atom.Small: {}, atom.Source: {},
	atom.Span: {}, atom.Strong: {}, atom.Style: {}, atom.Sub: {},
	atom.Summary: {}, atom.Sup: {}, atom.Table: {}, atom.Tbody: {},
	atom.Td: {}, atom.Template: {}, atom.Textarea: {}, atom.Tfoot: {},
	atom.Th: {}, atom.Thead: {}, atom.Time: {}, atom.Title: {},
	atom.Tr: {}, atom.Track: {}, atom.U: {}, atom.Ul: {},
	atom.Var: {}, atom.Video: {}, atom.Wbr: {}, atom.Frame: {},
	atom.Frameset: {}, atom.Malignmark: {}, atom.Manifest: {}, atom.Rb: {},
	atom.Rtc: {},

	// SVG
	atom.Svg: {}, atom.Desc: {}, atom.Foreignobject: {}, atom.Image: {},

	// MathML
	atom.Math: {}, atom.Mglyph: {}, atom.Mi: {}, atom.Mn: {},
	atom.Mo: {}, atom.Ms: {}, atom.Mtext: {},

	// Legacy
	atom.Acronym:       {}, // obsolete (use <abbr>)
	atom.Xmp:           {}, // obsolete or non-standard
	atom.Applet:        {}, // deprecated
	atom.Annotation:    {}, // technically MathML, but sometimes XML
	atom.AnnotationXml: {}, // technically MathML, but sometimes XML
	atom.Basefont:      {}, // deprecated
	atom.Bgsound:       {}, // obsolete or non-standard
	atom.Big:           {}, // deprecated
	atom.Blink:         {}, // obsolete or non-standard
	atom.Center:        {}, // deprecated
	atom.Font:          {}, // deprecated
	atom.Isindex:       {}, // obsolete or non-standard
	atom.Listing:       {}, // obsolete or non-standard
	atom.Marquee:       {}, // obsolete or non-standard
	atom.Nobr:          {}, // obsolete or non-standard
	atom.Noembed:       {}, // obsolete or non-standard
	atom.Noframes:      {}, // obsolete or non-standard
	atom.Plaintext:     {}, // obsolete or non-standard
	atom.Scoped:        {}, // also exists as attribute
	atom.Spacer:        {}, // obsolete or non-standard
	atom.Strike:        {}, // deprecated
	atom.Tt:            {}, // deprecated
}

var camelAttrTags = map[atom.Atom]struct{}{
	atom.A:        {},
	atom.Area:     {},
	atom.Audio:    {},
	atom.Button:   {},
	atom.Del:      {},
	atom.Form:     {},
	atom.Iframe:   {},
	atom.Img:      {},
	atom.Input:    {},
	atom.Ins:      {},
	atom.Label:    {},
	atom.Link:     {},
	atom.Meta:     {},
	atom.Object:   {},
	atom.Output:   {},
	atom.Script:   {},
	atom.Select:   {},
	atom.Source:   {},
	atom.Td:       {},
	atom.Textarea: {},
	atom.Th:       {},
	atom.Time:     {},
	atom.Track:    {},
	atom.Video:    {},
}

var camelAttrTable = map[uint64]string{
	0x63c060002780c: "autoComplete",    // select + autocomplete
	0x26e040002780c: "autoComplete",    // form + autocomplete
	0x26e0400001a0e: "acceptCharset",   // form + accept-charset
	0x26e0400028d07: "encType",         // form + enctype
	0x26e040002b20a: "noValidate",      // form + novalidate
	0x204030002b808: "dateTime",        // ins + datetime
	0x92020001c407:  "colSpan",         // td + colspan
	0x920200009c07:  "rowSpan",         // td + rowspan
	0x352080002780c: "autoComplete",    // textarea + autocomplete
	0x3520800033d09: "maxLength",       // textarea + maxlength
	0x3520800034709: "minLength",       // textarea + minlength
	0x3520800035708: "readOnly",        // textarea + readonly
	0x1ba050005f907: "srcLang",         // track + srclang
	0x2f10500013c08: "autoPlay",        // video + autoplay
	0x2f1050001fb0b: "crossOrigin",     // video + crossorigin
	0x2f1050001400b: "playsInline",     // video + playsinline
	0x1910600026e0a: "formAction",      // button + formaction
	0x191060002890b: "formEncType",     // button + formenctype
	0x191060002a40a: "formMethod",      // button + formmethod
	0x191060002ae0e: "formNoValidate",  // button + formnovalidate
	0x191060002c00a: "formTarget",      // button + formtarget
	0x2680600059e06: "useMap",          // object + usemap
	0x218060001fb0b: "crossOrigin",     // script + crossorigin
	0x218060000a208: "noModule",        // script + nomodule
	0x218060003d10e: "referrerPolicy",  // script + referrerpolicy
	0x2180600002107: "charSet",         // script + charset
	0x378060006f906: "srcSet",          // source + srcset
	0x1150500013c08: "autoPlay",        // audio + autoplay
	0x115050001fb0b: "crossOrigin",     // audio + crossorigin
	0x307030001fb0b: "crossOrigin",     // img + crossorigin
	0x307030003d10e: "referrerPolicy",  // img + referrerpolicy
	0x307030006f906: "srcSet",          // img + srcset
	0x3070300059e06: "useMap",          // img + usemap
	0x1560200009c07: "rowSpan",         // th + rowspan
	0x156020001c407: "colSpan",         // th + colspan
	0x42040002b808:  "dateTime",        // time + datetime
	0x100002107:     "charSet",         // a + charset
	0x10002e008:     "hrefLang",        // a + hreflang
	0x10003d10e:     "referrerPolicy",  // a + referrerpolicy
	0x356040003d10e: "referrerPolicy",  // area + referrerpolicy
	0x452030002b808: "dateTime",        // del + datetime
	0x2fc060003d10e: "referrerPolicy",  // iframe + referrerpolicy
	0x2fc060005c006: "srcDoc",          // iframe + srcdoc
	0x2fc060002080f: "allowFullScreen", // iframe + allowfullscreen
	0x44b0500026e0a: "formAction",      // input + formaction
	0x44b050002a40a: "formMethod",      // input + formmethod
	0x44b050002ae0e: "formNoValidate",  // input + formnovalidate
	0x44b0500033d09: "maxLength",       // input + maxlength
	0x44b0500034709: "minLength",       // input + minlength
	0x44b050002890b: "formEncType",     // input + formenctype
	0x44b050002c00a: "formTarget",      // input + formtarget
	0x44b0500035708: "readOnly",        // input + readonly
	0x44b050002780c: "autoComplete",    // input + autocomplete
	0x174040001fb0b: "crossOrigin",     // link + crossorigin
	0x174040002e008: "hrefLang",        // link + hreflang
	0x174040003d10e: "referrerPolicy",  // link + referrerpolicy
	0x4b80400002107: "charSet",         // meta + charset
	0x4b8040002e80a: "httpEquiv",       // meta + http-equiv
}

var nonSpecCamelAttrTags = map[atom.Atom]struct{}{
	atom.Img:      {},
	atom.Link:     {},
	atom.Script:   {},
	atom.Select:   {},
	atom.Frame:    {},
	atom.Input:    {},
	atom.Object:   {},
	atom.Table:    {},
	atom.Textarea: {},
	atom.Video:    {},
	atom.Audio:    {},
	atom.Button:   {},
	atom.Form:     {},
	atom.Iframe:   {},
}

var nonSpecCamelAttrTable = map[uint64]string{
	0x307031bc09f5e: "fetchPriority",           // img + fetchpriority
	0x174041bc09f5e: "fetchPriority",           // link + fetchpriority
	0x17404afadfdeb: "imageSizes",              // link + imagesizes
	0x1740486aaa7bd: "imageSrcSet",             // link + imagesrcset
	0x218061bc09f5e: "fetchPriority",           // script + fetchpriority
	0x63c06a01a4778: "defaultValue",            // select + defaultvalue (React)
	0x8b057d8d7e5c:  "marginWidth",             // frame + marginwidth
	0x8b05dc68add1:  "frameBorder",             // frame + frameborder
	0x8b05f62c8d4b:  "marginHeight",            // frame + marginheight
	0x44b055053ea34: "defaultChecked",          // input + defaultchecked (React)
	0x44b05a01a4778: "defaultValue",            // input + defaultvalue (React)
	0x44b05941b2242: "popoverTarget",           // input + popovertarget
	0x44b058c8f1e30: "popoverTargetAction",     // input + popovertargetaction
	0x44b05495cbf9f: "autoCapitalize",          // input + autocapitalize
	0x44b05f0fc73fc: "autoSave",                // input + autosave (obsolete, relevant in Apple/Safari)
	0x26806dff20beb: "classID",                 // object + classid
	0x595058b87a8ff: "cellPadding",             // table + cellpadding
	0x595058af90373: "cellSpacing",             // table + cellspacing
	0x35208495cbf9f: "autoCapitalize",          // textarea + autocapitalize
	0x35208e8c5b4f5: "autoCorrect",             // textarea + autocorrect
	0x35208a01a4778: "defaultValue",            // textarea + defaultvalue (React)
	0x2f10570805ce0: "controlsList",            // video + controlslist
	0x2f105319cac13: "disablePictureInPicture", // video + disablepictureinpicture
	0x2f1056f2dd125: "disableRemotePlayback",   // video + disableremoteplayback
	0x1150570805ce0: "controlsList",            // audio + controlslist
	0x115056f2dd125: "disableRemotePlayback",   // audio + disableremoteplayback
	0x19106941b2242: "popoverTarget",           // button + popovertarget
	0x191068c8f1e30: "popoverTargetAction",     // button + popovertargetaction
	0x26e04495cbf9f: "autoCapitalize",          // form + autocapitalize
	0x2fc06dc68add1: "frameBorder",             // iframe + frameborder
	0x2fc06f62c8d4b: "marginHeight",            // iframe + marginheight
	0x2fc067d8d7e5c: "marginWidth",             // iframe + marginwidth
}

var globalCamelAttrTable = map[atom.Atom]string{
	atom.Accesskey:       "accessKey",
	atom.Autofocus:       "autoFocus", // <select>, <textarea>, <button>, <input> + contenteditable
	atom.Class:           "className", // applies to *, also a special keyword
	atom.Contenteditable: "contentEditable",
	atom.For:             "htmlFor",   // applies to <label> and <output>, but special keyword
	atom.Inputmode:       "inputMode", // applies to <input> + contenteditable
	atom.Itemid:          "itemID",
	atom.Itemprop:        "itemProp",
	atom.Itemref:         "itemRef",
	atom.Itemscope:       "itemScope",
	atom.Itemtype:        "itemType",
	atom.Spellcheck:      "spellCheck", // applies to <textarea> + contenteditable
	atom.Tabindex:        "tabIndex",   // applies to <input> + relevant to contenteditable
	// Events
	atom.Onabort:                   "onAbort",
	atom.Onafterprint:              "onAfterPrint",
	atom.Onautocomplete:            "onAutoComplete",
	atom.Onautocompleteerror:       "onAutoCompleteError",
	atom.Onauxclick:                "onAuxClick",
	atom.Onbeforeprint:             "onBeforePrint",
	atom.Onbeforeunload:            "onBeforeUnload",
	atom.Onblur:                    "onBlur",
	atom.Oncancel:                  "onCancel",
	atom.Oncanplay:                 "onCanPlay",
	atom.Oncanplaythrough:          "onCanPlayThrough",
	atom.Onchange:                  "onChange",
	atom.Onclick:                   "onClick",
	atom.Onclose:                   "onClose",
	atom.Oncontextmenu:             "onContextMenu",
	atom.Oncopy:                    "onCopy",
	atom.Oncuechange:               "onCueChange",
	atom.Oncut:                     "onCut",
	atom.Ondblclick:                "onDoubleClick",
	atom.Ondrag:                    "onDrag",
	atom.Ondragend:                 "onDragEnd",
	atom.Ondragenter:               "onDragEnter",
	atom.Ondragexit:                "onDragExit",
	atom.Ondragleave:               "onDragLeave",
	atom.Ondragover:                "onDragOver",
	atom.Ondragstart:               "onDragStart",
	atom.Ondrop:                    "onDrop",
	atom.Ondurationchange:          "onDurationChange",
	atom.Onemptied:                 "onEmptied",
	atom.Onended:                   "onEnded",
	atom.Onerror:                   "onError",
	atom.Onfocus:                   "onFocus",
	atom.Onhashchange:              "onHashChange",
	atom.Oninput:                   "onInput",
	atom.Oninvalid:                 "onInvalid",
	atom.Onkeydown:                 "onKeyDown",
	atom.Onkeypress:                "onKeyPress",
	atom.Onkeyup:                   "onKeyUp",
	atom.Onlanguagechange:          "onLanguageChange",
	atom.Onload:                    "onLoad",
	atom.Onloadeddata:              "onLoadedData",
	atom.Onloadedmetadata:          "onLoadedMetadata",
	atom.Onloadend:                 "onLoadEnd",
	atom.Onloadstart:               "onLoadStart",
	atom.Onmessage:                 "onMessage",
	atom.Onmessageerror:            "onMessageError",
	atom.Onmousedown:               "onMouseDown",
	atom.Onmouseenter:              "onMouseEnter",
	atom.Onmouseleave:              "onMouseLeave",
	atom.Onmousemove:               "onMouseMove",
	atom.Onmouseout:                "onMouseOut",
	atom.Onmouseover:               "onMouseOver",
	atom.Onmouseup:                 "onMouseUp",
	atom.Onmousewheel:              "onMouseWheel",
	atom.Onoffline:                 "onOffline",
	atom.Ononline:                  "onOnline",
	atom.Onpagehide:                "onPageHide",
	atom.Onpageshow:                "onPageShow",
	atom.Onpaste:                   "onPaste",
	atom.Onpause:                   "onPause",
	atom.Onplay:                    "onPlay",
	atom.Onplaying:                 "onPlaying",
	atom.Onpopstate:                "onPopState",
	atom.Onprogress:                "onProgress",
	atom.Onratechange:              "onRateChange",
	atom.Onrejectionhandled:        "onRejectionHandled",
	atom.Onreset:                   "onReset",
	atom.Onresize:                  "onResize",
	atom.Onscroll:                  "onScroll",
	atom.Onsecuritypolicyviolation: "onSecurityPolicyViolation",
	atom.Onseeked:                  "onSeeked",
	atom.Onseeking:                 "onSeeking",
	atom.Onselect:                  "onSelect",
	atom.Onshow:                    "onShow",
	atom.Onsort:                    "onSort",
	atom.Onstalled:                 "onStalled",
	atom.Onstorage:                 "onStorage",
	atom.Onsubmit:                  "onSubmit",
	atom.Onsuspend:                 "onSuspend",
	atom.Ontimeupdate:              "onTimeUpdate",
	atom.Ontoggle:                  "onToggle",
	atom.Onunhandledrejection:      "onUnhandledRejection",
	atom.Onunload:                  "onUnload",
	atom.Onvolumechange:            "onVolumeChange",
	atom.Onwaiting:                 "onWaiting",
	atom.Onwheel:                   "onWheel",

// Following attributes always receive "true" regardless of the value set.
// https://html.spec.whatwg.org/#boolean-attributes
var boolAttrs = map[atom.Atom]struct{}{
	atom.Allowfullscreen:     {},
	atom.Allowpaymentrequest: {},
	atom.Async:               {},
	atom.Autofocus:           {},
	atom.Autoplay:            {},
	atom.Checked:             {},
	atom.Controls:            {},
	atom.Default:             {},
	atom.Defer:               {},
	atom.Disabled:            {},
	atom.Formnovalidate:      {},
	atom.Hidden:              {},
	atom.Ismap:               {},
	atom.Itemscope:           {},
	atom.Loop:                {},
	atom.Multiple:            {},
	atom.Muted:               {},
	atom.Nomodule:            {},
	atom.Novalidate:          {},
	atom.Open:                {},
	atom.Playsinline:         {},
	atom.Readonly:            {},
	atom.Required:            {},
	atom.Reversed:            {},
	atom.Selected:            {},
}

