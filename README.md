## Практическое занятие №5. Колчин Степан Сергеевич, ЭФМО-02-25. Подключение к PostgreSQL через database/sql. Выполнение простых запросов (INSERT, SELECT)

### Окружение

- Go: go version go1.25.3 windows/amd64

- PostgreSQL: psql (PostgreSQL) 14+

- ОС: Windows 11

### Скриншоты

- создание БД/таблицы в psql


- успешный вывод go run . (вставка и список задач);





- SELECT * FROM tasks; в psql



### Код:

- [db.go](./db.go)
- [repository.go](./repository.go)

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

### 6 Обоснование транзакций и настроек пула:

- `SetMaxOpenConns(10)`:

    - Для чего: Ограничение максимального числа одновременных соединений
    - Обоснование: Предотвращение перегрузки БД при высокой нагрузке

- `SetMaxIdleConns(5)`:

    - Для чего: Поддержание соединений в "простое" для быстрого переиспользования
    - Обоснование: Уменьшение задержек на установку нового соединения

- `SetConnMaxLifetime(30 * time.Minute)`:

    - Для чего: Ограничение времени жизни соединения
    - Обоснование: Предотвращение использования "устаревших" соединений