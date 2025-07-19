# tcpchecker

### How-to build

    # Get sources
    $ git clone https://github.com/savier89/tcpchecker.git
    $ cd tcpchecker
    $ go mod download
    
    # Build server APP
    $ go build -o tcp_checker_server server/server.go
    
    # Build client APP
    go build -o tcp_checker client/client.go

### Change Telegramm Bot id and Group id in client/client.go file
    flag.StringVar(&TG_API_KEY, "api_key", "Change TG API KET here", "Telegram BOT API KEY")
	flag.Int64Var(&TG_CHAT_ID, "chat_id", "Change TG Group id here", "ID группы где постить уведомления.")

### Server usage commands

    $ ./tcp_checker_server -h 
    Usage of ./tcp_checker_server:
    -p string
        Server TCP port for listen. (default "8080")
    -s string
        Server ip address or hostname for listen. (default "0.0.0.0")

### Client usage commands

    $ ./tcp_checker -h 
    Usage of ./tcp_checker:
      -from string
            Адрес e-mail с которого слать уведомления. (default "user@example.com")
      -log string
            Путь до файла куда писать логи. (default "/var/log/tcp_checker_client.log")
      -p string
            TCP Порт tcp_ckecker (default "8080")
      -password string
            Пароль от e-mail адреса. (default "password")
      -port int
            SMTP Server port (default 587)
      -s string
            IP адрес сервера на котором установлен tcp_checker_server. (default "127.0.0.1")
      -sm
            Отправлять почту? По умолчанию: false.
      -smtp string
            SMTP Server. (default "smtp.gmail.com")
      -to string
            Адрес e-mail куда слать уведомления. (default "notify@example.com")
