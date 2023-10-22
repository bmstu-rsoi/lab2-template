package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"lab2/src/bonus-service/internal/models"
)

type Repository interface {
	GetPrvilegeByUsername(username string) (*models.Privilege, error)
	GetHistoryById(ticketUID string) ([]*models.PrivilegeHistory, error)
	CreateHistoryRecord(*models.PrivilegeHistory) error
	CreatePrivilege(*models.Privilege) error
}

type BonusRepository struct {
	db *sql.DB
}

func NewMySqlRepo(db *sql.DB) *BonusRepository {
	return &BonusRepository{
		db: db,
	}
}

func (r *BonusRepository) CreateHistoryRecord(record *models.PrivilegeHistory) error {

	_, err := r.db.Query(
		`INSERT INTO privilege_history (privilege_id, ticket_uid, datetime, balance_diff, operation_type) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		record.PrivilegeID,
		record.TicketUID,
		record.Date,
		record.BalanceDiff,
		record.OperationType,
	)

	return err
}

func (r *BonusRepository) CreatePrivilege(record *models.Privilege) error {

	_, err := r.db.Query(
		`INSERT INTO privilege (username, balance) VALUES ($1, $2) RETURNING id;`,
		record.Username,
		record.Balance,
	)

	return err
}

func (r *BonusRepository) GetPrvilegeByUsername(username string) (*models.Privilege, error) {

	var privilege models.Privilege

	log.Printf(">>>> username: %s", username)
	row := r.db.QueryRow(`SELECT * FROM privilege where username = $1;`, username)
	err := row.Scan(&privilege.ID, &privilege.Username, &privilege.Status, &privilege.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &privilege, err
		}
	}

	return &privilege, nil
}

func (r *BonusRepository) GetHistoryById(privilegeID string) ([]*models.PrivilegeHistory, error) {

	var history []*models.PrivilegeHistory
	rows, err := r.db.Query(`SELECT * FROM privilege_history where privilege_id = $1;`, privilegeID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute the query: %s", err)
	}

	for rows.Next() {
		row := new(models.PrivilegeHistory)
		rows.Scan(
			&row.ID,
			&row.PrivilegeID,
			&row.TicketUID,
			&row.Date,
			&row.BalanceDiff,
			&row.OperationType,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to execute the query: %s", err)
		}

		history = append(history, row)
	}

	return history, nil
}
