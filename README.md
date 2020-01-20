# Convert-url

В коде прописано обращаться к базе urlstb от пользователя rodp с паролем qwe .
Для этого нужно создать такого пользователя и задать ему пароль
(делать это надо через того пользователя который имеет права создавать новых пользователей бд)
```
createuser --pwprompt rodp

```
после выполнения этой команды он попросит ввести пароль, надо ввести пароль "qwe" 


Создание базы 
```
sudo -u postgres createdb -O rodp urlstb
```
Создание таблицы 

```
CREATE TABLE urlstb(long varchar(255) primary key, short varchar(64), date date);
```
