package restache

import (
	"golang.org/x/net/html/atom"
)

var elementAtoms = map[atom.Atom]bool{
	atom.A:          true,
	atom.Abbr:       true,
	atom.Address:    true,
	atom.Area:       true,
	atom.Article:    true,
	atom.Aside:      true,
	atom.Audio:      true,
	atom.B:          true,
	atom.Base:       true,
	atom.Bdi:        true,
	atom.Bdo:        true,
	atom.Blockquote: true,
	atom.Body:       true,
	atom.Br:         true,
	atom.Button:     true,
	atom.Canvas:     true,
	atom.Caption:    true,
	atom.Cite:       true,
	atom.Code:       true,
	atom.Col:        true,
	atom.Colgroup:   true,
	atom.Command:    true,
	atom.Data:       true,
	atom.Datalist:   true,
	atom.Dd:         true,
	atom.Del:        true,
	atom.Details:    true,
	atom.Dfn:        true,
	atom.Dialog:     true,
	atom.Div:        true,
	atom.Dl:         true,
	atom.Dt:         true,
	atom.Em:         true,
	atom.Embed:      true,
	atom.Fieldset:   true,
	atom.Figcaption: true,
	atom.Figure:     true,
	atom.Footer:     true,
	atom.Form:       true,
	atom.H1:         true,
	atom.H2:         true,
	atom.H3:         true,
	atom.H4:         true,
	atom.H5:         true,
	atom.H6:         true,
	atom.Head:       true,
	atom.Header:     true,
	atom.Hgroup:     true,
	atom.Hr:         true,
	atom.Html:       true,
	atom.I:          true,
	atom.Iframe:     true,
	atom.Img:        true,
	atom.Input:      true,
	atom.Ins:        true,
	atom.Kbd:        true,
	atom.Keygen:     true,
	atom.Label:      true,
	atom.Legend:     true,
	atom.Li:         true,
	atom.Link:       true,
	atom.Main:       true,
	atom.Map:        true,
	atom.Mark:       true,
	atom.Menu:       true,
	atom.Menuitem:   true,
	atom.Meta:       true,
	atom.Meter:      true,
	atom.Nav:        true,
	atom.Noscript:   true,
	atom.Object:     true,
	atom.Ol:         true,
	atom.Optgroup:   true,
	atom.Option:     true,
	atom.Output:     true,
	atom.P:          true,
	atom.Param:      true,
	atom.Picture:    true,
	atom.Pre:        true,
	atom.Progress:   true,
	atom.Q:          true,
	atom.Rp:         true,
	atom.Rt:         true,
	atom.Ruby:       true,
	atom.S:          true,
	atom.Samp:       true,
	atom.Script:     true,
	atom.Section:    true,
	atom.Select:     true,
	atom.Slot:       true,
	atom.Small:      true,
	atom.Source:     true,
	atom.Span:       true,
	atom.Strong:     true,
	atom.Style:      true,
	atom.Sub:        true,
	atom.Summary:    true,
	atom.Sup:        true,
	atom.Table:      true,
	atom.Tbody:      true,
	atom.Td:         true,
	atom.Template:   true,
	atom.Textarea:   true,
	atom.Tfoot:      true,
	atom.Th:         true,
	atom.Thead:      true,
	atom.Time:       true,
	atom.Title:      true,
	atom.Tr:         true,
	atom.Track:      true,
	atom.U:          true,
	atom.Ul:         true,
	atom.Var:        true,
	atom.Video:      true,
	atom.Wbr:        true,
}
