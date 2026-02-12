package velty

import "sort"

// ApplyPatches applies a set of textual replacements to src according to
// Patch spans and returns a new byte slice. Spans are interpreted using the
// same convention as parser spans: [Start, End] (inclusive end offset).
//
// Callers are expected to provide non-overlapping patches relative to the
// original src. Overlapping or out-of-range patches are ignored.
func ApplyPatches(src []byte, patches []Patch) []byte {
	if len(patches) == 0 {
		// preserve src immutability by returning a copy
		dst := make([]byte, len(src))
		copy(dst, src)
		return dst
	}

	// work on a copy to avoid mutating caller slice
	ps := make([]Patch, len(patches))
	copy(ps, patches)

	sort.Slice(ps, func(i, j int) bool { return ps[i].Span.Start < ps[j].Span.Start })

	// conservative capacity: original size; may grow if replacements are larger
	out := make([]byte, 0, len(src))
	last := 0

	for _, p := range ps {
		start := p.Span.Start
		// spans from parser are [start,end] inclusive; convert to [start,endExcl)
		endExcl := p.Span.End + 1

		// basic validation and non-overlap guarantee relative to original src
		if start < last || start < 0 || endExcl < start || endExcl > len(src) {
			// skip invalid/overlapping patch
			continue
		}

		// copy unchanged region
		if start > last {
			out = append(out, src[last:start]...)
		}

		// apply replacement
		if len(p.Replacement) > 0 {
			out = append(out, p.Replacement...)
		}

		last = endExcl
	}

	// append remaining tail
	if last < len(src) {
		out = append(out, src[last:]...)
	}

	return out
}
