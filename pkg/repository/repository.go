package repository

import "github.com/jmoiron/sqlx"

 type AddCategoryDB interface{

 }
 
 type AddPurchaseDB interface{
 
 }
 
 type Repository struct{
	AddCategoryDB
	AddPurchaseDB
 }
 
 func NewRepository(configDB *sqlx.DB) *Repository{
	return &Repository{}
 }