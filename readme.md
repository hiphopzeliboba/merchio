## merchio 

### Запуск сервиса

1. `git clone https://github.com/hiphopzeliboba/merchio.git `
2. `cd merchio`
3. `docker-compose up -d --build`

### Что удалось реализовать:
1. проект запускается локально (сам го-сервис не может подключиться к postgres по conn_url, где-то я накосячил но могу найти где, в Goland подключается из контенера нет)
2. есть что-то похожее на выдачу JWT токена и провера на валидность
3. есть покрытие unit тестами (51,9% если смотреть в goland на каверидж по директории merchio), но большинство не проходят из-за ошибок логики кода)))
4. из ручек, только api/auth


Спасибо за возможность поучаствовать и за ваше время)

