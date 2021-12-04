package adaptor

import (
	"database/sql"
	"finance/entity"
	"fmt"
)

type CategoryAdaptorInterface interface {
	GetAllCategory() error
}

type CategoryMySQLAdaptor struct {
	conn *sql.DB
}

func NewCategoryMySQLAdaptor() CategoryMySQLAdaptor {
	return CategoryMySQLAdaptor{conn: GetMySQLConnection()}
}

func (c CategoryMySQLAdaptor) GetAllCategory() ([]entity.Category, error) {
	rows, err := c.conn.Query(`SELECT * FROM category`)
	if err != nil {
		err = fmt.Errorf("[Error] sql: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []entity.Category
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(&category.CategoryID, &category.Name, &category.IsExpense, &category.IconName); err != nil {
			err = fmt.Errorf("[Error] sql: %v", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, err
}
