package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func StartServer(args Args) error {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	address := args.Address()
	writer, err := args.Writer()
	if err != nil {
		logger.Printf("args.Writer() %s", err.Error())
		return err
	}

	http.Handle("/", Handler(logger, writer))

	hello := strings.TrimSpace(fmt.Sprintf(`
/---------------------------------------
/ grace server start %s
/---------------------------------------
`, address))

	fmt.Fprintln(os.Stdout, hello)

	return http.ListenAndServe(address, nil)
}
