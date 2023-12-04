// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

// InnerJoin sets inner join
func (b *Builder) InnerJoin(joinTable, joinCond interface{}, alias string) *Builder {
	return b.Join("INNER", joinTable, joinCond, alias)
}

// LeftJoin sets left join SQL
func (b *Builder) LeftJoin(joinTable, joinCond interface{}, alias string) *Builder {
	return b.Join("LEFT", joinTable, joinCond, alias)
}

// RightJoin sets right join SQL
func (b *Builder) RightJoin(joinTable, joinCond interface{}, alias string) *Builder {
	return b.Join("RIGHT", joinTable, joinCond, alias)
}

// CrossJoin sets cross join SQL
func (b *Builder) CrossJoin(joinTable, joinCond interface{}, alias string) *Builder {
	return b.Join("CROSS", joinTable, joinCond, alias)
}

// FullJoin sets full join SQL
func (b *Builder) FullJoin(joinTable, joinCond interface{}, alias string) *Builder {
	return b.Join("FULL", joinTable, joinCond, alias)
}

// Join sets join table and conditions
func (b *Builder) Join(joinType string, joinTable, joinCond interface{}, alias string) *Builder {
	switch joinCond.(type) {
	case Cond:
		b.joins = append(b.joins, join{joinType, joinTable, joinCond.(Cond), alias})
	case string:
		b.joins = append(b.joins, join{joinType, joinTable, Expr(joinCond.(string)), alias})
	}

	return b
}
