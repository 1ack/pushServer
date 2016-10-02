package main
 
import (
	"fmt"
    "net/http"
    "os"
    "time"
    "net"
    "bufio"
)
 
type connStatus struct{
 	conn net.Conn
 	name string
 	push_chan chan bool
}

type serverStatus struct{
	num  int
 	conn_pool map[string](*connStatus)

}


var server_status = &serverStatus{num:0,conn_pool:make(map[string](*connStatus))}

func SayHello(w http.ResponseWriter, req *http.Request) {
	for _,conn_status := range server_status.conn_pool{
	    conn_status.conn.Write([]byte("Hello#"))
	}
}
 
func http_listen() {
    http.HandleFunc("/hello", SayHello)
    err := http.ListenAndServe(":8001", nil)
    if err != nil {  
        fmt.Println("ListenAndServe error: ", err.Error())  
    }  
}

func handle_conn(conn_stat *connStatus, name_del chan string){
	conn := conn_stat.conn
	try := 0
	//ticker := time.NewTicker(3 * time.Second)
	data := make([]byte,300)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(time.Second*10)); err != nil {
			return
		}
		if _, err := conn.Read(data); err != nil {
			fmt.Println(err.Error())
			e, ok := err.(net.Error)
			if !ok || !e.Timeout() || try >= 3 {
				name_del <- conn_stat.name
				return 
			}
			try ++
			fmt.Println("try:",try)
			continue
		}
		try = 0 
		//fmt.Println(string(data))
		conn.Write([]byte("heartbeat_pong#"))

	}
}
 
func del_conn( name_timeout <- chan string){

	for{
		select{

			case name := <- name_timeout:
				delete(server_status.conn_pool,name)
				fmt.Printf("delete conn: %s\n",name)

		}
	}



}

func main() {
	go http_listen()
	name_timeout := make(chan string,10)
    go del_conn(name_timeout)
    service := ":8002"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err)
    listener, err := net.ListenTCP("tcp4", tcpAddr)
    checkError(err)
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        name, err := bufio.NewReader(conn).ReadString('#')
        fmt.Printf("client %s connected\n",name)
        server_status.num ++
        fmt.Println(server_status.num)
        conn_stat := &connStatus{conn:conn, name:name, push_chan:make(chan bool)}
        server_status.conn_pool[name] = conn_stat
        go handle_conn(conn_stat, name_timeout )
        for c_name, _:= range server_status.conn_pool{
        	fmt.Println(c_name)
        }

        //go tcp_handle()
        //daytime := time.Now().String()
        //conn.Write([]byte(daytime)) // don't care about return value
        //conn.Close()                // we're finished with this client
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
        os.Exit(1)
    }
}