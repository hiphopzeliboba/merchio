package user

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"log"
	"merchio/internal/model"
	"merchio/internal/repository"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(ctx context.Context, dsn string) repository.UserRepository {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to connect to db")
		errors.Errorf("failed to connect to db: %v", err)
	}
	return &repo{db: dbc}
}

func (r *repo) IsUserPresent(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func (r *repo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx,
		"SELECT id, username, password, coins, created_at FROM users WHERE username = $1",
		username).Scan(&user.ID, &user.Username, &user.Password, &user.Coins, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) CreateUser(ctx context.Context, username, password string) (int64, error) {
	var id int64

	err := r.db.QueryRow(ctx,
		"INSERT INTO users (username, password, coins) VALUES ($1, $2, 1000) RETURNING id",
		username, password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//
//func (r *repo) TransferCoins(ctx context.Context, fromID, toID uint, amount int) error {
//	tx, err := r.db.BeginTx(ctx, nil)
//	if err != nil {
//		return err
//	}
//	defer tx.Rollback()
//
//	// Проверяем баланс отправителя
//	var senderCoins int
//	err = tx.QueryRowContext(ctx, "SELECT coins FROM users WHERE id = $1 FOR UPDATE", fromID).Scan(&senderCoins)
//	if err != nil {
//		return err
//	}
//
//	if senderCoins < amount {
//		return errors.New("insufficient coins")
//	}
//
//	// Обновляем балансы
//	_, err = tx.ExecContext(ctx, "UPDATE users SET coins = coins - $1 WHERE id = $2", amount, fromID)
//	if err != nil {
//		return err
//	}
//
//	_, err = tx.ExecContext(ctx, "UPDATE users SET coins = coins + $1 WHERE id = $2", amount, toID)
//	if err != nil {
//		return err
//	}
//
//	// Записываем транзакцию
//	_, err = tx.ExecContext(ctx,
//		"INSERT INTO coin_transactions (from_id, to_id, amount) VALUES ($1, $2, $3)",
//		fromID, toID, amount)
//	if err != nil {
//		return err
//	}
//
//	return tx.Commit()
//}
//
//func (r *repo) BuyMerch(ctx context.Context, userID uint, merchName string) error {
//	tx, err := r.db.BeginTx(ctx, nil)
//	if err != nil {
//		return err
//	}
//	defer tx.Rollback()
//
//	// Получаем информацию о товаре
//	var merchID uint
//	var price int
//	err = tx.QueryRowContext(ctx,
//		"SELECT id, price FROM merch_items WHERE name = $1",
//		merchName).Scan(&merchID, &price)
//	if err != nil {
//		return err
//	}
//
//	// Проверяем баланс пользователя
//	var userCoins int
//	err = tx.QueryRowContext(ctx,
//		"SELECT coins FROM users WHERE id = $1 FOR UPDATE",
//		userID).Scan(&userCoins)
//	if err != nil {
//		return err
//	}
//
//	if userCoins < price {
//		return errors.New("insufficient coins")
//	}
//
//	// Обновляем баланс
//	_, err = tx.ExecContext(ctx,
//		"UPDATE users SET coins = coins - $1 WHERE id = $2",
//		price, userID)
//	if err != nil {
//		return err
//	}
//
//	// Добавляем товар в инвентарь
//	_, err = tx.ExecContext(ctx,
//		`INSERT INTO user_inventory (user_id, merch_id, quantity)
//         VALUES ($1, $2, 1)
//         ON CONFLICT (user_id, merch_id)
//         DO UPDATE SET quantity = user_inventory.quantity + 1`,
//		userID, merchID)
//	if err != nil {
//		return err
//	}
//
//	return tx.Commit()
//}
