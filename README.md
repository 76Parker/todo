# Приложение TODO LIST
Приложение для ведения задач

Поддерживает тегирование, сортировку и фильтрацию задач

<img width="2030" height="1168" alt="image" src="https://github.com/user-attachments/assets/c3740264-2cac-44ed-8362-5e73a3553de9" />


## Таблица покрытия тестами основных пакетов
| Название пакета | % покрытия |
|---|---:|
| `todo/internal/api` | **68.9%** |
| `todo/internal/api/middleware` | **100.0%** | 
| `todo/internal/api/router` | **100.0%** |
| `todo/internal/api/validator` | **39.1%** |
| `todo/internal/app` | **18.6%** |
| `todo/internal/entities/domain` | **88.2%** |
| `todo/internal/usecase/task` | **93.3%** |
## Сборка
Сборка приложения не требуют никаких внешних зависимостей и делается через `Makefile`:

`make docker-build` - сборка Docker-образа

`make linux-build` - сборка бинарного файла для платформы`linux/amd64`

`make arm-build` - сборка бинарного файла для платформы `darwin/arm64` (MacOS)

## Основные используемые библиотеки

`go-playground/validator` - библиотека для валидации структур по тегам. Используется для валидации входящих DTO в API приложения

`jackc/pgx` - библиотека для взаимодействия с PostgreSQL. Используется в реализации Repository-интерфейсов

`go.uber.org/mock` - библиотека для генерации моков которые используется в тестах

`Masterminds/squirrel` - удобный SQL-билдер для построения динамических SQL-запросов




