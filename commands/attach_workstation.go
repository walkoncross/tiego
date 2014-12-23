package commands

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/luan/teapot"
	"github.com/luan/tiego/say"
	"github.com/pkg/term"
	"github.com/tmtk75/cli"
)

func AttachWorkstation(c *cli.Context) {
	teapotAddr := c.GlobalString("teapot")
	client := teapot.NewClient(teapotAddr)
	name, _ := c.ArgFor("name")

	ws, err := client.AttachWorkstation(name)
	if err != nil {
		say.Print(0, say.Bold(say.Red("FAILED: ")))
		var errorMessage string
		if len(err.Error()) > 0 {
			fieldRegexp, _ := regexp.Compile("([^:]*: )([^,]*)(,?)")
			errorMessage = fieldRegexp.ReplaceAllString(err.Error(), "$1"+say.Cyan("$2")+"$3")
		} else {
			errorMessage = "Could not talk to Teapot, did you set the " + say.Cyan("TEAPOT") + " url correctly?"
		}
		say.Println(0, errorMessage)
		os.Exit(1)
	}
	defer ws.Close()

	var in io.Reader

	term, err := term.Open(os.Stdin.Name())
	if err == nil {
		err = term.SetRaw()
		if err != nil {
			log.Fatalln("failed to set raw:", term)
		}

		defer term.Restore()

		in = term
	} else {
		in = os.Stdin
	}

	done := make(chan bool)
	go readLoop(ws, os.Stdout, done)
	go writeLoop(ws, in, done)
	<-done
}

func readLoop(c *websocket.Conn, w io.Writer, done chan bool) {
	for {
		_, m, err := c.ReadMessage()
		if err != nil {
			done <- true
			return
		}

		w.Write(m)
	}
}

func writeLoop(c *websocket.Conn, r io.Reader, done chan bool) {
	br := bufio.NewReader(r)
	for {
		x, size, err := br.ReadRune()
		if size <= 0 || err != nil {
			done <- true
			return
		}

		p := make([]byte, size)
		utf8.EncodeRune(p, x)

		err = c.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			done <- true
			return
		}
	}
}
