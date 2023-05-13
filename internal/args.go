package internal

import (
	"fmt"
	"os"
)

type Output string

const (
	OutputStdout = Output("stdout")
	OutputFile   = Output("file")
)

type Args struct {
	Port   int    `short:"p" long:"port" default:"8080" env:"PORT" description:"listen port"`
	Host   string `short:"h" long:"host" default:"127.0.0.1" description:"listen address"`
	Output Output `short:"o" long:"output" default:"stdout" choice:"stdout" choice:"file" description:"output format"`
	Dir    string `short:"d" long:"dir" default:"./" default-mask:"current directory" description:"file output directory"`
}

func (a Args) Address() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a Args) Writer() (Writer, error) {
	if a.Output == OutputStdout {
		return &stdWriter{wc: os.Stdout}, nil
	}
	stat, err := os.Stat(a.Dir)
	if err != nil {
		return nil, err
	}

	return &fileWriter{dir: stat}, nil
}
