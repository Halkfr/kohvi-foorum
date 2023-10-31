TODO: 

Small stuff:
Fetch real count of post replies or delete them
Clear notifications on userlistBtn click?
Handle case when new user was created and sent msg to user with open userlist?
Make total of 20 default users and 10 default posts
Cделать редирект с localhost на 127.0.0.1:8080/
Пофиксить нижнюю линию в окне чата
Пофиксить посты, они иногда подгружаются в одном порядке, а иногда в другом
Писать время создания поста на странице поста тоже, а не только в home
Хэндлить случай, когда у поста нет комментариев
Убирать Post Image, если у поста нет изображения


Good practice:
fix routing to load only ones
Make smooth transition on loading old messages with scroll event
Too much auth checks when loading userlist, remove unnessasery
Remove fetch username and creation date, add it to responce struct
fix notification badge dropping chat icon
Handle errors
Clear code from unused code from previous forum


Optional cool fetures:
Генерировать кофейный ник
Add animations
Change to dark mode on chrome experimental mode enable
Открытие чата с нужным пользователем при нажатии на username (например автора поста). Возможно предварительно открывать меню - view profile, send message
