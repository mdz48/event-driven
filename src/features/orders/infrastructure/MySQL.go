package infrastructure

import (
	"database/sql"
	"errors"
	"event-driven/src/features/orders/domain"
)

type MySQL struct {
	db *sql.DB
}

func NewMySQL(db *sql.DB) *MySQL {
	return &MySQL{db: db}
}

func (m *MySQL) GetOrder(orderID int) (domain.Order, error) {
	var order domain.Order
	query := "SELECT id, table_id, product, quantity, status FROM orders WHERE id = ?"
	err := m.db.QueryRow(query, orderID).Scan(&order.ID, &order.TableID, &order.Product, &order.Quantity, &order.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Order{}, errors.New("orden no encontrada")
		}
		return domain.Order{}, err
	}
	return order, nil
}

func (m *MySQL) CreateOrder(order domain.Order) (domain.Order, error) {
	query := "INSERT INTO orders (table_id, product, quantity, status) VALUES (?, ?, ?, ?)"
	result, err := m.db.Exec(query, order.TableID, order.Product, order.Quantity, order.Status)
	if err != nil {
		return domain.Order{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Order{}, err
	}

	order.ID = int(id)
	return order, nil
}

func (m *MySQL) UpdateStatus(orderID int, status string) (domain.Order, error) {
	query := "UPDATE orders SET status = ? WHERE id = ?"
	_, err := m.db.Exec(query, status, orderID)
	if err != nil {
		return domain.Order{}, err
	}

	return m.GetOrder(orderID)
}

func (m *MySQL) GetAll() ([]domain.Order, error) {
	var orders []domain.Order
	query := "SELECT id, table_id, product, quantity, status FROM orders"

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.TableID, &order.Product, &order.Quantity, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (m *MySQL) DeleteOrder(orderID int) error {
	query := "DELETE FROM orders WHERE id = ?"
	result, err := m.db.Exec(query, orderID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("orden no encontrada")
	}

	return nil
}