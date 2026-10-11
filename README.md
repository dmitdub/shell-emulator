# shell-emulator

Эмулятор оболочки UNIX-подобной операционной системы на языке Go.

Учебный проект по дисциплине «Конфигурационное управление».
Вариант 7, ИКБО-10-25, РТУ МИРЭА.

## Текущий статус

Реализованы все пять этапов.

- **Этап 1 — REPL.** Диалог с пользователем, разбор аргументов,
  раскрытие переменных окружения.
- **Этап 2 — Конфигурация.** Параметры командной строки,
  отладочный вывод, стартовый скрипт, команда `conf-dump`.
- **Этап 3 — VFS.** Загрузка виртуальной файловой системы
  из CSV-файла, работа с деревом в памяти, команда `vfs-init`.
- **Этап 4 — Основные команды.** Реальные `ls`, `cd`, `uptime`, `cat`.
- **Этап 5 — Дополнительные команды.** Копирование `cp`.

## Возможности

- Интерактивный REPL с приглашением вида `username@hostname:/path$`.
- Приглашение формируется из реальных данных ОС и текущего пути в VFS.
- Раскрытие переменных окружения (`$HOME`, `$USER` и другие).
- Разбор аргументов с поддержкой одинарных и двойных кавычек.
- Редактирование строки в любом месте (курсор, backspace).
- История команд (стрелки вверх/вниз).
- Параметры командной строки `--vfs` и `--script`.
- Отладочный вывод конфигурации при запуске.
- Команда `conf-dump` для просмотра конфигурации и состояния VFS.
- Стартовый скрипт с поддержкой комментариев.
- Виртуальная файловая система в памяти, загружаемая из CSV.
- Команда `vfs-init` для сброса VFS на дефолтную.
- Команды `ls` (с флагами `-l`, `-a`, `-h`), `cd`, `uptime`, `cat`, `cp`.

## Команды

| Команда | Назначение |
|---|---|
| `ls [-lah] [путь...]` | Показать содержимое директории |
| `cd [путь]` | Сменить текущую директорию |
| `cat файл...` | Показать содержимое файла |
| `cp [-r] источник цель` | Скопировать файл или директорию |
| `uptime` | Показать время работы эмулятора |
| `conf-dump` | Вывести конфигурацию и состояние VFS |
| `vfs-init` | Сбросить VFS на дефолтную |
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
│   ├── commands/        Реализация команд (ls, cd, uptime, cat, cp)
│   ├── config/          Парсинг параметров командной строки
│   ├── prompt/          Формирование приглашения
│   ├── parser/          Разбор командной строки
│   ├── repl/            Цикл REPL и диспетчер команд
│   ├── script/          Загрузка стартовых скриптов
│   └── vfs/             Виртуальная файловая система
├── tests/               Модульные тесты
├── testdata/            Тестовые данные
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

```bash
make build
```

Бинарник появится в `bin/shell-emulator` (или `bin/shell-emulator.exe` в Windows).

## Запуск

```bash
make run
make run-vfs
make run-script
make run-flags
make run-vfs-minimal
make run-vfs-deep
make run-vfs-broken
make run-commands
make run-cp
```

Или вручную:

```bash
go run ./src/
go run ./src/ --vfs testdata/vfs.csv
go run ./src/ --script testdata/startup.txt
go run ./src/ --vfs testdata/vfs.csv --script testdata/startup.txt
```

> **Примечание для Windows.** При запуске через `go run ./src/`
> команда `exit` может сработать только со второго раза — это
> особенность обёртки `go run`. Для реальной работы используйте
> `make build` и бинарник `bin/shell-emulator.exe`.

## Тесты

```bash
make test
make test-v
```

Или вручную:

```bash
go test ./tests/
go test -v ./tests/
```

## Этап 2 — Конфигурация

### Отладочный вывод

При запуске эмулятор выводит параметры конфигурации в stderr:

```
configuration:
vfs.path = testdata/vfs.csv
script.path = 
vfs.root = /
vfs.children = 2
```

### Стартовый скрипт

Файл с одной командой на строке. Строки `//` — комментарии.
Пустые строки игнорируются.

Пример `testdata/startup.txt`:

```
ls /
cd home
ls
cat readme.txt
cp readme.txt /tmp
cd /tmp
ls
cp -r /home /tmp
ls
cd /
uptime
exit
```

Каждая команда эхо-печатается с приглашением.

### Скрипты реальной ОС

В папке `scripts/`:

