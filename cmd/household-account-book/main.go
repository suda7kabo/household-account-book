package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/suda7kabo/household-account-book/handler"
	"github.com/suda7kabo/household-account-book/infrastructure"
	"github.com/suda7kabo/household-account-book/usecase"
)

func main() {
	e := echo.New()

	db, err := infrastructure.NewDB()
	if err != nil {
		log.Fatalln("$$failed to connect db", err)
	}
	defer db.Close()
	expenseRepository := infrastructure.NewExpenseRepository(db)
	expenseUseCase := usecase.NewExpenseUseCase(expenseRepository)
	expenseHandler := handler.NewExpenseHandler(expenseUseCase)

	e.POST("/expenses", expenseHandler.CreateExpense)
	go func() {
		if err := e.Start(":1323"); err != http.ErrServerClosed {
			log.Fatalln("Server closed with error:", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		panic(fmt.Errorf("failed to graceful shutdown:%w", err))
	}
	log.Println("Server shutdown")
}
