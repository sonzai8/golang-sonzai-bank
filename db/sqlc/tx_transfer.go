package db

import "context"

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := s.execTx(ctx, func(queries *Queries) error {
		ok := CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		}
		var err error
		//txName := ctx.Value(txKey)

		result.Transfer, err = queries.CreateTransfer(ctx, ok)

		if err != nil {
			return err
		}
		//fmt.Println(txName, "Create entry 1")
		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		//fmt.Println(txName, "Create entry 2")
		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// todo : update account' balance
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err =
				addMoney(ctx, queries, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

		} else {
			result.ToAccount, result.FromAccount, err =
				addMoney(ctx, queries, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

func addMoney(ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,

	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})

	return
}
