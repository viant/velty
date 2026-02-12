package velty

// Severity level for diagnostics.
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
)

// ErrorRecord captures a diagnostic with optional span/position.
type ErrorRecord struct {
	Severity Severity
	Code     string
	Message  string
	File     string
	Span     Span
	Position Position
	Policy   string
	// Optional node/context info
	NodeKind string
	NodeID   int
	// SelectorChain may carry $foo.bar.baz chain segments if available
	SelectorChain []string
}

// Reporter accumulates diagnostics and can flush them to a sink.
type Reporter struct {
	records []ErrorRecord
}

func NewReporter() *Reporter { return &Reporter{records: make([]ErrorRecord, 0, 16)} }

// Report appends a diagnostic record.
func (r *Reporter) Report(rec ErrorRecord) { r.records = append(r.records, rec) }

// Records returns current diagnostics.
func (r *Reporter) Records() []ErrorRecord { return r.records }

// Reset clears buffered diagnostics.
func (r *Reporter) Reset() { r.records = r.records[:0] }
