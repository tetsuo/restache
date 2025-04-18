package restache

import (
	"golang.org/x/net/html/atom"
)

// voidElements only have a start tag; ends tags are not specified.
var voidElements = map[atom.Atom]bool{
	atom.Area: true, atom.Br: true, atom.Embed: true, atom.Img: true,
	atom.Input: true, atom.Wbr: true, atom.Col: true, atom.Hr: true,
	atom.Link: true, atom.Track: true, atom.Source: true,
}

var commonElements = map[atom.Atom]bool{
	// HTML
	atom.A: true, atom.Abbr: true, atom.Address: true, atom.Area: true,
	atom.Article: true, atom.Aside: true, atom.Audio: true, atom.B: true,
	atom.Base: true, atom.Bdi: true, atom.Bdo: true, atom.Blockquote: true,
	atom.Body: true, atom.Br: true, atom.Button: true, atom.Canvas: true,
	atom.Caption: true, atom.Cite: true, atom.Code: true, atom.Col: true,
	atom.Colgroup: true, atom.Command: true, atom.Data: true, atom.Datalist: true,
	atom.Dd: true, atom.Del: true, atom.Details: true, atom.Dfn: true,
	atom.Dialog: true, atom.Div: true, atom.Dl: true, atom.Dt: true,
	atom.Em: true, atom.Embed: true, atom.Fieldset: true, atom.Figcaption: true,
	atom.Figure: true, atom.Footer: true, atom.Form: true, atom.H1: true,
	atom.H2: true, atom.H3: true, atom.H4: true, atom.H5: true,
	atom.H6: true, atom.Head: true, atom.Header: true, atom.Hgroup: true,
	atom.Hr: true, atom.Html: true, atom.I: true, atom.Iframe: true,
	atom.Img: true, atom.Input: true, atom.Ins: true, atom.Kbd: true,
	atom.Keygen: true, atom.Label: true, atom.Legend: true, atom.Li: true,
	atom.Link: true, atom.Main: true, atom.Map: true, atom.Mark: true,
	atom.Menu: true, atom.Menuitem: true, atom.Meta: true, atom.Meter: true,
	atom.Nav: true, atom.Noscript: true, atom.Object: true, atom.Ol: true,
	atom.Optgroup: true, atom.Option: true, atom.Output: true, atom.P: true,
	atom.Param: true, atom.Picture: true, atom.Pre: true, atom.Progress: true,
	atom.Q: true, atom.Rp: true, atom.Rt: true, atom.Ruby: true,
	atom.S: true, atom.Samp: true, atom.Script: true, atom.Section: true,
	atom.Select: true, atom.Slot: true, atom.Small: true, atom.Source: true,
	atom.Span: true, atom.Strong: true, atom.Style: true, atom.Sub: true,
	atom.Summary: true, atom.Sup: true, atom.Table: true, atom.Tbody: true,
	atom.Td: true, atom.Template: true, atom.Textarea: true, atom.Tfoot: true,
	atom.Th: true, atom.Thead: true, atom.Time: true, atom.Title: true,
	atom.Tr: true, atom.Track: true, atom.U: true, atom.Ul: true,
	atom.Var: true, atom.Video: true, atom.Wbr: true, atom.Frame: true,
	atom.Frameset: true, atom.Malignmark: true, atom.Manifest: true, atom.Rb: true,
	atom.Rtc: true,

	// SVG
	atom.Svg: true, atom.Desc: true, atom.Foreignobject: true, atom.Image: true,

	// MathML
	atom.Math: true, atom.Mglyph: true, atom.Mi: true, atom.Mn: true,
	atom.Mo: true, atom.Ms: true, atom.Mtext: true,

	// Legacy
	atom.Acronym:       true, // obsolete (use <abbr>)
	atom.Xmp:           true, // obsolete or non-standard
	atom.Applet:        true, // deprecated
	atom.Annotation:    true, // technically MathML, but sometimes XML
	atom.AnnotationXml: true, // technically MathML, but sometimes XML
	atom.Basefont:      true, // deprecated
	atom.Bgsound:       true, // obsolete or non-standard
	atom.Big:           true, // deprecated
	atom.Blink:         true, // obsolete or non-standard
	atom.Center:        true, // deprecated
	atom.Font:          true, // deprecated
	atom.Isindex:       true, // obsolete or non-standard
	atom.Listing:       true, // obsolete or non-standard
	atom.Marquee:       true, // obsolete or non-standard
	atom.Nobr:          true, // obsolete or non-standard
	atom.Noembed:       true, // obsolete or non-standard
	atom.Noframes:      true, // obsolete or non-standard
	atom.Plaintext:     true, // obsolete or non-standard
	atom.Scoped:        true, // also exists as attribute
	atom.Spacer:        true, // obsolete or non-standard
	atom.Strike:        true, // deprecated
	atom.Tt:            true, // deprecated
}

var camelAttrTags = map[atom.Atom]bool{
	atom.A:        true,
	atom.Area:     true,
	atom.Audio:    true,
	atom.Button:   true,
	atom.Del:      true,
	atom.Form:     true,
	atom.Iframe:   true,
	atom.Img:      true,
	atom.Input:    true,
	atom.Ins:      true,
	atom.Label:    true,
	atom.Link:     true,
	atom.Meta:     true,
	atom.Object:   true,
	atom.Output:   true,
	atom.Script:   true,
	atom.Select:   true,
	atom.Source:   true,
	atom.Td:       true,
	atom.Textarea: true,
	atom.Th:       true,
	atom.Time:     true,
	atom.Track:    true,
	atom.Video:    true,
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

var nonSpecCamelAttrTags = map[atom.Atom]bool{
	atom.Img:      true,
	atom.Link:     true,
	atom.Script:   true,
	atom.Select:   true,
	atom.Frame:    true,
	atom.Input:    true,
	atom.Object:   true,
	atom.Table:    true,
	atom.Textarea: true,
	atom.Video:    true,
	atom.Audio:    true,
	atom.Button:   true,
	atom.Form:     true,
	atom.Iframe:   true,
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
}
