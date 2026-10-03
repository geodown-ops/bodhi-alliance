package guide

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
)

var errUnsupported = errors.New("只接受 .txt、.md、.docx、.pdf 檔")

// ExtractText turns an uploaded file into plain text for the knowledge base.
func ExtractText(filename string, data []byte) (text string, err error) {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".txt", ".md", ".markdown":
		data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
		if !utf8.Valid(data) {
			return "", errors.New("文字檔必須是 UTF-8 編碼")
		}
		text = string(data)
	case ".docx":
		text, err = extractDocx(data)
	case ".pdf":
		text, err = extractPDF(data)
	default:
		return "", errUnsupported
	}
	if err != nil {
		return "", err
	}
	text = normalize(text)
	if strings.TrimSpace(text) == "" {
		return "", errors.New("檔案裡讀不到文字（掃描檔請先轉成文字）")
	}
	return text, nil
}

func extractDocx(data []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", errors.New("無法開啟 Word 檔")
	}
	for _, f := range zr.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		return docxBodyText(io.LimitReader(rc, 50<<20))
	}
	return "", errors.New("Word 檔裡找不到內文")
}

// docxBodyText keeps paragraph breaks and turns Heading styles into markdown headings,
// so chunking can follow the document's own sections.
func docxBodyText(r io.Reader) (string, error) {
	dec := xml.NewDecoder(r)
	var (
		out     strings.Builder
		para    strings.Builder
		heading int
		inText  bool
	)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", errors.New("Word 檔格式損壞")
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "p":
				para.Reset()
				heading = 0
			case "pStyle":
				for _, a := range t.Attr {
					if a.Name.Local == "val" {
						var n int
						if _, err := fmt.Sscanf(strings.ToLower(a.Value), "heading%d", &n); err == nil && n >= 1 && n <= 6 {
							heading = n
						}
					}
				}
			case "t":
				inText = true
			case "tab":
				para.WriteString("\t")
			case "br":
				para.WriteString("\n")
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "t":
				inText = false
			case "p":
				line := strings.TrimSpace(para.String())
				if line != "" && heading > 0 {
					line = strings.Repeat("#", heading) + " " + line
				}
				out.WriteString(line)
				out.WriteString("\n\n")
			}
		case xml.CharData:
			if inText {
				para.Write(t)
			}
		}
	}
	return out.String(), nil
}

func extractPDF(data []byte) (text string, err error) {
	// The PDF library panics on some malformed files.
	defer func() {
		if recover() != nil {
			text, err = "", errors.New("無法讀取這個 PDF")
		}
	}()
	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", errors.New("無法讀取這個 PDF")
	}
	var b strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		s, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		b.WriteString(s)
		b.WriteString("\n\n")
	}
	return b.String(), nil
}

// normalize unifies line endings and collapses runs of blank lines.
func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	var out []string
	blank := 0
	for _, l := range lines {
		l = strings.TrimRight(l, " \t　")
		if l == "" {
			blank++
			if blank > 1 {
				continue
			}
		} else {
			blank = 0
		}
		out = append(out, l)
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}
