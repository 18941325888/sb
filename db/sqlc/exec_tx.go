package db

import (
	"context"
	"fmt"
)

// 创建一个函数执行通用数据库事务。它需要一个上下文和一个回调函数作为输入，然后它会开始一个新的数据库事务，使用该事务创建一个新的查询对象，并使用创建的查询对象调用回调函数，最后基于该函数返回的错误,提交或回滚事务.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err:%v,rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
