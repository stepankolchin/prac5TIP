package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// .env не обязателен; если файла нет --- ошибка игнорируется
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// fallback --- прямой DSN в коде (только для учебного стенда!)
		dsn = "postgres://postgres:353722Zz@localhost:5432/todo?sslmode=disable"
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// 1) Вставим пару задач
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	// 2) Прочитаем список задач
	ctxList, cancelList := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err := repo.ListTasks(ctxList)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}
	
	// ТЕСТИРУЕМ НОВЫЕ ФУНКЦИИ
	
	// 1. Тест ListDone - выведем только невыполненные задачи
	fmt.Println("\n=== Невыполненные задачи (ListDone) ===")
	undoneTasks, err := repo.ListDone(ctxList, false)
	if err != nil {
		log.Printf("ListDone error: %v", err)
	} else {
		for _, t := range undoneTasks {
			fmt.Printf("#%d | %s\n", t.ID, t.Title)
		}
	}

	// 2. Тест FindByID - найдем задачу с ID=2
	fmt.Println("\n=== Поиск задачи с ID=2 (FindByID) ===")
	task, err := repo.FindByID(ctxList, 2)
	if err != nil {
		log.Printf("FindByID error: %v", err)
	} else {
		fmt.Printf("Найдена: #%d | %-24s | done=%-5v | %s\n", 
			task.ID, task.Title, task.Done, task.CreatedAt.Format(time.RFC3339))
	}

	// 3. Тест CreateMany - массовое добавление задач
	fmt.Println("\n=== Массовое добавление (CreateMany) ===")
	newTitles := []string{"Новая задача 1", "Новая задача 2", "Новая задача 3"}
	err = repo.CreateMany(ctxList, newTitles)
	if err != nil {
		log.Printf("CreateMany error: %v", err)
	} else {
		fmt.Println("Успешно добавлено 3 новые задачи")
		
		// Покажем обновленный список
		fmt.Println("\n=== Обновленный список всех задач ===")
		allTasks, err := repo.ListTasks(ctxList)
		if err != nil {
			log.Printf("ListTasks error: %v", err)
		} else {
			for _, t := range allTasks {
				fmt.Printf("#%d | %-24s | done=%-5v | %s\n", 
					t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
			}
		}
	}
}
	
