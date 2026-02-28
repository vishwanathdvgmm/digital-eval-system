package block

// transaction.go kept minimal for now. Transaction struct defined in block.go.
// Add helpers here if needed (validation, canonicalization).

// ValidateTransaction performs basic required field checks.
func ValidateTransaction(tx *Transaction) bool {
	if tx == nil {
		return false
	}
	if tx.ScriptID == "" || tx.USN == "" || tx.CourseID == "" || tx.Semester == "" || tx.CID == "" {
		return false
	}
	return true
}
