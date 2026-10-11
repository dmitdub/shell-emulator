# shell-emulator

Эмулятор оболочки UNIX-подобной операционной системы на языке Go.

Учебный проект по дисциплине «Конфигурационное управление».
Вариант 7, ИКБО-10-25, РТУ МИРЭА.

## Текущий статус

Реализованы этапы 1, 2 и 3.

- **Этап 1 — REPL.** Диалог с пользователем, разбор аргументов,
  раскрытие переменных окружения, команды-заглушки.
- **Этап 2 — Конфигурация.** Параметры командной строки,
  отладочный вывод, стартовый скрипт, команда `conf-dump`.
- **Этап 3 — VFS.** Загрузка виртуальной файловой системы
  из CSV-файла, работа с деревом в памяти, команда `vfs-init`.

## Возможности

- Интерактивный REPL с приглашением вида `username@hostname:~$`.
- Приглашение формируется из реальных данных операционной системы.
- Раскрытие переменных окружения (`$HOME`, `$USER` и другие).
- Разбор аргументов с поддержкой одинарных и двойных кавычек.
- Сообщение об ошибке при неизвестной команде или незакрытой кавычке.
- Редактирование строки в любом месте (курсор, backspace в середине).
- История команд (стрелки вверх/вниз).
- Параметры командной строки `--vfs` и `--script`.
- Отладочный вывод конфигурации при запуске.
- Команда `conf-dump` для просмотра конфигурации.
- Стартовый скрипт с поддержкой комментариев.
- Загрузка виртуальной файловой системы из CSV в память.
- Команда `vfs-init` для сброса VFS на дефолтную.

## Команды

| Команда | Назначение |
|---|---|
| `ls [аргументы...]` | Заглушка: печатает своё имя и аргументы |
| `cd [аргументы...]` | Заглушка: печатает своё имя и аргументы |
| `conf-dump` | Вывести текущие параметры конфигурации и VFS |
| `vfs-init` | Сбросить VFS на дефолтную (и перезаписать CSV) |
| `exit` | Завершить эмулятор |

## Параметры командной строки

| Флаг | Назначение | Обязательный |
|---|---|---|
| `--vfs` | Путь к виртуальной файловой системе | нет |
| `--script` | Путь к стартовому скрипту | нет |

Оба флага можно писать через пробел (`--vfs path`) или через
знак равенства (`--vfs=path`).

## Структура проекта

```
shell-emulator/
├── src/
│   ├── main.go          Точка входа
│   ├── config/          Парсинг параметров командной строки
│   ├── prompt/          Формирование приглашения
│   ├── parser/          Разбор командной строки
│   ├── repl/            Цикл REPL и диспетчер команд
│   ├── script/          Загрузка стартовых скриптов
│   └── vfs/             Виртуальная файловая система
├── tests/               Модульные тесты
├── testdata/            Тестовые данные
│   ├── startup.txt      Пример стартового скрипта
│   ├── vfs.csv          Основной VFS
│   ├── vfs-minimal.csv  Минимальный VFS
│   ├── vfs-deep.csv     VFS с 4 уровнями вложенности
│   └── vfs-broken.csv   Битый VFS для проверки ошибок
├── scripts/             Bash-скрипты для запуска эмулятора
├── Makefile             Сборка, тесты, запуск
├── go.mod               Модуль Go
└── README.md
```

## Требования

- Go 1.21 или новее.
- Git (для клонирования репозитория).
- Bash (для скриптов в `scripts/`).

## Сборка

Собрать исполняемый файл:

```bash
make build
```

Бинарник появится в `bin/shell-emulator` (или `bin/shell-emulator.exe` в Windows).

## Запуск

Запустить из исходников:

```bash
make run
```

Запустить с флагами:

```bash
make run-vfs
make run-script
make run-flags
```

Запустить с разными VFS:

```bash
make run-vfs-minimal
make run-vfs-deep
make run-vfs-broken
```

Или вручную:

```bash
go run ./src/
go run ./src/ --vfs testdata/vfs.csv
go run ./src/ --script testdata/startup.txt
go run ./src/ --vfs testdata/vfs.csv --script testdata/startup.txt
go run ./src/ --vfs testdata/vfs-minimal.csv
go run ./src/ --vfs testdata/vfs-deep.csv
go run ./src/ --vfs testdata/vfs-broken.csv
```

> **Примечание для Windows.** При запуске через `go run ./src/`
> команда `exit` может сработать только со второго раза — это
> особенность обёртки `go run`, а не эмулятора. Для реальной
> работы используйте `make build` и запускайте бинарник
> `bin/shell-emulator.exe` напрямую.

## Тесты

Запустить все тесты:

```bash
make test
```

С подробным выводом:

```bash
make test-v
```

Или вручную:

```bash
go test ./tests/
go test -v ./tests/
```

