package adaptor

import (
	"database/sql"
	"fmt"
)

type Category struct {
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
	IsExpense  bool   `json:"is_expense"`
	IconName   string `json:"icon_name"`
}

type CategoryAdaptorInterface interface {
	GetAllCategory() ([]Category, error)
}

type CategoryMySQLAdaptor struct {
	conn *sql.DB
}

func NewCategoryMySQLAdaptor() CategoryMySQLAdaptor {
	return CategoryMySQLAdaptor{conn: GetMySQLConnection()}
}

func (c CategoryMySQLAdaptor) GetAllCategory() ([]Category, error) {
	rows, err := c.conn.Query(`SELECT * FROM Category`)
	if err != nil {
		err = fmt.Errorf("[Error] GetAllCategory: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.CategoryID, &category.Name, &category.IsExpense, &category.IconName); err != nil {
			err = fmt.Errorf("[Error] GetAllCategory: %v", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, err
}
