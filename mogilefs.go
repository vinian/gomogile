package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

type mogileFS struct {
	host   string //better use []string support more than on tracker
	domain string
	class  string
	sock   interface{}
}

func (mog *mogileFS) connect() {
	mogAddr, err := net.ResolveTCPAddr("tcp4", mog.host)
	check_error(err)

	conn, err := net.DialTCP("tcp", nil, mogAddr)
	check_error(err)

	mog.sock = conn
}

func (mog *mogileFS) Cmd(cmd string, args map[string]string) map[string]string {
	mog.connect()
	opts := []string{}
	for key, value := range args {
		str := strings.Join([]string{key, value}, "=")
		opts = append(opts, str)
	}

	opt := strings.Join(opts, "&")
	req := cmd + " " + opt + "\r\n"

	sock := mog.sock.(*net.TCPConn)
	_, err := sock.Write([]byte(req))
	check_error(err)
	var result [512]byte
	_, err = sock.Read(result[0:])
	check_error(err)

	file_info := string_to_hash(string(result[:]))

	return file_info
}

func string_to_hash(str string) map[string]string {
	regex := regexp.MustCompile(`\w+`)
	find := regex.FindString(str)

	req := str[len(find)+1:]

	items := strings.Split(req, "&")
	res := map[string]string{}

	for _, para := range items {
		data := strings.Split(para, "=")
		res[data[0]] = data[1]
	}

	return res
}

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	mog := mogileFS{"host:port", "domain", "class", ""}

	cmd_map := map[string]string{
		"FileInfo":  "file_info",
		"FilePaths": "get_paths",
	}
	request := map[string]string{"key": "mogilekey", "domain": "domain", "device": "0"}
	file_info := mog.Cmd(cmd_map["FileInfo"], request)

	for key, value := range file_info {
		fmt.Printf("%20s : %-30s\n", key, value)
	}
}
