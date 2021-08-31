package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/ngrok/ngrok-api-go"
)

func main() {
	cmd := exec.Command("ngrok", "http", "-auth=\"jack:sparrow\"", "file:///home/aci1dzero")
	fmt.Println("1")
	closeCmd := func(c *exec.Cmd) {
		if err := c.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
	}
	fmt.Println("B")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("C")
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	fmt.Println("D")
	req, err := http.NewRequest(http.MethodGet, "http://localhost:4040/api/tunnels", nil)
	fmt.Println("E")
	if err != nil {
		panic(err)
	}
	fmt.Println("F")
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("G")
	time.Sleep(5 * time.Second)
	res, err := client.Do(req)
	fmt.Println("H")
	if err != nil {
		panic(err)
	}
	fmt.Println("I")
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var m ngrok.TunnelList
	jsonErr := json.Unmarshal(body, &m)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(m.Tunnels[0].PublicURL)

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	closeCmd(cmd)
}
