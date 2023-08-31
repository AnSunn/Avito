Запуск через команду ```docker-compose up```

Также необходимо создать файл configs/.env со следующим содержанием:
```
Bind_addr = ":8181"
Logger_level = "debug"

database_url = "host=host port=port user=user password=password dbname=name sslmode=disable"
```
* Используется PostgresSQL, в которой необходимо создать базу данных с таблицами, которые описаны в файлах: ```tablesCreation.sql addDataToUserTbl.sql```

**(!) Заметки:**
* доп задание 1: требовалось на выходе получить csv-файл - не успела доделать. У меня на выходе json-файл со всей необходимой информацией. Если в поле "Operation" значение "Active" - это операция добавления, тогда в поле Date будет дата добавления. Если в поле Operation значение "Removed" - это операция удаления, тогда в поле Date будет дата окончания действия сегмента для пользователя. 

* Таблица Actions (здесь хранится информация о сегментах, назначенных пользователю) хранится start_date и end_date - требуется для логирования и доп заданий 1 и 2. Если добавляют сегмент пользователю на ограниченное кол-во дней, то в поле end_date поставится дата окончания. Если на неограниченный срок добавляют сегмент, то null.
Далее по тексту буду использовать выражение "активный сегмент у пользователя", если end_date больше system date (текущая дата) или null.
Иначе - "Неактивный сегмент у пользователя".
3. При удалении сегмента проставляется признак false в поле "Статус" - это означает, что сегмент неактивный. При этом если удаляемый сегмент был присвоен кому-то из пользователей, то для всех активных сегментов у пользователя в таблице Actions проставляется end_date = дате удаления сегмента из таблицы с сегментами (segments).


**(?) Список вопросов:**
1. Работа с БД - создана схема БД, все манипуляции, которые не предусмотрены через API приложения, делаются напрямую через БД.
2. доп задание 2: требовалось создать возможность добавлять пользователя в сегмент на ограниченный срок.
Возник вопрос: Можно добавлять только на определенное кол-во дней или должна быть возможность еще указать точное время, когда должно быть закончено действие сегмента у пользователя. В примерах везде написано про дни, поэтому можно делать запрос на добавление сегмента пользователю на определенное кол-во дней без точного времени.

**Как использовать Postman requests:**
* GetAllSegments - получить все сегменты. Строку запроса не надо менять
* Get segment by title - получить сегмент по заданному названию сегмента. В строке запроса последний параметр - это название искомого сегмента. 
* Delete segment by title - удалить сегмент по заданному названию сегмента. В строке запроса последний параметр - это название удаляемого сегмента. 
* Create segment - создание нового сегмента. Строку запроса не надо менять. В теле запроса надо указать название добавляемого сегмента. ID присвоится автоматически. В поле Status при добавлении нового сегмента по умолчанию записывается true.
* Create actions - добавление для/удаление у пользователя сегментов. При добавлении сегментов пользователю можно указать срок (кол-во дней), на который добавляется этот сегмент. Строку запроса не надо менять. В теле запроса указывается список добавляемых и удаляемых сегментов. Структура json ниже:
````
{
    "user_id": 4,
    "add_list":
    [
        {
            "title": "AVITO_VOICE_MESSAGES",
            "days": 3
        },
        {
            "title": "AVITO_DISCOUNT_30"
        }
    ],
    "remove_list": ["AVITO_DISCOUNT_90", "AVITO_DISCOUNT_50"]

}
````

* Get Active Segments By User_id - получение активных сегментов пользователя. В строке запроса последний параметр - это id пользователя.
* Get data from actions for a period - получение всех событий по добавлению/удалению сегментов у пользователя за заданный год и месяц. Не успела сделать вывод в виде csv. В строке запроса указывается первым параметром год, вторым месяц. Например: http://localhost:8181/api/v1/data/2023/8, где 2023 - это год, 8 - номер месяца.


