package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ExtractClipboard extracts the name of the clipboard from the config filename, e.g. the name of a clipboard with
// the config file /etc/pcopy/work.conf is "work".
func ExtractClipboard(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), suffixConf)
}

// ExpandServerAddr expands the server address with the default port if no port is provided to a full URL, including
// protocol prefix. For instance: "myhost" will become "https://myhost:2586", and "myhost:443" will become "https://myhost",
// but "http://myhost:1234" will remain unchanged.
func ExpandServerAddr(serverAddr string) string {
	if strings.HasPrefix(serverAddr, "http://") || strings.HasPrefix(serverAddr, "https://") {
		return serverAddr
	}
	if !strings.Contains(serverAddr, ":") {
		serverAddr = fmt.Sprintf("%s:%d", serverAddr, DefaultPort)
	}
	return fmt.Sprintf("https://%s", strings.ReplaceAll(serverAddr, ":443", ""))
}

// CollapseServerAddr removes the default port from the given server address if the address contains
// the default port, but leaves the address unchanged if it doesn't contain it.
func CollapseServerAddr(serverAddr string) string {
	if strings.HasPrefix(serverAddr, "http://") {
		return serverAddr
	}
	if strings.HasPrefix(serverAddr, "https://") {
		u, err := url.Parse(serverAddr)
		if err != nil {
			return serverAddr
		}
		return strings.TrimSuffix(u.Host, fmt.Sprintf(":%d", DefaultPort))
	}
	return strings.TrimSuffix(serverAddr, fmt.Sprintf(":%d", DefaultPort))
}

// DefaultCertFile returns the default path to the certificate file, relative to the config file. If mustExist is
// true, the function returns an empty string if the file does not exist.
func DefaultCertFile(configFile string, mustExist bool) string {
	return defaultFileWithNewExt(suffixCert, configFile, mustExist)
}

// DefaultKeyFile returns the default path to the key file, relative to the config file. If mustExist is
// true, the function returns an empty string.
func DefaultKeyFile(configFile string, mustExist bool) string {
	return defaultFileWithNewExt(suffixKey, configFile, mustExist)
}

func defaultFileWithNewExt(newExtension string, configFile string, mustExist bool) string {
	file := strings.TrimSuffix(configFile, suffixConf) + newExtension
	if mustExist {
		if _, err := os.Stat(file); err != nil {
			return ""
		}
	}

	return file
}

// TODO move this to util package
func parseSize(s string) (int64, error) {
	matches := sizeStrRegex.FindStringSubmatch(s)
	if matches == nil {
		return -1, fmt.Errorf("invalid size %s", s)
	}
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, fmt.Errorf("cannot convert number %s", matches[1])
	}
	switch strings.ToUpper(matches[2]) {
	case "G":
		return int64(value) * 1024 * 1024 * 1024, nil
	case "M":
		return int64(value) * 1024 * 1024, nil
	case "K":
		return int64(value) * 1024, nil
	default:
		return int64(value), nil
	}
}
