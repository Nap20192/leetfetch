# leetfetch

`leetfetch` — CLI-утилита, которая забирает задачу с LeetCode по slug/URL и создает готовую папку-шаблон для решения на Go:
- `README.md` с условием задачи (HTML → Markdown)
- `solution.go` со сниппетом из LeetCode
- `solution_test.go` с автосгенерированными тест-кейсами из примеров
- `go.mod` для отдельного модуля задачи

## Быстрый старт

```bash
git clone https://github.com/Nap20192/leetfetch.git
cd leetfetch
go run . two-sum
```

После этого появится директория `0001-two-sum/` с файлами для старта решения.

## Использование

```bash
leetfetch [flags] <slug|url>
```

Примеры:

```bash
# По slug
leetfetch two-sum

# По URL
leetfetch https://leetcode.com/problems/two-sum/
leetfetch https://leetcode.com/problems/two-sum/description/

# Кастомная директория вывода
leetfetch -o ~/leetcode two-sum

# Перезаписать уже существующую директорию
leetfetch -f two-sum

# Взять сниппет в другом языке (если есть в API LeetCode)
leetfetch -lang python3 add-two-numbers
```

## Флаги

| Флаг | По умолчанию | Описание |
| --- | --- | --- |
| `-o` | `.` | Куда создавать папку задачи |
| `-f` | `false` | Перезаписать существующую папку задачи |
| `-lang` | `golang` | Язык сниппета из LeetCode (`golang`, `python3`, `cpp`, `java`, ...) |

## Что генерируется

Для `two-sum`:

```text
0001-two-sum/
├── README.md
├── solution.go
├── solution_test.go
└── go.mod
```

`solution_test.go` генерируется по примерам из LeetCode, но ожидаемые значения (`want`) помечаются как `TODO` и их нужно заполнить вручную.

## Ограничения

- Платные задачи (`isPaidOnly`) не скачиваются.
- Если в выбранном `-lang` нет сниппета, в `solution.go` будет заглушка с `TODO`.
- Для `ListNode` / `TreeNode` добавляются подсказки, но сами типы нужно определить вручную.

## Разработка

```bash
go test ./...
go run . two-sum
```

![Go Report Card](https://goreportcard.com/badge/github.com/Nap20192/leetfetch)
