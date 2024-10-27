# ТЗ на разработку сервиса "Превьювер изображений"

## Общее описание
Сервис предназначен для изготовления preview (создания изображения
с новыми размерами на основе имеющегося изображения).

#### Пример превьюшек в папке [examples](./examples/image-previewer)

## Архитектура
Сервис представляет собой web-сервер (прокси), загружающий изображения,
масштабирующий/обрезающий их до нужного формата и возвращающий пользователю.

## Основной обработчик
http://cut-service.com/fill/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg

<---- микросервис ----><- размеры превью -><--------- URL исходного изображения --------------------------------->

В URL выше мы видим:
- http://cut-service.com/fill/300/200/ - endpoint нашего сервиса,
  в котором 300x200 - это размеры финального изображения.
- https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg -
  адрес исходного изображения; сервис должен скачать его, произвести resize, закэшировать и отдать клиенту.

Сервис должен получить URL исходного изображения, скачать его, изменить до необходимых размеров и вернуть как HTTP-ответ.

- Работаем только с HTTP.
- Ошибки удалённого сервиса или проксируем как есть, или логируем и отвечаем клиенту 502 Bad Gateway.
- Поддержка JPEG является минимальным и достаточным требованием.

**Важно**: необходимо проксировать все заголовки исходного HTTP запроса к целевому сервису (raw.githubusercontent.com в примере).

Сервис должен сохранить (кэшировать) полученное preview на локальном диске и при повторном запросе
отдавать изображение с диска, без запроса к удаленному HTTP-серверу.

Поскольку размер места для кэширования ограничен, то для удаления редко используемых изображений
необходимо использовать алгоритм **"Least Recent Used"**.

## Конфигурация
Основной параметр конфигурации сервиса - разрешенный размер LRU-кэша.

Он может измеряться как количеством закэшированных изображений, так и суммой их байт (на выбор разработчика).

## Развертывание
Развертывание микросервиса должно осуществляться командой `make run` (внутри `docker compose up`)
в директории с проектом.

## Тестирование
Реализацию алгоритма LRU нужно покрыть unit-тестами.

Для интеграционного тестирования можно использовать контейнер с Nginx в качестве удаленного HTTP-сервера,
раздающего вам заданный набор изображений.

Необходимо проверить работу сервера в разных сценариях:
* картинка найдена в кэше;
* удаленный сервер не существует;
* удаленный сервер существует, но изображение не найдено (404 Not Found);
* удаленный сервер существует, но изображение не изображение, а скажем, exe-файл;
* удаленный сервер вернул ошибку;
* удаленный сервер вернул изображение;
* изображение меньше, чем нужный размер;
  и пр.

## Разбалловка
Максимум - **15 баллов**
(при условии выполнения [обязательных требований](./README.md)):

* Реализован HTTP-сервер, проксирующий запросы к удаленному серверу - 2 балла.
* Реализована нарезка изображений - 2 балла.
* Кэширование нарезанных изображений на диске - 1 балл.
* Ограничение кэша одним из способов (LRU кэш) - 1 балл.
* Прокси сервер правильно передает заголовки запроса - 1 балл.
* Написаны интеграционные тесты - 3 балла.
* Тесты адекватны и полностью покрывают функциональность - 1 балл.
* Проект возможно собрать через `make build`, запустить через `make run`
  и протестировать через `make test` - 1 балл.
* Понятность и чистота кода - до 3 баллов.

#### Зачёт от 10 баллов

### Обязательные требования для каждого проекта
* Наличие юнит-тестов на ключевые алгоритмы (core-логику) сервиса.
* Наличие валидных Dockerfile и Makefile/Taskfile для сервиса.
* Ветка master успешно проходит пайплайн в CI-CD системе
  (на ваш вкус, GitHub Actions, Circle CI, Travis CI, Jenkins, GitLab CI и пр.).
  **Пайплайн должен в себе содержать**:
  - запуск последней версии `golangci-lint` на весь проект с
    [конфигом, представленным в данном репозитории](./.golangci.yml);
  - запуск юнит тестов командой вида `go test -race -count 100`;
  - сборку бинаря сервиса для версии Go не ниже 1.22.

При невыполнении хотя бы одного из требований выше - максимальная оценка за проект **4 балла**
(незачёт), несмотря на, например, полностью написанный код сервиса.

Более подробная разбалловка представлена в описании конкретной темы.

### Использование сторонних библиотек для core-логики
Допускается только в следующих темах:
- Анти-брутфорс
- Превьювер изображений

Но:
- вы должны иметь представление о том, что происходит внутри.
- точно ли подходит данная библиотека для решения вашей задачи?
- не станет ли библиотека узким местом сервиса?
- не полезнее ли написать функционал, которые вы хотите получить от библиотеки, самому?

---

Для упрощения проверки вашего репозитория, рекомендуем использовать значки GitHub
([GitHub badges](https://github.com/dwyl/repo-badges)), а также
[Go Report Card](https://goreportcard.com/).

---
Авторы ТЗ:
- [Дмитрий Смаль](https://github.com/mialinx)
- [Антон Телышев](https://github.com/Antonboom)


---
golangci-lint run ./...
WARN The linter 'ifshort' is deprecated (since v1.48.0) due to: The repository of the linter has been deprecated by the owner.  
WARN The linter 'deadcode' is deprecated (since v1.49.0) due to: The owner seems to have abandoned the linter. Replaced by unused.
WARN The linter 'varcheck' is deprecated (since v1.49.0) due to: The owner seems to have abandoned the linter. Replaced by unused.
WARN The linter 'exportloopref' is deprecated (since v1.60.2) due to: Since Go1.22 (loopvar) this linter is no longer relevant. Replaced by copyloopvar.
WARN The linter 'structcheck' is deprecated (since v1.49.0) due to: The owner seems to have abandoned the linter. Replaced by unused.
ERRO [linters_context] deadcode: This linter is fully inactivated: it will not produce any reports.
ERRO [linters_context] ifshort: This linter is fully inactivated: it will not produce any reports.
ERRO [linters_context] structcheck: This linter is fully inactivated: it will not produce any reports.
ERRO [linters_context] varcheck: This linter is fully inactivated: it will not produce any reports. 



internal/config/config.go:8:2: import 'github.com/joho/godotenv' is not allowed from list 'Main' (depguard)
"github.com/joho/godotenv"
^
internal/http/client/client.go:10:2: import 'github.com/randomurban/image-previewer/internal/model' is not allowed from list 'Main' (depguard)
"github.com/randomurban/image-previewer/internal/model"
^
internal/service/service.go:10:2: import 'github.com/disintegration/imaging' is not allowed from list 'Main' (depguard)
"github.com/disintegration/imaging"
^
internal/service/service.go:11:2: import 'github.com/randomurban/image-previewer/internal/http/client' is not allowed from list 'Main' (depguard)
"github.com/randomurban/image-previewer/internal/http/client"
^
internal/http/server/server.go:9:2: import 'github.com/randomurban/image-previewer/internal/service' is not allowed from list 'Main' (depguard)
"github.com/randomurban/image-previewer/internal/service"
^
cmd/previewer/main.go:13:2: import 'github.com/randomurban/image-previewer/internal/config' is not allowed from list 'Main' (depguard)
"github.com/randomurban/image-previewer/internal/config"
^
cmd/previewer/main.go:14:2: import 'github.com/randomurban/image-previewer/internal/http/server' is not allowed from list 'Main' (depguard)
"github.com/randomurban/image-previewer/internal/http/server"
^
