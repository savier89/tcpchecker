package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"time"

	gomail "gopkg.in/mail.v2"
)

var (
	LOG_FILE    string
	EMAIL_USER  string
	EMAIL_PASS  string
	MAIL_TO     string
	SMTP_SERVER string
	SMTP_PORT   int
	SENDMAIL    bool
	ip_server   string
	ip_port     string
)

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

func MakeConnection(server_addr *string, server_port *string) (connection net.Conn, raddr string, laddr string, err error) {
	conn, err := net.Dial("tcp", *server_addr+":"+*server_port)
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
	if err != nil {
		c.Close()
		if SENDMAIL {
			FSendMail("Connection to server - " + c.LocalAddr().String() + " <-> " + c.RemoteAddr().String() + " lost!")
		}

		FLogging("Connection to server - " + c.LocalAddr().String() + " <-> " + c.RemoteAddr().String() + " lost!")
		return false
	}
	return true
}

func main() {

	flag.StringVar(&EMAIL_USER, "from", "user@example.com", "Адрес e-mail с которого слать уведомления.")
	flag.StringVar(&EMAIL_PASS, "password", "password", "Пароль от e-mail адреса.")
	flag.StringVar(&MAIL_TO, "to", "notify@example.com", "Адрес e-mail куда слать уведомления.")
	flag.StringVar(&SMTP_SERVER, "smtp", "smtp.gmail.com", "SMTP Server.")
	flag.IntVar(&SMTP_PORT, "port", 587, "SMTP Server port")
	flag.BoolVar(&SENDMAIL, "sm", false, "Отправлять почту? По умолчанию: false.")

	flag.StringVar(&ip_server, "s", "127.0.0.1", "IP адрес сервера на котором установлен tcp_checker_server.")
	flag.StringVar(&ip_port, "p", "8080", "TCP Порт tcp_ckecker")
	flag.StringVar(&LOG_FILE, "log", "/var/log/tcp_checker_client.log", "Путь до файла куда писать логи.")

	flag.Parse()

	for {
		conn, raddr, laddr, err := MakeConnection(&ip_server, &ip_port)

		if err != nil {
			FLogging("Connection to server - " + ip_server + " refused!")
			time.Sleep(3 * time.Second)
			continue
		}

		if SENDMAIL {
			FSendMail("NEW Connection to server - " + raddr + " from: " + laddr + "!")
		}

		FLogging("NEW Connection to server - " + raddr + " from: " + laddr + "!")

		for isTCPWorking(conn) {
			time.Sleep(1 * time.Second)
		}

		time.Sleep(1 * time.Second)
	}
}