## Примеры использования

### Простой вызов команды

```
User@host:~$ ls -la /tmp
ls
  -la
  /tmp
```

### Аргумент с пробелами

```
User@host:~$ cd "/home/my folder"
cd
  /home/my folder
```

### Раскрытие переменных окружения

```
User@host:~$ cd $HOME
cd
  /home/user
```

### Неизвестная команда

```
User@host:~$ echo hello
command not found: echo
```

### Ошибка разбора

```
User@host:~$ ls "unclosed
parse error: unbalanced quote
```

### Просмотр конфигурации

```
User@host:~$ conf-dump
vfs.path = testdata/vfs.csv
script.path = 
vfs.root = /
vfs.children = 2
```

### Выход

```
User@host:~$ exit
```

## Этап 2 — Конфигурация

### Отладочный вывод

При запуске эмулятор выводит все параметры конфигурации в stderr:

```
configuration:
vfs.path = testdata/vfs.csv
script.path = testdata/startup.txt
vfs.root = /
vfs.children = 2
```

Отладочный вывод идёт именно в stderr, чтобы не смешиваться
с выводом команд.

### Стартовый скрипт

Файл с одной командой эмулятора на строке. Строки, начинающиеся
с `//`, считаются комментариями и не выполняются. Пустые строки
игнорируются. Пробелы по краям строк обрезаются.

Пример `testdata/startup.txt`:

```
// sample startup script
ls /tmp
cd /home/user
ls "my documents"
conf-dump
vfs-init
conf-dump
exit
```

Каждая команда эхо-печатается с приглашением перед выполнением,
имитируя диалог с пользователем.

### Скрипты реальной ОС

В папке `scripts/` лежат bash-скрипты, которые запускают эмулятор
с разными комбинациями флагов:

| Скрипт | Что проверяет |
|---|---|
| `run-no-flags.sh` | запуск без параметров |
| `run-with-vfs.sh` | флаг `--vfs` |
| `run-with-script.sh` | флаг `--script` |
| `run-all-flags.sh` | оба флага |
| `run-missing-script.sh` | ошибка при отсутствующем скрипте |
| `run-vfs-minimal.sh` | VFS минимальный |
| `run-vfs-deep.sh` | VFS глубокий |
| `run-vfs-broken.sh` | битый VFS → ошибка |
| `run-vfs-init.sh` | сброс VFS через `vfs-init` |
| `run-all.sh` | прогон всех неинтерактивных скриптов |

Запуск всех скриптов:

```bash
make demo
```

Или вручную:

```bash
bash scripts/run-all.sh
```

## Этап 3 — VFS

Виртуальная файловая система загружается из CSV-файла в память.
Файл на диске **не модифицируется**, кроме случая команды `vfs-init`.

### Формат CSV

```
id,parent,name,type,content
0,-1,/,dir,
1,0,home,dir,
2,1,user,dir,
3,2,notes.txt,file,Hello world
```

| Колонка | Значение |
|---|---|
| `id` | уникальный идентификатор узла |
| `parent` | идентификатор родителя; `-1` — корень |
| `name` | имя файла или директории |
| `type` | `dir` или `file` |
| `content` | содержимое файла (пусто для директорий) |

### Тестовые VFS

В папке `testdata/` лежат три варианта VFS:

| Файл | Что содержит |
|---|---|
| `vfs-minimal.csv` | только корень `/` |
| `vfs.csv` | несколько папок и файлов |
| `vfs-deep.csv` | 4 уровня вложенности |

Плюс `vfs-broken.csv` — для проверки обработки ошибок.

### Запуск с разными VFS

```bash
make run-vfs           # основной (vfs.csv)
make run-vfs-minimal   # минимальный
make run-vfs-deep      # глубокий
make run-vfs-broken    # битый → ошибка
```

Или вручную:

```bash
go run ./src/ --vfs testdata/vfs-minimal.csv
go run ./src/ --vfs testdata/vfs-deep.csv
go run ./src/ --vfs testdata/vfs-broken.csv
```

### Команда vfs-init

```
User@host:~$ vfs-init
VFS reset to default and saved to testdata/vfs.csv
```

- Заменяет VFS в памяти на дефолтную
- **Перезаписывает** исходный CSV-файл (по требованию методички)
- Если `--vfs` не задан — работает только в памяти

### Ошибки загрузки VFS

| Ситуация | Сообщение |
|---|---|
| Файл не найден | `load VFS: open ...: no such file or directory` |
| Битый CSV | `load VFS: read rows: ...` |
| Неверный тип | `load VFS: row N: invalid type "xxx"` |
| Нет корня | `load VFS: root node not found` |
| Битый parent | `load VFS: node N: parent M not found` |

## Лицензия

Учебный проект, свободное использование в образовательных целях.