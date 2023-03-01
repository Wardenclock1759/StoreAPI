# StoreAPI
В общем, после месяца изучения go я так и не приблизился к истине. По этой причине пришлось собирать код из нескольких источников, а позже, еще и модифицировать.
Что где искать:
  1. internal/storage - интерфейсы и логика бд. Для каждой модели прописан свой репозиторий, через который и осуществвляется связь с бд.
  2. internal/model - собственно сами модели.
  3. internal/apiserver - ядро программы. Там и все роутинги, контроллеры и подобное
  4. migrations - sql скрипты.
  
Что не сделано, а хотелось бы добавить:
  1. Рассылка по почте - пример работает отлично. А вот если заменить тестовую строку, то почему-то высылается пустой email. Протестить и устранить пока не хватило времени.
  2. Тоже с почтой, что нет повторной отправки + нет параллельного процесса под это(
  3. Слишком много мороки с jwt и сессионным куки. В запросах надо указывать то, что дублируется в них.
  4. Так и не смог победить swagger. Даже с их официальными примерами онлайн интерпритатор не позволяет вставить примеры. Так что эта документация по сути бесполезна.
  https://app.swaggerhub.com/apis/Wardenclock1759/of-store_api/1.0.0 Так же можно скопировать из файла сюда https://editor.swagger.io


Небольшой роутинг для теста
Домен: https://gentle-beach-94096.herokuapp.com
  1. /user/sign-up (email, password)
  2. /user/sign-in (email, password) Тут надо взять нужный нам jwt и вставлять в заголовок запроса с названием Token для каждого из следующих запросов (живет 30 минут вроде)
  3. /user/role/grant-role (userid, role) id берем из /private/whoami. В роле пишем "seller"
  4. /store/publisher/game (name, price, user_id) постим игру в систему и приписываем юзера
  5. /store/publisher/key (game_id, code) постим ключ с id игры
  6. /store/payment/session (game_id, email) какую игру купить и email покупателя. Получаем куку (живет 5 минут). Не знаю, важно ли название, но я всегда в заголовке писал Set-Cookie
  7. /store/payment/procedure (card) Номер карты. В заголовке тоже jwt и куки.
  
После этого цикл закончился. На email продавца и покупателя отправиться письмо (правда пустое :biblethump). Пока оставлю сервис включенным до поры до времени. 