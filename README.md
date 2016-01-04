# Avito Realty Notifier

*Avito Realty Notifier* создан на основе кодовой базы https://github.com/Lexty/avito-notifier. Люто плюсуем [@Lexty](https://github.com/Lexty) за это :)

*Avito Realty Notifier* — это кастомная Slack-интеграция для уведомления в канал о новых объявлениях недвижимости (продажа/аренда, частная/коммерческая и т.д.).

Для подготовки к запуску потребуется сделать следующее:

1) Получить ссылку с листингом объявлений ([пример](https://www.avito.ru/moskva/kvartiry/sdam/na_dlitelnyy_srok/1-komnatnye?bt=0&pmax=50000)).

2) Добавить и настроить [Incoming WebHook](https://strsqr.slack.com/apps/new/A0F7XDUAZ-incoming-webhooks) интеграцию для канала (любой, хоть general :)).

Далее запускаем скрипт с двумя обязательными параметрами:
- `-s`: URL с результатами поиска недвижимости;
- `-w`: WebHook URL, который сгенерирует Slack для "пушей" сообщений.

`go run main.go -s "https://www.avito.ru/moskva/kvartiry/sdam/na_dlitelnyy_srok/1-komnatnye?bt=0&pmax=50000" -w "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/wowitsrandomguid"`

Или предварительно собираем приложение, а после запускаем бинарник:

`go build main.go`

`./main -s "https://www.avito.ru/moskva/kvartiry/sdam/na_dlitelnyy_srok/1-komnatnye?bt=0&pmax=50000" -w "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/wowitsrandomguid"`

В результате, в вашем Slack-канале вы найдете следующее:

![screen](images/screen.png)

Чтобы далеко не ходить за логотипом Avito для бота, просто возьмите его здесь:

![avito-logo](images/avito-logo.png)

## TODO:

- [ ] Нужно больше информации и аттачменты!
- [ ] Рефакторинг и оптимизация кода
- [ ] Возможность уведомлять не только о недвижимости
