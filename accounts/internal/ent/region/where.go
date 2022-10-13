// Code generated by ent, DO NOT EDIT.

package region

import (
	"entgo.io/ent/dialect/sql"
	"github.com/transerver/accounts/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		v := make([]any, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		v := make([]any, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Code applies equality check predicate on the "code" field. It's identical to CodeEQ.
func Code(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCode), v))
	})
}

// Area applies equality check predicate on the "area" field. It's identical to AreaEQ.
func Area(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldArea), v))
	})
}

// Img applies equality check predicate on the "img" field. It's identical to ImgEQ.
func Img(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldImg), v))
	})
}

// CodeEQ applies the EQ predicate on the "code" field.
func CodeEQ(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCode), v))
	})
}

// CodeNEQ applies the NEQ predicate on the "code" field.
func CodeNEQ(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCode), v))
	})
}

// CodeIn applies the In predicate on the "code" field.
func CodeIn(vs ...string) predicate.Region {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldCode), v...))
	})
}

// CodeNotIn applies the NotIn predicate on the "code" field.
func CodeNotIn(vs ...string) predicate.Region {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldCode), v...))
	})
}

// CodeGT applies the GT predicate on the "code" field.
func CodeGT(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCode), v))
	})
}

// CodeGTE applies the GTE predicate on the "code" field.
func CodeGTE(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCode), v))
	})
}

// CodeLT applies the LT predicate on the "code" field.
func CodeLT(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCode), v))
	})
}

// CodeLTE applies the LTE predicate on the "code" field.
func CodeLTE(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCode), v))
	})
}

// CodeContains applies the Contains predicate on the "code" field.
func CodeContains(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCode), v))
	})
}

// CodeHasPrefix applies the HasPrefix predicate on the "code" field.
func CodeHasPrefix(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCode), v))
	})
}

// CodeHasSuffix applies the HasSuffix predicate on the "code" field.
func CodeHasSuffix(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCode), v))
	})
}

// CodeEqualFold applies the EqualFold predicate on the "code" field.
func CodeEqualFold(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCode), v))
	})
}

// CodeContainsFold applies the ContainsFold predicate on the "code" field.
func CodeContainsFold(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCode), v))
	})
}

// AreaEQ applies the EQ predicate on the "area" field.
func AreaEQ(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldArea), v))
	})
}

// AreaNEQ applies the NEQ predicate on the "area" field.
func AreaNEQ(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldArea), v))
	})
}

// AreaIn applies the In predicate on the "area" field.
func AreaIn(vs ...string) predicate.Region {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldArea), v...))
	})
}

// AreaNotIn applies the NotIn predicate on the "area" field.
func AreaNotIn(vs ...string) predicate.Region {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldArea), v...))
	})
}

// AreaGT applies the GT predicate on the "area" field.
func AreaGT(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldArea), v))
	})
}

// AreaGTE applies the GTE predicate on the "area" field.
func AreaGTE(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldArea), v))
	})
}

// AreaLT applies the LT predicate on the "area" field.
func AreaLT(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldArea), v))
	})
}

// AreaLTE applies the LTE predicate on the "area" field.
func AreaLTE(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldArea), v))
	})
}

// AreaContains applies the Contains predicate on the "area" field.
func AreaContains(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldArea), v))
	})
}

// AreaHasPrefix applies the HasPrefix predicate on the "area" field.
func AreaHasPrefix(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldArea), v))
	})
}

// AreaHasSuffix applies the HasSuffix predicate on the "area" field.
func AreaHasSuffix(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldArea), v))
	})
}

// AreaEqualFold applies the EqualFold predicate on the "area" field.
func AreaEqualFold(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldArea), v))
	})
}

// AreaContainsFold applies the ContainsFold predicate on the "area" field.
func AreaContainsFold(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldArea), v))
	})
}

// ImgEQ applies the EQ predicate on the "img" field.
func ImgEQ(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldImg), v))
	})
}

// ImgNEQ applies the NEQ predicate on the "img" field.
func ImgNEQ(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldImg), v))
	})
}

// ImgIn applies the In predicate on the "img" field.
func ImgIn(vs ...string) predicate.Region {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldImg), v...))
	})
}

// ImgNotIn applies the NotIn predicate on the "img" field.
func ImgNotIn(vs ...string) predicate.Region {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldImg), v...))
	})
}

// ImgGT applies the GT predicate on the "img" field.
func ImgGT(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldImg), v))
	})
}

// ImgGTE applies the GTE predicate on the "img" field.
func ImgGTE(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldImg), v))
	})
}

// ImgLT applies the LT predicate on the "img" field.
func ImgLT(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldImg), v))
	})
}

// ImgLTE applies the LTE predicate on the "img" field.
func ImgLTE(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldImg), v))
	})
}

// ImgContains applies the Contains predicate on the "img" field.
func ImgContains(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldImg), v))
	})
}

// ImgHasPrefix applies the HasPrefix predicate on the "img" field.
func ImgHasPrefix(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldImg), v))
	})
}

// ImgHasSuffix applies the HasSuffix predicate on the "img" field.
func ImgHasSuffix(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldImg), v))
	})
}

// ImgEqualFold applies the EqualFold predicate on the "img" field.
func ImgEqualFold(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldImg), v))
	})
}

// ImgContainsFold applies the ContainsFold predicate on the "img" field.
func ImgContainsFold(v string) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldImg), v))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Region) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Region) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Region) predicate.Region {
	return predicate.Region(func(s *sql.Selector) {
		p(s.Not())
	})
}
