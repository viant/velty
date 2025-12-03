package velty

// ForEachInfo holds loop context similar to Apache Velocity $foreach.
// Fields are accessible as: $foreach.Index, $foreach.Count, $foreach.HasNext, $foreach.First, $foreach.Last
type ForEachInfo struct {
	Index   int  // zero-based index
	Count   int  // one-based counter
	HasNext bool // true if there is a next element
	First   bool // true if this is the first element
	Last    bool // true if this is the last element
}
