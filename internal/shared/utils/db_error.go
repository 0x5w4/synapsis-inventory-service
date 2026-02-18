package utils

import "inventory-service/internal/shared/exception"

// handleDBError is a utility function to handle database errors.
func handleDBError(err error, table string, operation string) error {
	return exception.NewDBError(err, table, operation)
}
