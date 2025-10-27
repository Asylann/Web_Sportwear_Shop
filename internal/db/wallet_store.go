package db

import (
	"WebSportwareShop/internal/models"
	"context"
	"fmt"
)

func CreateWalletByUserId(ctx context.Context, userId int) (int, error) {
	var id int
	err := db.QueryRowContext(ctx, "INSERT INTO wallets(user_id) VALUES($1)", userId).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetWalletByUserId(ctx context.Context, userId int) (models.Wallet, error) {
	var wallet models.Wallet
	err := db.QueryRowxContext(ctx, "SELECT * FROM wallets WHERE user_id=$1", userId).Scan(&wallet.Id, &wallet.UserId, &wallet.Balance, &wallet.Currency)
	if err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func MakeAPayment(ctx context.Context, FromWhoId int, ToWhomId int, amount float64) error {
	if FromWhoId == ToWhomId {
		return fmt.Errorf("cannot transfer to same user")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be > 0")
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var w1, w2 models.Wallet

	if err = tx.GetContext(ctx, &w1, "SELECT id, user_id, balance FROM wallets WHERE user_id=$1 FOR UPDATE", FromWhoId); err != nil {
		return err
	}
	if err = tx.GetContext(ctx, &w2, "SELECT id, user_id, balance FROM wallets WHERE user_id=$1 FOR UPDATE", ToWhomId); err != nil {
		return err
	}

	// map to from/to
	var fromWalletID, toWalletID int
	var fromBalance, toBalance float64
	fromWalletID, fromBalance = w1.Id, w1.Balance
	toWalletID, toBalance = w2.Id, w2.Balance

	if fromBalance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// update balances
	if _, err := tx.ExecContext(ctx, "UPDATE wallets SET balance = $1 WHERE id = $2", fromBalance-amount, fromWalletID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE wallets SET balance = $1 WHERE id = $2", toBalance+amount, toWalletID); err != nil {
		return err
	}

	// audit row
	if _, err := tx.ExecContext(ctx, "INSERT INTO transactions (from_wallet_id, to_wallet_id, amount) VALUES ($1, $2, $3)",
		fromWalletID, toWalletID, amount); err != nil {
		return err
	}

	return tx.Commit()
}
