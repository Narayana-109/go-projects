package api

import (
	"fmt"
	"context"
	"net/http"
	"trailog/database"
	"github.com/gin-gonic/gin"
)

func LoadRouter() {
	router := gin.Default()
	router.POST("/employee/add", AddEmployee)
	router.PUT("/employee/update/:email",UpdateEmployee)
	router.Run(":8000")
}

type Employee struct{
	employee_id string 		`json:"employee_id"`
	first_name string 		`json:"first_name"`
	last_name string 		`json:"last_name"`
	email string 			`json:"email"`
}

func AddEmployee(c *gin.Context) {
	var emp Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Access the exported Pool from the database package
	query := "INSERT INTO employees (employee_id, first_name, last_name, email) VALUES ($1, $2, $3, $4)"
	_, err := database.DB.Exec(context.Background(), query, emp.employee_id, emp.first_name, emp.last_name, emp.email)
	
	if err != nil {
		fmt.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Employee added successfully"})
}

func UpdateEmployee(c *gin.Context) {
	var emp Employee
	currentEmail := c.Param("email")
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	query := "UPDATE employees SET employee_id=$1, first_name=$2, last_name=$3, email=$4 where email=$5"
	result, err := database.DB.Exec(context.Background(), query, emp.employee_id, emp.first_name, emp.last_name, emp.email, currentEmail)

	if err != nil {
		fmt.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Employee not found"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Employee updated successfully"})
}