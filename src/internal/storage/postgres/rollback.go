package postgres

func (t *StorageTx) Rollback() error {
	return t.tx.Rollback()
}
