package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetManpowerRequests(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`SELECT request_id, doc_no, position_title, num_required, overall_status FROM manpower_requests`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		for rows.Next() {
			var id int
			var docNo, positionTitle, status string
			var num int
			if err := rows.Scan(&id, &docNo, &positionTitle, &num, &status); err == nil {
				results = append(results, map[string]interface{}{
					"id":       id,
					"doc_no":   docNo,
					"title":    positionTitle,
					"num":      num,
					"status":   status,
				})
			}
		}

		c.JSON(http.StatusOK, results)
	}
}

func CreateManpowerRequest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			DocNo    string `json:"doc_no"`
			Title    string `json:"title"`
			Num      int    `json:"num"`
			Reason   string `json:"reason"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(
			`INSERT INTO manpower_requests (doc_no, doc_date, department_id, requested_by, position_title, num_required, reason)
			 VALUES ($1, CURRENT_DATE, 1, 1, $2, $3, $4)`,
			req.DocNo, req.Title, req.Num, req.Reason,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "request created"})
	}
}
