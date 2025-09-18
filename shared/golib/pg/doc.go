// Package pg provides enterprise-grade PostgreSQL database connectivity using pgx/v5 native interface.
//
// Features:
//   - Connection pooling with health monitoring
//   - Automatic retries with exponential backoff
//   - Transaction management with rollback safety
//   - OpenTelemetry integration (traces, metrics, logs)
//   - Comprehensive error handling
//   - Test utilities and mocking support
//
// Basic usage:
//
//	cfg := pg.Config{
//		Host:     "localhost",
//		Port:     5432,
//		Database: "mydb",
//		Username: "user",
//		Password: "pass",
//	}
//
//	client, err := pg.NewClient(ctx, cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	var count int
//	err = client.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
package pg
