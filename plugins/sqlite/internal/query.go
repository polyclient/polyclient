// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package internal

import (
	commontypes "github.com/polyclient/polyclient/bindings/polyclient/sql/common-types"
	"github.com/polyclient/polyclient/bindings/polyclient/sql/query"
	"go.bytecodealliance.org/cm"
)

func NewExecute(connID uint32, queryText string, params cm.List[commontypes.SQLValue]) cm.Result[query.ResultSetShape, query.ResultSet, commontypes.Error] {
	// Create column info
	columns := []query.ColumnInfo{
		{
			Name:            "name",
			TypeName:        "text",
			IsPrimaryKey:    false,
			IsAutoIncrement: false,
			IsUnique:        false,
			IsNullable:      false,
			DefaultValue:    cm.Option[commontypes.SQLValue]{},
			Precision:       cm.Option[uint32]{},
			Scale:           cm.Option[float32]{},
			Comment:         cm.Option[string]{},
			Length:          cm.Option[uint32]{},
			Collation:       cm.Option[string]{},
			OrdinalPosition: cm.Option[uint32]{},
		},
		{
			Name:            "email",
			TypeName:        "text",
			IsPrimaryKey:    false,
			IsAutoIncrement: false,
			IsUnique:        true,
			IsNullable:      false,
			DefaultValue:    cm.Option[commontypes.SQLValue]{},
			Precision:       cm.Option[uint32]{},
			Scale:           cm.Option[float32]{},
			Comment:         cm.Option[string]{},
			Length:          cm.Option[uint32]{},
			Collation:       cm.Option[string]{},
			OrdinalPosition: cm.Option[uint32]{},
		},
		{
			Name:            "birthday",
			TypeName:        "timestamp",
			IsPrimaryKey:    false,
			IsAutoIncrement: false,
			IsUnique:        false,
			IsNullable:      true,
			DefaultValue:    cm.Option[commontypes.SQLValue]{},
			Precision:       cm.Option[uint32]{},
			Scale:           cm.Option[float32]{},
			Comment:         cm.Option[string]{},
			Length:          cm.Option[uint32]{},
			Collation:       cm.Option[string]{},
			OrdinalPosition: cm.Option[uint32]{},
		},
	}

	// Create rows
	listRows := cm.ToList([]query.Row{})

	// Create metadata - Fix 1: Convert slice to cm.List correctly
	metadata := query.ResultMetadata{
		Columns:      cm.ToList(columns),
		AffectedRows: cm.Some(uint64(len(listRows.Slice()))),
	}

	// Create result set
	resultSet := query.ResultSet{
		Metadata: metadata,
		Rows:     listRows,
	}

	return cm.OK[cm.Result[query.ResultSetShape, query.ResultSet, commontypes.Error]](resultSet)
}
