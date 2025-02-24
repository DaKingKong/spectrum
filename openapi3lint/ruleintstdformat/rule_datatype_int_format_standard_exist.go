package ruleintstdformat

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleDatatypeIntFormatStandardExist struct {
	name     string
	severity string
}

func NewRule(sev string) RuleDatatypeIntFormatStandardExist {
	return RuleDatatypeIntFormatStandardExist{
		name:     lintutil.RulenameDatatypeIntFormatStandardExist,
		severity: sev}
}

func (rule RuleDatatypeIntFormatStandardExist) Name() string {
	return rule.name
}

func (rule RuleDatatypeIntFormatStandardExist) Severity() string {
	return rule.severity
}

func (rule RuleDatatypeIntFormatStandardExist) Scope() string {
	return lintutil.ScopeSpecification
}

func (rule RuleDatatypeIntFormatStandardExist) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}

func (rule RuleDatatypeIntFormatStandardExist) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	vios := []lintutil.PolicyViolation{}
	openapi3.VisitTypesFormats(
		spec,
		func(jsonPointerRoot, oasType, oasFormat string) {
			if oasType == openapi3.TypeInteger &&
				oasFormat != openapi3.FormatInt32 &&
				oasFormat != openapi3.FormatInt64 {
				vios = append(vios, lintutil.PolicyViolation{
					RuleName: rule.Name(),
					Location: urlutil.JoinAbsolute(pointerBase+jsonPointerRoot, "format"),
					Value:    oasFormat})
			}
		},
	)
	return vios
}
