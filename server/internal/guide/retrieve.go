package guide

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// Chunk is one searchable piece of a published document.
type Chunk struct {
	DocumentID string
	Title      string
	Category   string
	Seq        int
	Heading    string
	Body       string
}

const chunkTarget = 600 // runes per chunk, roughly one or two paragraphs of Chinese

// SplitChunks cuts a document body into paragraph-aligned chunks, carrying the latest
// markdown heading along so each chunk keeps its section context.
func SplitChunks(body string) []Chunk {
	var (
		out     []Chunk
		heading string
		buf     []string
		size    int
	)
	flush := func() {
		if len(buf) == 0 {
			return
		}
		out = append(out, Chunk{Seq: len(out), Heading: heading, Body: strings.Join(buf, "\n\n")})
		buf, size = nil, 0
	}
	for _, para := range strings.Split(body, "\n\n") {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		if strings.HasPrefix(para, "#") && !strings.Contains(para, "\n") {
			flush()
			heading = strings.TrimSpace(strings.TrimLeft(para, "#"))
			continue
		}
		n := len([]rune(para))
		if size > 0 && size+n > chunkTarget {
			flush()
		}
		// A single paragraph longer than two chunks is cut at sentence ends.
		for n > 2*chunkTarget {
			head, rest := cutSentence(para, chunkTarget)
			buf = append(buf, head)
			flush()
			para = rest
			n = len([]rune(para))
		}
		buf = append(buf, para)
		size += n
	}
	flush()
	return out
}

func cutSentence(s string, near int) (string, string) {
	r := []rune(s)
	for i := near; i < len(r) && i < near+200; i++ {
		if strings.ContainsRune("。！？!?.\n", r[i]) {
			return string(r[:i+1]), strings.TrimSpace(string(r[i+1:]))
		}
	}
	return string(r[:near]), string(r[near:])
}

// terms splits text into search terms: CJK character bigrams (and single characters),
// plus lower-cased latin words and numbers.
func terms(s string) []string {
	var (
		out  []string
		word []rune
		prev rune
	)
	endWord := func() {
		if len(word) > 1 {
			out = append(out, strings.ToLower(string(word)))
		}
		word = word[:0]
	}
	for _, r := range s {
		switch {
		case unicode.Is(unicode.Han, r):
			endWord()
			out = append(out, string(r))
			if prev != 0 {
				out = append(out, string([]rune{prev, r}))
			}
			prev = r
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			word = append(word, r)
			prev = 0
		default:
			endWord()
			prev = 0
		}
	}
	endWord()
	return out
}

// Rank returns up to k chunks most relevant to query, scored by IDF-weighted term overlap.
// Title and heading count as part of the chunk so section names match.
func Rank(chunks []Chunk, query string, k int) []Chunk {
	q := map[string]bool{}
	for _, t := range terms(query) {
		q[t] = true
	}
	if len(q) == 0 || len(chunks) == 0 {
		return nil
	}
	docTerms := make([]map[string]int, len(chunks))
	df := map[string]int{}
	for i, c := range chunks {
		m := map[string]int{}
		for _, t := range terms(c.Title + "\n" + c.Heading + "\n" + c.Body) {
			if q[t] {
				m[t]++
			}
		}
		for t := range m {
			df[t]++
		}
		docTerms[i] = m
	}
	type scored struct {
		i     int
		score float64
	}
	var ranked []scored
	n := float64(len(chunks))
	for i, m := range docTerms {
		var s float64
		for t, tf := range m {
			idf := math.Log(1 + n/float64(df[t]))
			w := idf * (1 + math.Log(float64(tf)))
			if len([]rune(t)) > 1 {
				w *= 2 // bigrams and words are far more specific than single characters
			}
			s += w
		}
		if s > 0 {
			ranked = append(ranked, scored{i, s})
		}
	}
	sort.SliceStable(ranked, func(a, b int) bool { return ranked[a].score > ranked[b].score })
	var out []Chunk
	for _, r := range ranked {
		if len(out) == k {
			break
		}
		out = append(out, chunks[r.i])
	}
	return out
}
