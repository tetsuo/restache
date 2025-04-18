package restache

import (
	"golang.org/x/net/html/atom"
)

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
