package rulepathparamstyle

import (
	"fmt"
	"regexp"
	"strconv"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/simplego/text/stringcase"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RulePathParamStyle struct {
	name       string
	stringCase string
	severity   string
}

func NewRule(sev, requiredStringCase string) (RulePathParamStyle, error) {
	canonicalCase, err := stringcase.Parse(requiredStringCase)
	if err != nil {
		return RulePathParamStyle{},
			fmt.Errorf("invalid string case [%s]", requiredStringCase)
	}
	rule := RulePathParamStyle{
		stringCase: canonicalCase,
		severity:   sev}
	switch canonicalCase {
	case stringcase.CamelCase:
		rule.name = lintutil.RulenamePathParamStyleCamelCase
	case stringcase.KebabCase:
		rule.name = lintutil.RulenamePathParamStyleKebabCase
	case stringcase.PascalCase:
		rule.name = lintutil.RulenamePathParamStylePascalCase
	case stringcase.SnakeCase:
		rule.name = lintutil.RulenamePathParamStyleSnakeCase
	default:
		return rule, fmt.Errorf("invalid string case [%s]", canonicalCase)
	}
	return rule, nil
}

func (rule RulePathParamStyle) Name() string {
	return rule.name
}

func (rule RulePathParamStyle) Scope() string {
	return lintutil.ScopeOperation
}

func (rule RulePathParamStyle) Severity() string {
	return rule.severity
}

func (rule RulePathParamStyle) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return nil
}

var rxParams = regexp.MustCompile(`{([^/]*)}`)

func (rule RulePathParamStyle) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	vios := []lintutil.PolicyViolation{}

	for pathUrl := range spec.Paths {
		m := rxParams.FindAllStringSubmatch(pathUrl, -1)
		if len(m) == 0 {
			continue
		}
		jsPtr := jsonutil.PointerSubEscapeAll("%s#/paths/%s",
			pointerBase, pathUrl)
		for _, mi := range m {
			isWantCase, err := stringcase.IsCase(rule.stringCase, mi[1])
			if err != nil {
				vios = append(vios, lintutil.PolicyViolation{
					RuleName: rule.Name(),
					Location: jsPtr,
					Value:    err.Error()})
			} else if !isWantCase {
				vios = append(vios, lintutil.PolicyViolation{
					RuleName: rule.Name(),
					Location: jsPtr})
			}
		}
	}

	openapi3.VisitOperations(spec, func(path, method string, op *oas3.Operation) {
		if op == nil {
			return
		}
		jsPtr := jsonutil.PointerSubEscapeAll(
			"%s#/paths/%s/%s/parameters/",
			pointerBase,
			path,
			method)
		for i, paramRef := range op.Parameters {
			if paramRef == nil || paramRef.Value == nil {
				continue
			}
			isWantCase, err := stringcase.IsCase(rule.stringCase, paramRef.Value.Name)
			if err != nil {
				// error should not occur as stringCase should be pre-validated.
				vios = append(vios, lintutil.PolicyViolation{
					RuleName: rule.Name(),
					Location: jsPtr + strconv.Itoa(i),
					Value:    paramRef.Value.Name + " " + err.Error()})
			} else if !isWantCase {
				vios = append(vios, lintutil.PolicyViolation{
					RuleName: rule.Name(),
					Location: jsPtr + strconv.Itoa(i),
					Value:    paramRef.Value.Name})
			}
		}
	})

	return vios
}
