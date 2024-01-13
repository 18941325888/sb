package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

type SQLStore struct {
	*Queries
	//sqlc 生成的查询结构,Queries 提供的所有单个查询功能都将可用于 Store.
	db *sql.DB
	//我们可以通过向新结构添加更多功能来支持事务。为此，我们需要 Store 有一个 sql.DB 对象。所以需要创建一个新的数据库事务。
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
		//通过传入的sql.db,用sqlc生成的New函数创建Queries查询结构
	}
}

// 创建一个函数执行通用数据库事务。它需要一个上下文和一个回调函数作为输入，然后它会开始一个新的数据库事务，使用该事务创建一个新的查询对象，并使用创建的查询对象调用回调函数，最后基于该函数返回的错误,提交或回滚事务.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v,rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
