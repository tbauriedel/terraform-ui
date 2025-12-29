package database

import (
	"reflect"
	"testing"
)

func TestFilterToSQL(t *testing.T) {
	filter := Filter{
		Key:      "dummy",
		Operator: "=",
		Value:    "foobar",
	}

	f, args, _, err := filter.ToSQL(1)
	if err != nil {
		t.Fatal(err)
	}

	if f != "dummy = $1" {
		t.Fatal("wrong sql generated")
	}

	if args[0] != "foobar" {
		t.Fatal("wrong args returned")
	}
}

func TestLogicalFilterToSQL(t *testing.T) {
	filter := LogicalFilter{
		Operator: "OR",
		Filters: []FilterExpr{
			&Filter{
				Key:      "car",
				Operator: "=",
				Value:    "volvo",
			},
			&Filter{
				Key:      "car",
				Operator: "=",
				Value:    "bmw",
			},
		},
	}

	f, args, _, err := filter.ToSQL(1)
	if err != nil {
		t.Fatal(err)
	}

	if f != "(car = $1 OR car = $2)" {
		t.Fatalf("\nactual: %s\nexpected: (car = ? OR car = ?)", f)
	}

	if !reflect.DeepEqual(args, []any{"volvo", "bmw"}) {
		t.Fatal("wrong args returned")
	}
}

func TestBuildWhereFilter(t *testing.T) {
	filter := Filter{
		Key:      "car",
		Operator: "=",
		Value:    "volvo",
	}

	s, args, err := BuildWhere(filter)
	if err != nil {
		t.Fatal(err)
	}

	if s != " WHERE car = $1" {
		t.Fatalf("\nactual: %s\nexpected: WHERE car = $1", s)
	}

	if !reflect.DeepEqual(args, []any{"volvo"}) {
		t.Fatalf("wrong args returned: %v", args)
	}
}

func TestBuildWhereLogicalFilter(t *testing.T) {
	filter := LogicalFilter{
		Operator: "AND",
		Filters: []FilterExpr{
			LogicalFilter{
				Operator: "OR",
				Filters: []FilterExpr{
					Filter{
						Key:      "car",
						Operator: "=",
						Value:    "volvo",
					},
					Filter{
						Key:      "car",
						Operator: "=",
						Value:    "bmw",
					},
				},
			},
			Filter{
				Key:      "car",
				Operator: "=",
				Value:    "audi",
			},
		},
	}

	s, args, err := BuildWhere(filter)
	if err != nil {
		t.Fatal(err)
	}

	expected := " WHERE ((car = $1 OR car = $2) AND car = $3)"
	if s != expected {
		t.Fatalf("\nactual: %s\nexpected: %s", s, expected)
	}

	if !reflect.DeepEqual(args, []any{"volvo", "bmw", "audi"}) {
		t.Fatalf("wrong args returned: %v", args)
	}
}
