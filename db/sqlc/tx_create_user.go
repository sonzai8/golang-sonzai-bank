package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction

func (s *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := s.execTx(ctx, func(queries *Queries) error {
		var err error
		result.User, err = queries.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.User)

	})
	return result, err
}
