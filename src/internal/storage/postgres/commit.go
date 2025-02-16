package postgres

func (t *StorageTx) Commit() error {
	return t.tx.Commit()
}
