package server

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
)

const (
    BUFFER_SIZE = 1024
    C_INVALID = -1
    C_GET = 0
    C_SET = 1
)

var instance    *server
var data        map[string]string

type server struct {
    serverHost  string // "localhost"
    serverPort  string // "3333"
    serverType  string // "tcp"
    running     bool
}

// singleton pattern
func GetInstance() *server {
    if instance == nil {
        instance = &server{
            serverHost: "localhost",
            serverPort: "3333",
            serverType: "tcp",
            running:    false,
        }
    }
    return instance
}

// boot the server and make it start listening
func Boot(s *server) {
    if s.running {
        fmt.Println("server already running")
        return
    } // if running
    fmt.Println("server booting")

    // create key value map
    data = make(map[string]string)

    server, err := net.Listen(s.serverType, s.serverHost + ":" + s.serverPort)
    if err != nil {
        fmt.Println("err listening: ", err.Error())
        os.Exit(1)
    } // if err
    fmt.Println("server listening on " + s.serverHost + ":" + s.serverPort)

    defer server.Close()
    s.running = true
    defer cleanUp(s)

    // listen for clients
    for {
        connection, err := server.Accept()
        if err != nil {
            fmt.Println("err accepting: ", err.Error())
            os.Exit(1)
        } // if
        fmt.Println("client connected: " + connection.RemoteAddr().String())
        go serveClient(connection)
    } // for
} // Boot()

// set members to non-running values
func cleanUp(s *server) {
    s.running = false
    data = nil
} // cleanUp()

// complete requests from a client
func serveClient(connection net.Conn) {
    defer connection.Close()
    for {
        response := "you shouldn't see this"
        buffer := make([]byte, BUFFER_SIZE)
        mLen, err := connection.Read(buffer)
        if err != nil {
            if err.Error() == "EOF" {
                fmt.Println("EOF recieved. Closing connection")
                os.Exit(0)
            } else {
                fmt.Println("err reading: ", err.Error())
            } // if
            response = "err reading"
            continue
        } else if mLen == 0 {
            // if the buffer is empty
            fmt.Println("err reading")
            response = "err reading"
            continue
        } // if err reading

        text := string(buffer)
        // cut the unused characters from the buffer
        text = text[:strings.Index(text, "\n")]
        //fmt.Println(strconv.Quote(text)) // debug
        //text = strings.Replace(text, "\n", "", -1)

        //fmt.Println(string(buffer)) // debug

        //fmt.Println("recieved: " + string(buffer))//debug
        tokens := strings.Split(text, " ")

        for _, t := range tokens {
            fmt.Println(strconv.Quote(t))
        } // for

        switch tokens[0] {
        case "get":
            if len(tokens) != 2 {
                response = "invalid syntax"
                break
            } // if
            fmt.Println("getting v for k: " + tokens[1]) // debug
            fmt.Printf("size of data: %d\n", len(data)) // debug
            //fmt.Println("keys of data: " + Keys(data)) // debug
            for k := range data {
                fmt.Println(k)
            } // for k
            value, ok := data[strings.TrimSpace(tokens[1])]
            fmt.Printf("value of ok: %t\n", ok) // debug
            if ok {
                fmt.Println("value exists and is: " + value) // debug
                response = "value: " + value
            } else {
                response = "value does not exist"
            }
        case "set":
            if len(tokens) != 3 {
                response = "invalid syntax"
                break
            } // if
            data[strings.TrimSpace(tokens[1])] = tokens[2]
            fmt.Println("value added:" + data[strings.TrimSpace(tokens[1])]) // debug
            response = "value " + tokens[1] + " set to: " + data[strings.TrimSpace(tokens[1])]
        default:
            response = "invalid command"
        } // switch opcode

        fmt.Println("sending to client: " + response) // debug
        _, err = connection.Write([]byte(response))
        if err != nil {
            fmt.Println("err writing: ", err.Error())
        } // if

        /*
        opcode, tokens := parseRequest(buffer)
        switch opcode {
        case C_INVALID:
            response = "err reading: parsing error"
        case C_GET:
            response = "value: "
            response += data[tokens[1]]
        case C_SET:
            data[tokens[1]] = tokens[2]
            response = "value " + tokens[1] + " set to: " + data[tokens[1]]
        default:
            response = "err reading"
        } // switch opcode
        _, err = connection.Write([]byte(response))
        */
    } // for (recieving loop)
} // serveClient()

// returns (opcode, tokens {command, K, V})
func parseRequest(buffer []byte) (int, [3]string) {
    var tokens [3]string
    // parse the request
    rawTokens := strings.Split(string(buffer), " ")

    // if not enough tokens
    if len(rawTokens) < 2 {
        return C_INVALID, tokens
    } // if

    tokens[0] = rawTokens[0] // command
    tokens[1] = rawTokens[1] // K (if get,set)

    // if set, find value (V)
    if (tokens[0] == "get") {
        return C_GET, tokens
    } // if

    // flags for parsing
    openQuoteFound := false
    extractionFinished := false
    var value string = ""
    for index, element := range rawTokens {
        // ignore first two tokens (command, K)
        if index < 2 {
            continue
        } // if

        // search for quote
        if (!openQuoteFound && strings.HasPrefix(element, "\"")) {
            value += element[1:]
            openQuoteFound = true
        } else if (openQuoteFound && strings.HasSuffix(element, "\"")) {
            value += element[:len(element) - 1]
            extractionFinished = true
        } // if

        if extractionFinished {
            fmt.Println("finished") // debug
            tokens[2] = value
            break
        }
    } // for

    if extractionFinished {
        return C_SET, tokens
    } else {
        fmt.Println("extraction failed") // debug
        return C_INVALID, tokens
    } // if
} // parseRequest()

// honestly a func just for testing
func Info(s *server) (string) {
    return s.serverHost
} // info()
