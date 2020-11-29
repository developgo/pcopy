package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"pcopy"
)

const (
	defaultSystemConfigFile = "/etc/pcopy/pcopy.conf"
	defaultLocalConfigFile = "$HOME/.config/pcopy/pcopy.conf"
)

func main() {
	if len(os.Args) < 2 {
		printSyntaxAndExit()
	}

	command := os.Args[1]
	switch command {
	case "copy":
		execCopy()
	case "paste":
		execPaste()
	case "serve":
		execServe()
	default:
		printSyntaxAndExit()
	}
}

func execCopy() {
	config, fileId := parseClientArgs("copy")
	client := pcopy.NewClient(config)

	if err := client.Copy(os.Stdin, fileId); err != nil {
		fail(err)
	}
}

func execPaste()  {
	config, fileId := parseClientArgs("paste")
	client := pcopy.NewClient(config)

	if err := client.Paste(os.Stdout, fileId); err != nil {
		fail(err)
	}
}

func parseClientArgs(command string) (*pcopy.Config, string) {
	flags := flag.NewFlagSet(command, flag.ExitOnError)
	serverUrl := flags.String("server", "", "Server URL")
	configFile := flags.String("config", "", "Config file")
	if err := flags.Parse(os.Args[2:]); err != nil {
		fail(err)
	}

	config, err := loadConfig(*configFile)
	if err != nil {
		fail(err)
	}

	if *serverUrl != "" {
		config.ServerUrl = *serverUrl
	}

	if config.ServerUrl == "" {
		fail(errors.New("server address missing, specify -server flag or add 'ServerUrl' to config"))
	}

	fileId := "default"
	if flags.NArg() > 0 {
		fileId = flags.Arg(0)
	}

	return config, fileId
}

func execServe()  {
	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	listenAddr := flags.String("listen", "", "Listen address")
	configFile := flags.String("config", "", "Config file")
	if err := flags.Parse(os.Args[2:]); err != nil {
		fail(err)
	}

	config, err := loadConfig(*configFile)
	if err != nil {
		fail(err)
	}

	if *listenAddr != "" {
		config.ListenAddr = *listenAddr
	}

	if config.ListenAddr == "" {
		fail(errors.New("listen address missing, specify -listen flag or add 'ListenAddr' to config"))
	}

	if err := pcopy.Serve(config); err != nil {
		fail(err)
	}
}

func loadConfig(configFile string) (*pcopy.Config, error) {
	var err error
	var config *pcopy.Config

	if configFile != "" {
		config, err = pcopy.LoadConfig(configFile)
		if err != nil {
			return nil, err
		}
	} else if _, err := os.Stat(os.ExpandEnv(defaultLocalConfigFile)); err == nil {
		config, err = pcopy.LoadConfig(os.ExpandEnv(defaultLocalConfigFile))
		if err != nil {
			return nil, err
		}
	} else if _, err := os.Stat(defaultSystemConfigFile); err == nil {
		config, err = pcopy.LoadConfig(defaultSystemConfigFile)
		if err != nil {
			return nil, err
		}
	} else {
		config = pcopy.DefaultConfig
	}

	return config, nil
}

func printSyntaxAndExit() {
	fmt.Println("Syntax:")
	fmt.Println("  pcopy serve [-listen :1986]")
	fmt.Println("    Start server")
	fmt.Println()
	fmt.Println("  pcopy copy [-config CONFIG] [-server http://...] < myfile.txt")
	fmt.Println("    Copy myfile.txt to the remote clipboard")
	fmt.Println()
	fmt.Println("  pcopy paste [-config CONFIG] [-server http://...] > myfile.txt")
	fmt.Println("    Paste myfile.txt from the remote clipboard")
	os.Exit(1)
}

func fail(err error) {
	fmt.Println(err.Error())
	os.Exit(2)
}