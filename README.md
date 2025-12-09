## Практическое занятие №5. Колчин Степан Сергеевич, ЭФМО-02-25. Подключение к PostgreSQL через database/sql. Выполнение простых запросов (INSERT, SELECT)

### Окружение

- Go: go version go1.25.3 windows/amd64

- PostgreSQL: psql (PostgreSQL) 14+

- ОС: Windows 11

### Скриншоты

- создание БД/таблицы в psql

<img width="481" height="201" alt="screen1" src="https://github.com/user-attachments/assets/f8b4aa42-33c5-4159-bf30-692ef7b0775b" />

- успешный вывод go run . (вставка и список задач);

<img width="569" height="457" alt="screen2" src="https://github.com/user-attachments/assets/30b70bd1-930e-4851-94f7-3c8198215d7f" />

- SELECT * FROM tasks; в psql

<img width="608" height="219" alt="screen3" src="https://github.com/user-attachments/assets/0cbb3552-c809-428a-af55-de0996f1d325" />

### Код:

```GO
// db.go
package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Проверка соединения с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}
```
```GO
// Реализация фунции ListDone

// ListDone возвращает только выполненные (done=true) или невыполненные (done=false) задачи
func (r *Repo) ListDone(ctx context.Context, done bool) ([]Task, error) {
    const q = `SELECT id, title, done, created_at FROM tasks WHERE done = $1 ORDER BY id;`
    
    rows, err := r.DB.QueryContext(ctx, q, done)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []Task
    for rows.Next() {
        var t Task
        if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
            return nil, err
        }
        out = append(out, t)
    }
    return out, rows.Err()
}
```

```GO
//Фрагмент main.go
func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:123@localhost:5432/todo?sslmode=disable"
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// Вставка задач
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

	// Вывод всех задач
	tasks, err := repo.ListTasks(ctx)
	// ... вывод результатов

	// ТЕСТИРОВАНИЕ ПРОВЕРОЧНЫХ ЗАДАНИЙ
	// ... код тестирования ListDone, FindByID, CreateMany
}
```

### Краткие ответы:
- Что такое пул соединений *sql.DB и зачем его настраивать?

**Пул соединений `*sql.DB`** — это набор заранее установленных подключений к БД, которые переиспользуются между запросами. Настройка нужна для оптимизации производительности: ограничение максимального числа соединений предотвращает перегрузку БД, а поддержание соединений в простое уменьшает задержки при новых запросах.

- Почему используем плейсхолдеры $1, $2?

**Плейсхолдеры $1, $2** - для защиты от SQL-инъекций. 

- Чем `Query`, `QueryRow` и `Exec` отличаются?

`Query` — для множества строк

`QueryRow` — для одной строки

`Exec` — для 0 строк

### Обоснование транзакций и настроек пула:

- `SetMaxOpenConns(10)`:

    - Для чего: Ограничение максимального числа одновременных соединений
    - Обоснование: Предотвращение перегрузки БД при высокой нагрузке

- `SetMaxIdleConns(5)`:

    - Для чего: Поддержание соединений в "простое" для быстрого переиспользования
    - Обоснование: Уменьшение задержек на установку нового соединения

- `SetConnMaxLifetime(30 * time.Minute)`:

    - Для чего: Ограничение времени жизни соединения
    - Обоснование: Предотвращение использования "устаревших" соединений
