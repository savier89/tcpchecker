package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gomail "gopkg.in/mail.v2"
)

var (
	LOG_FILE    string
	TG_API_KEY  string
	TG_CHAT_ID  int64
	EMAIL_USER  string
	EMAIL_PASS  string
	MAIL_TO     string
	SMTP_SERVER string
	SMTP_PORT   int
	SENDMAIL    bool
	DEBUG       bool
	IP_SERVER   string
	IP_PORT     int
	BOT         *tgbotapi.BotAPI
)

func FNewTgBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(TG_API_KEY)
	if err != nil {
		log.Printf("TG ERROR: AUTH FAIL! Check your api_key %s", err.Error())
	}

	bot.Debug = DEBUG
	// FLogging("TELEGRAM INFO: Успешная авторизация бота: ")

	if DEBUG {
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}

	return bot
}

func FSendTG(Message string) {
	msg := tgbotapi.NewMessage(TG_CHAT_ID, Message)
	BOT.Send(msg)
}

func FSendMail(Message string) {

	m := gomail.NewMessage()
	m.SetHeader("From", EMAIL_USER)
	m.SetHeader("To", MAIL_TO)
	m.SetHeader("Subject", "TCP Checker NOTIFY: "+Message)
	m.SetBody("text/plain", Message)

	d := gomail.NewDialer(SMTP_SERVER, SMTP_PORT, EMAIL_USER, EMAIL_PASS)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		FLogging(err.Error())
	}
}

func FLogging(Message string) {
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println(Message)
}

func MakeConnection(server_addr *string, server_port *int) (connection net.Conn, raddr string, laddr string, err error) {
	conn, err := net.Dial("tcp", *server_addr+":"+strconv.Itoa(*server_port))
	if err != nil {
		return
	}

	conn.(*net.TCPConn).SetKeepAlive(true)
	conn.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)

	remoteAddr := conn.RemoteAddr().String()
	localAddr := conn.LocalAddr().String()

	return conn, remoteAddr, localAddr, nil
}

func isTCPWorking(c net.Conn) bool {
	_, err := c.Write([]byte("Client send data from local address: " + c.LocalAddr().String()))

	if DEBUG {
		FLogging("Client send data from local address: " + c.LocalAddr().String())
		FSendTG("Client send data from local address: " + c.LocalAddr().String())
	}

	if err != nil {
		c.Close()
		message := "Connection to server - " + c.LocalAddr().String() + " <-> " + c.RemoteAddr().String() + " lost!"
		if SENDMAIL {
			FSendMail(message)
		}

		FLogging(message)
		FSendTG(message)
		return false
	}
	return true
}

func main() {

	// FOR TG MESSAGES
	flag.StringVar(&TG_API_KEY, "api_key", "Change TG API KET HERE", "Telegram BOT API KEY")
	flag.Int64Var(&TG_CHAT_ID, "chat_id", "Change TG Group id here", "ID группы где постить уведомления.")
	flag.BoolVar(&DEBUG, "debug", false, "Включить Debug? По умолчанию: false.")

	// BOT Init
	BOT = FNewTgBot()

	// FOR MAIL SEND MESSAGE
	flag.StringVar(&EMAIL_USER, "from", "user@example.com", "Адрес e-mail с которого слать уведомления.")
	flag.StringVar(&EMAIL_PASS, "password", "password", "Пароль от e-mail адреса.")
	flag.StringVar(&MAIL_TO, "to", "notify@example.com", "Адрес e-mail куда слать уведомления.")
	flag.StringVar(&SMTP_SERVER, "smtp", "smtp.gmail.com", "SMTP Server.")
	flag.IntVar(&SMTP_PORT, "smtp_port", 587, "SMTP Server port")
	flag.BoolVar(&SENDMAIL, "sm", false, "Отправлять почту? По умолчанию: false.")

	// FOR TCP CHECK
	flag.StringVar(&IP_SERVER, "s", "127.0.0.1", "IP адрес сервера на котором установлен tcp_checker_server.")
	flag.IntVar(&IP_PORT, "p", 8080, "TCP Порт tcp_ckecker")
	flag.StringVar(&LOG_FILE, "log", "./logs/tcp_checker_client.log", "Путь до файла куда писать логи.")

	flag.Parse()

	for {
		conn, raddr, laddr, err := MakeConnection(&IP_SERVER, &IP_PORT)

		if err != nil {
			FLogging("Connection to server - " + IP_SERVER + " refused!")
			time.Sleep(3 * time.Second)
			continue
		}

		message := "NEW Connection to server - " + raddr + " from: " + laddr + "!"

		// Send MAIL
		if SENDMAIL {
			FSendMail(message)
		}

		// Send Message to Telegram
		FSendTG(message)

		// Logging Message
		FLogging(message)

		for isTCPWorking(conn) {
			time.Sleep(1 * time.Second)
		}

		time.Sleep(1 * time.Second)
	}
}
