package sqlbuilder

// InsertStatement represents a INSERT statement.
type InsertStatement struct {
	columns ColumnList
	values  []literal
	into    Table
}

// Insert returns new INSERT statement. The table is Table object for into.
func Insert(into Table) *InsertStatement {
	return &InsertStatement{
		into:    into,
		columns: make(ColumnList, 0),
		values:  make([]literal, 0),
	}
}

// Columns sets columns for insert.  This overwrite old results of Columns() or Set().
// If not set this, get error on ToSql().
func (b *InsertStatement) Columns(columns ...Column) *InsertStatement {
	b.columns = ColumnList(columns)
	return b
}

// Values sets VALUES clause. This overwrite old results of Values() or Set().
func (b *InsertStatement) Values(values ...interface{}) *InsertStatement {
	sl := make([]literal, len(values))
	for i := range values {
		sl[i] = toLiteral(values[i])
	}
	b.values = sl
	return b
}

// Set sets the column and value togeter.
// Set cannot be called with Columns() or Values() in a statement.
func (b *InsertStatement) Set(column Column, value interface{}) *InsertStatement {
	b.columns = append(b.columns, column)
	b.values = append(b.values, toLiteral(value))
	return b
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *InsertStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder()
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()

	// INSERT
	bldr.Append("INSERT")

	// INTO Table
	bldr.Append(" INTO ")
	if b.into != nil {
		bldr.AppendItem(b.into)
	} else {
		bldr.SetError(newError("into is nil"))
		return
	}

	// (COLUMN)
	if len(b.columns) == 0 {
		b.columns = b.into.Columns()
	}
	bldr.Append(" ( ")
	bldr.AppendItem(b.columns)
	bldr.Append(" )")

	// VALUES
	if len(b.columns) != len(b.values) {
		bldr.SetError(newError("%d values needed, but got %d", len(b.columns), len(b.values)))
		return
	}
	for i := range b.columns {
		if !b.columns[i].acceptType(b.values[i]) {
			bldr.SetError(newError("%s column not accept %T",
				b.columns[i].config().Type().String(),
				b.values[i].Raw()))
			return
		}
	}
	bldr.Append(" VALUES ( ")
	values := make([]serializable, len(b.values))
	for i := range values {
		values[i] = b.values[i]
	}
	bldr.AppendItems(values, ", ")
	bldr.Append(" )")

	return
}
