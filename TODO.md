Small stuff:
Make beautiful filler for empty chats
Clear notifications on userlistBtn click?
Handle case when new user was created and sent msg to user with open userlist?
Make total of 20 default users and 10 default posts
Пофиксить нижнюю линию в окне чата
Пофиксить посты, они иногда подгружаются в одном порядке, а иногда в другом


Good practice:
fix routing to load only ones
Make smooth transition on loading old messages with scroll event
Too much auth checks when loading userlist, remove unnessasery
Remove fetch username and creation date, add it to responce struct
fix notification badge dropping chat icon
Handle errors
Clear code from unused code from previous forum


Optional cool fetures:
Можно сделать закладки в окне чата для быстрого открытия переписки
Можно не закрывать чат по нажатию на userlistBtn, оставлять открытым, закрывать на крестик, сворачивать итд
Генерировать кофейный ник
Add animations
Change to dark mode on chrome experimental mode enable
Открытие чата с нужным пользователем при нажатии на username (например автора поста). Возможно предварительно открывать меню - view profile, send message