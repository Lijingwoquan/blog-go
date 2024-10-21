package mysql

import "fmt"

func noAffectedRowErr(affectedRows int64, oldErr error, msg string) error {
	if oldErr != nil {
		return oldErr
	}
	if affectedRows == 0 {
		return fmt.Errorf(msg)
	}
	return nil
}
