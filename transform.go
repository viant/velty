package velty

import (
	"github.com/viant/velty/parser"
)

// TransformTemplate parses a template with spans, runs the supplied adjuster,
// and applies accumulated text patches to return transformed source.
func TransformTemplate(src []byte, adjuster NodeAdjuster, evalCfg ...EvaluateConfig) ([]byte, error) {
	if len(src) == 0 || adjuster == nil {
		dst := make([]byte, len(src))
		copy(dst, src)
		return dst, nil
	}
	root, err := parser.ParseWithSpans(src)
	if err != nil {
		return nil, err
	}
	cfg := EvaluateConfig{}
	if len(evalCfg) > 0 {
		cfg = evalCfg[0]
	}
	// Use adjuster chain so action patches are accumulated into parser context.
	_, ctx, err := applyParserHooksWithConfig("", src, root, nil, NewAdjusterChain(adjuster), cfg)
	if err != nil {
		return nil, err
	}
	return ApplyPatches(src, ctx.Patches()), nil
}
