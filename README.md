# Проект "Поиск пути на Википедии"

## Описание проекта
Этот проект выполняет поиск пути от одной страницы Википедии до другой, используя алгоритм поиска в ширину (BFS). Программа парсит страницы Википедии, извлекает ссылки и находит кратчайший путь между двумя указанными страницами. Результаты поиска и текст абзацев, содержащих найденные ссылки, выводятся в консоль. Также ведется лог посещенных страниц.

## Основные компоненты

### Основной файл
Основной файл программы содержит функцию `main`, которая выполняет следующие шаги:
1. Считывает начальную и целевую URL-адреса с клавиатуры.
2. Запускает поиск пути с помощью функции `bfs`.
3. Если путь найден, извлекает текст абзацев, содержащих ссылки, и выводит результаты в консоль.

### Функции
- `fetchLinks(pageURL string) ([]string, error)`: Извлекает все ссылки со страницы по указанному URL.
- `UpgradeURL(links *[]string) []string`: Преобразует относительные URL в абсолютные.
- `linkFilter(links []string) []string`: Фильтрует ссылки, оставляя только те, которые ведут на статьи Википедии.
- `bfs(startURL, targetURL string, traceLimit int) bool`: Выполняет поиск в ширину, чтобы найти путь от начальной страницы до целевой.
- `readableURL(encodedURL *string) (string, error)`: Декодирует URL, чтобы сделать его читаемым.
- `fetchParagraphWithLink(currentURL, targetURL string) (string, error)`: Извлекает текст абзаца, содержащего ссылку на целевую страницу.
- `readURLFromInput(prompt string) string`: Считывает URL с клавиатуры.
- `initLogger() `: Настраивает логгер для записи посещенных страниц.

## Как использовать

1. Запустите программу.
2. Введите начальный URL страницы Википедии, с которой нужно начать поиск.
3. Введите целевой URL страницы Википедии, до которой нужно найти путь.
4. Программа выполнит поиск и выведет результаты в консоль.
5. Лог посещенных страниц будет записан в файл `visited_pages.log`.