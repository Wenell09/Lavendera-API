package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RLSMiddleware(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil tenant_id dari Gin Context yang diset oleh JWTMiddleware
		tenantIDVal, exists := c.Get("tenant_id")
		if !exists || tenantIDVal == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Konteks tenant_id tidak ditemukan"})
			return
		}

		tenantID := fmt.Sprintf("%v", tenantIDVal)
		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "tenant_id kosong"})
			return
		}

		ctx := c.Request.Context()

		// 2. Ambil koneksi dari PGX Pool
		conn, err := db.Acquire(ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database connection pool exhausted"})
			return
		}
		defer conn.Release()

		// 3. Mulai Transaksi PostgreSQL
		tx, err := conn.Begin(ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
			return
		}

		// 4. Inject tenant_id ke Session Variable PostgreSQL
		_, err = tx.Exec(ctx, "SET LOCAL app.current_tenant_id = $1;", tenantID)
		if err != nil {
			_ = tx.Rollback(ctx)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to set RLS context"})
			return
		}

		// 5. Simpan transaksi PostgreSQL ke Gin Context untuk digunakan di Handler
		c.Set("pgx_tx", tx)

		c.Next()

		// 6. Auto-Commit / Rollback sesuai HTTP status response
		if c.Writer.Status() < 400 {
			_ = tx.Commit(ctx)
		} else {
			_ = tx.Rollback(ctx)
		}
	}
}
