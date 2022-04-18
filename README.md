Установка
```
go build nltools.go
sudo cp ./nltools /usr/local/bin
```

Установка папки с репозиторием dev-server
```
nltools config --dir-dev-server=/path/to/dev-server
```

Команды конфигурации
```
nltools config --dir-dev-server=/path/to/dev-server // Установка папки с репозиторием dev-server
nltools config --show(-s) // просмотр текущей конфигурации
```

Команды для работы с приложением
```
nltools app --show-dev (-d) // покажет сервисы запущенные локально
nltools app --show-stabl (-s) // Показать сервисы не работающие локально
nltools app --show-all (-a) // Показать все сервисы
nltools app --disable=name_service // Переключить сервис в режим stable
nltools app --enable=name_service // Переключить сервис в режим dev
```