| Скрипт | Что проверяет |
|---|---|
| `run-no-flags.sh` | запуск без параметров |
| `run-with-vfs.sh` | флаг `--vfs` |
| `run-with-script.sh` | флаг `--script` |
| `run-all-flags.sh` | оба флага |
| `run-missing-script.sh` | ошибка при отсутствующем скрипте |
| `run-vfs-minimal.sh` | VFS минимальный |
| `run-vfs-deep.sh` | VFS глубокий |
| `run-vfs-broken.sh` | битый VFS |
| `run-vfs-init.sh` | сброс VFS через `vfs-init` |
| `run-commands.sh` | все команды этапа 4 |
| `run-cp.sh` | команды этапа 5 |
| `run-all.sh` | прогон всех неинтерактивных скриптов |

Запуск:

```bash
make demo
bash scripts/run-all.sh
```

## Этап 3 — VFS

Виртуальная файловая система загружается из CSV-файла в память.
Файл на диске **не модифицируется**, кроме команды `vfs-init`.

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
| `id` | уникальный идентификатор |
| `parent` | идентификатор родителя; `-1` — корень |
| `name` | имя файла или директории |
| `type` | `dir` или `file` |
| `content` | содержимое (пусто для директорий) |

### Тестовые VFS

| Файл | Что содержит |
|---|---|
| `vfs-minimal.csv` | только корень `/` |
| `vfs.csv` | несколько папок и файлов |
| `vfs-deep.csv` | 4 уровня вложенности |
| `vfs-broken.csv` | битый VFS для проверки ошибок |

### Команда vfs-init

```
User@host:/$ vfs-init
VFS reset to default and saved to testdata/vfs.csv
```

- Заменяет VFS в памяти на дефолтную
- **Перезаписывает** исходный CSV-файл
- Если `--vfs` не задан — работает только в памяти

### Ошибки загрузки VFS

| Ситуация | Сообщение |
|---|---|
| Файл не найден | `load VFS: open ...: no such file or directory` |
| Битый CSV | `load VFS: read rows: ...` |
| Неверный тип | `load VFS: row N: invalid type "xxx"` |
| Нет корня | `load VFS: root node not found` |
| Битый parent | `load VFS: node N: parent M not found` |

## Этап 4 — Основные команды

### ls

| Флаг | Что делает |
|---|---|
| `-l` | Длинный формат: права, размер, имя |
| `-a` | Показывать скрытые файлы |
| `-h` | Читаемый размер (работает с `-l`) |

Флаги комбинируются: `ls -lah`.

```
User@host:/$ ls /
home/  tmp/
User@host:/$ ls -l /
drwxr-xr-x  -  home
drwxr-xr-x  -  tmp
User@host:/$ ls -lah /home
-rw-r--r--  22  readme.txt
```

### cd

```
User@host:/$ cd home
User@host:/home$ cd ..
User@host:/$
```

Ошибки:
- `cd /no/such` → `cd: /no/such: No such file or directory`
- `cd /home/readme.txt` → `cd: /home/readme.txt: Not a directory`

### uptime

```
User@host:/$ uptime
up 2 minutes and 15 seconds
```

### cat

```
User@host:/$ cat /home/readme.txt
Welcome to the emulator.
User@host:/$ cat /home
cat: /home: Is a directory
User@host:/$ cat /no/such/file
cat: /no/such/file: No such file or directory
```

### Приглашение

```
User@host:/$ 
User@host:/home$ 
User@host:/home/user$ 
```

## Этап 5 — Копирование

### cp

Копирует файл или директорию внутри VFS. Все изменения
происходят **только в памяти** — исходные CSV-файлы не трогаются.

| Флаг | Что делает |
|---|---|
| `-r` / `-R` | Рекурсивное копирование директорий |

#### Копирование файла

```
User@host:/$ cp /home/readme.txt /tmp
User@host:/$ ls /tmp
readme.txt
User@host:/$ cat /tmp/readme.txt
Welcome to the emulator.
```

#### Копирование с новым именем

```
User@host:/$ cp /home/readme.txt /backup.txt
User@host:/$ ls /
backup.txt  home/  tmp/
```

#### Копирование директории

```
User@host:/$ cp -r /home /backup
User@host:/$ ls /backup
readme.txt
```

#### Ошибки

| Ситуация | Сообщение |
|---|---|
| Мало аргументов | `cp: usage: cp [-r] source target` |
| Источник не найден | `cp: /no/such: No such file or directory` |
| Директория без `-r` | `cp: -r not specified; omitting directory "/home"` |
| Копирование директории в себя | `cp: cannot copy a directory into itself` |

#### Пример скрипта

`testdata/startup.txt` содержит:

```
ls /
cd home
ls
cat readme.txt
cp readme.txt /tmp
cd /tmp
ls
cp -r /home /tmp
ls
cd /
uptime
exit
```

Запуск:

```bash
make run-cp
```

## Лицензия

Учебный проект, свободное использование в образовательных целях.