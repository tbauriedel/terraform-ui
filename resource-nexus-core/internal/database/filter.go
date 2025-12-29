package database

import (
	"fmt"
	"strings"
)

type FilterExpr interface {
	ToSQL(index int) (string, []any, int, error)
}

// Filter represents a simple filter for a database query.
//
// e.g. "name = 'dummy'" -> Filter{Key: "name", Value: "dummy", Operator: "="}.
type Filter struct {
	Key      string
	Operator string
	Value    any
}

// LogicalFilter represents a logical filter for a database query.
// Multiple LogicalFilter can be combined. Add them to the Filters slice.
type LogicalFilter struct {
	Operator string
	Filters  []FilterExpr
}

// ToSQL builds the SQL query string for the filter and returns it along with the arguments.
//
// index is used to generate the argument placeholder '$1'.
//
// // the returned lastIndex is the highest index used for the arguments.
func (f Filter) ToSQL(index int) (string, []any, int, error) {
	if f.Key == "" || f.Operator == "" {
		return "", nil, 0, fmt.Errorf("cant build filter for query")
	}

	return fmt.Sprintf("%s %s $%d", f.Key, f.Operator, index), []any{f.Value}, index + 1, nil
}

// ToSQL builds the SQL query string for the filter and returns it along with the arguments.
//
// index is used to generate the argument placeholder '$1', '$2', ...
// the provided index is used as "start number" for the index.
//
//	e.g. "(name = 'dummy' OR name = 'test')"
//	LogicalFilter{
//	 	Operator: "AND",
//	  	Filters: []FilterExpr{
//	  		Filter{Key: "name", Value: "dummy", Operator: "="},
//	  		Filter{Key: "name", Value: "test", Operator: "="}
//	  	}
//	 }
//
// the returned lastIndex is the highest index used for the arguments.
func (f LogicalFilter) ToSQL(index int) (string, []any, int, error) {
	// no filters added. nothing to combine into logical filter expression
	if len(f.Filters) == 0 {
		return "", nil, 0, nil
	}

	filters := make([]string, 0, len(f.Filters)) // single filter expressions that will be combined

	var arguments []any // arguments for the combined filters

	currentIdx := index // number for argument placeholders. Is incremented for each filter to have '$1', '$2', ...

	// loop over filters and combine them into a single logical expression
	for _, f := range f.Filters {
		// get filter string and argument for filter
		s, a, nextIdx, err := f.ToSQL(currentIdx)
		if err != nil {
			return "", nil, 0, fmt.Errorf("cant build logical filter: %w", err)
		}

		// append filter string and arguments to the combined filters
		filters = append(filters, s)
		arguments = append(arguments, a...)
		currentIdx = nextIdx
	}

	// combine filters by the desired operator
	filter := fmt.Sprintf("(%s)", strings.Join(filters, fmt.Sprintf(" %s ", f.Operator)))

	return filter, arguments, currentIdx, nil
}

// BuildWhere builds the WHERE clause for a database query from a FilterExpr.
//
// When no filter is given, an empty string is returned.
func BuildWhere(filter FilterExpr) (string, []any, error) {
	// no filter given. return empty string
	if filter == nil {
		return "", nil, nil
	}

	// get filter expression as string and arguments
	filterString, arguments, _, err := filter.ToSQL(1)
	if err != nil {
		return "", nil, fmt.Errorf("cant build where: %w", err)
	}

	// build where clause and return clause and arguments
	return fmt.Sprintf(" WHERE %s", filterString), arguments, nil
}
