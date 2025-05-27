package hypr

import (
	"net"
	"screen_tracker/pkg/utils"

	"strings"
)

func StartConnection(path string) net.Conn {
	connection, err := net.Dial("unix", path)
	utils.Check(err)

	return connection

}

func HyprIPC(sock_path string) ([]byte, error) {
	conn := StartConnection(sock_path)

	defer conn.Close()
	/*
		_, err := conn.Write([]byte(cmd))
		if err != nil {
			return nil, fmt.Errorf("Error writting: %v", err)
		}
	*/

	buffer := make([]byte, 4096)
	_, err := conn.Read(buffer)
	utils.Check(err)

	return buffer, nil

}

func GetActiveWindow(buf []byte) string {

	//activewindow >> class,title  :

	bufStr := string(buf)

	res := strings.HasPrefix(bufStr, "activewindow>>")

	if res {
		content := strings.TrimPrefix(bufStr, "activewindow>>")
		part := strings.SplitN(content, ",", 2)

		appClass := part[0]

		return appClass

	}
	return ""

}
