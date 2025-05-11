package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

func sampleLaplace(b float64) float64 {
	// Generate a uniform random number in [0,1)
	u := rand.Float64()
	// Use the inverse CDF method based on u:
	if u < 0.5 {
		return b * math.Log(2*u)
	}
	return -b * math.Log(2*(1-u))
}

func addNoise(trueValue float64, epsilon float64, sensitivity float64) float64 {
	// Calculate the scale parameter for the Laplace distribution
	b := sensitivity / epsilon

	// Sample from the Laplace distribution
	noise := sampleLaplace(b)

	// Add noise to the true value
	return trueValue + noise
}

// enforceKAnonymity removes any row whose combination of quasi‑identifiers
// appears fewer than k times. If quasiIDs is empty or nil, it defaults to
// using all visible columns in the table as quasi‑identifiers.
func enforceKAnonymity(table Table, quasiIDs []string, k int) Table {
	// If no quasi‑IDs specified, gather all visible column names
	if len(quasiIDs) == 0 {
		for _, col := range table.Columns {
			if col.Visible {
				quasiIDs = append(quasiIDs, col.Name)
			}
		}
	}

	// check if we have a hidden "count" column
	hasCountCol := false
	for _, col := range table.Columns {
		if col.Name == "count" {
			hasCountCol = true
			break
		}
	}

	var filtered []map[string]interface{}
	if hasCountCol {
		// use the precomputed count value
		for _, row := range table.Rows {
			// assume count is stored as float64
			cnt, _ := row["count"].(float64)
			if int(cnt) >= k {
				filtered = append(filtered, row)
			}
		}
	} else {
		// fallback: recompute frequencies per tuple
		freq := make(map[string]int, len(table.Rows))
		for _, row := range table.Rows {
			parts := make([]string, len(quasiIDs))
			for i, col := range quasiIDs {
				parts[i] = fmt.Sprintf("%v", row[col])
			}
			key := strings.Join(parts, "|")
			freq[key]++
		}
		for _, row := range table.Rows {
			parts := make([]string, len(quasiIDs))
			for i, col := range quasiIDs {
				parts[i] = fmt.Sprintf("%v", row[col])
			}
			key := strings.Join(parts, "|")
			if freq[key] >= k {
				filtered = append(filtered, row)
			}
		}
	}

	table.Rows = filtered
	return table
}

// enforceLDiversity filters out any row whose equivalence class (defined by
// all visible columns) has fewer than l distinct values of the specified
// sensitive attribute.
func enforceLDiversity(table Table, colsToCheck []string, l int) Table {

	// build a set for quick lookup of which cols to enforce
	checkSet := make(map[string]struct{}, len(colsToCheck))
	for _, name := range colsToCheck {
		checkSet[name] = struct{}{}
	}

	// 1) Count unique values for each column in colsToCheck
	uniqueVals := make(map[string]map[interface{}]struct{}, len(colsToCheck))
	for _, colName := range colsToCheck {
		uniqueVals[colName] = make(map[interface{}]struct{})
	}
	for _, row := range table.Rows {
		for _, colName := range colsToCheck {
			uniqueVals[colName][row[colName]] = struct{}{}
		}
	}

	// 2) Build new list of columns:
	//    - if col.Name is in checkSet, keep only if unique count >= l
	//    - otherwise always keep
	keptCols := make([]Column, 0, len(table.Columns))
	for _, col := range table.Columns {
		if _, toCheck := checkSet[col.Name]; toCheck {
			if len(uniqueVals[col.Name]) >= l {
				keptCols = append(keptCols, col)
			}
		} else {
			keptCols = append(keptCols, col)
		}
	}

	// 3) Filter each row to include only the kept columns
	newRows := make([]map[string]interface{}, len(table.Rows))
	for i, row := range table.Rows {
		newRow := make(map[string]interface{}, len(keptCols))
		for _, col := range keptCols {
			newRow[col.Name] = row[col.Name]
		}
		newRows[i] = newRow
	}

	table.Columns = keptCols
	table.Rows = newRows
	return table
}
