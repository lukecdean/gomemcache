package server

import (
    "fmt"
    "net"
    "os"
    "strings"
    "sync"
)

const (
    BUFFER_SIZE = 1024
    C_INVALID = -1
    C_GET = 0
    C_SET = 1
)

var instance    *server
var data        map[string]string
var dataMutex   sync.Mutex

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

    // set server to listening
    server, err := net.Listen(s.serverType, s.serverHost + ":" + s.serverPort)
    if err != nil {
        fmt.Println("err listening: ", err.Error())
        os.Exit(1)
    } // if err
    fmt.Println("server listening on " + s.serverHost + ":" + s.serverPort)

    defer server.Close()
    s.running = true
    defer cleanUp(s)

    // accept clients
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

// TODO break this down with helper funcs
// complete requests from a client
// format
//  set x y
//  get x
func serveClient(connection net.Conn) {
    defer connection.Close()
    for {
        // response var which will eventually be sent to client
        response := "you shouldn't see this"

        // take client's message
        buffer := make([]byte, BUFFER_SIZE)
        mLen, err := connection.Read(buffer)
        if err != nil {
            // if err
            if err.Error() == "EOF" {
                fmt.Println("EOF recieved. Closing connection")
                return
            } else {
                fmt.Println("err reading: ", err.Error())
            } // if
            response = "err reading"
            continue
        } else if mLen == 0 {
            // if the buffer is empty
            fmt.Println("err reading; buffer empty")
            response = "err reading; buffer empty"
            continue
        } // if err reading

        // process buffer
        text := string(buffer)
        // cut the unused characters from the buffer
        text = text[:strings.Index(text, "\n")]
        // tokenize the request
        tokens := strings.Split(text, " ")

        /*
        // debug: print literal chars of tokens
        for _, t := range tokens {
            fmt.Println(strconv.Quote(t))
        } // for
        */

        // switch based off command (first token)
        switch tokens[0] {
        case "get":
            if len(tokens) != 2 {
                // invalid get syntax
                response = "invalid syntax"
                break
            } // if

            dataMutex.Lock()
            value, ok := data[strings.TrimSpace(tokens[1])]
            dataMutex.Unlock()

            if ok {
                response = value
            } else {
                response = "value does not exist"
            }
        case "set":
            if len(tokens) != 3 {
                // invalid set syntax
                response = "invalid syntax"
                break
            } // if

            dataMutex.Lock()
            data[strings.TrimSpace(tokens[1])] = tokens[2]
            response = "value " + tokens[1] + " set to: " + data[strings.TrimSpace(tokens[1])]
            dataMutex.Unlock()
        default:
            response = "format:\n\tset <k> <v>\n\tget <k>"
        } // switch opcode

        // send response to the client
        _, err = connection.Write([]byte(response))
        if err != nil {
            fmt.Println("err writing: ", err.Error())
        } // if
    } // for (recieving loop)
} // serveClient()

// returns (opcode, tokens {command, K, V})
// now defunct
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
            //fmt.Println("finished") // debug
            tokens[2] = value
            break
        }
    } // for

    if extractionFinished {
        return C_SET, tokens
    } else {
        //fmt.Println("extraction failed") // debug
        return C_INVALID, tokens
    } // if
} // parseRequest()

// honestly a func just for testing
func Info(s *server) (string) {
    return s.serverHost
} // info()
