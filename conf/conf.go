package conf

import (
	"log"
	"io"
	"bufio"
	"os"
	"strings"
)

type Config map[string]map[string]string

func (c Config)Init(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("\033[1;31m[error]\033[0m:", err.Error(), "line:")
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var ssip string
	linenum := 1

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln("\033[1;31m[error]\033[0m:", err.Error(), "line:", linenum)
		}

		sline := strings.TrimSpace(string(line))
		if len(sline) == 0 {
			continue
		}

		if strings.Index(sline, "#") == 0 {
			linenum++
			continue
		}

		if strings.Index(sline, ";") == 0 {
			linenum++
			continue
		}

		n := strings.Index(sline, "[")
		m := strings.LastIndex(sline, "]")

		if m > n + 1 {
			ssip = strings.TrimSpace(sline[n + 1 : m])
			c[ssip] = make(map[string]string)
			linenum++
			continue
		}

		if ssip == "" {
			log.Fatalln("\033[1;31m[error]\033[0m: no caption defined. line:", linenum)
		}

		equalsign := strings.Index(sline, "=")
		if equalsign == -1 {
			log.Fatalln("\033[1;31m[error]\033[0m: no = in line. line:", linenum)
		}

		key := strings.TrimSpace(sline[:equalsign])
		if len(key) == 0 {
			log.Fatalln("\033[1;31m[error]\033[0m: no key defined. line:", linenum)
		}

		value := strings.TrimSpace(sline[equalsign + 1:])
		if len(value) == 0 {
			log.Fatalln("\033[1;31m[error]\033[0m: invalid value defined. line:", linenum)
		}

		pos := strings.Index(value, " #")
		if pos != -1 {
			value = value[:pos]
		}

		pos = strings.Index(value, "\t#")
		if pos != -1 {
			value = value[:pos]
		}

		pos = strings.Index(value, "\t;")
		if pos != -1 {
			value = value[:pos]
		}

		pos = strings.Index(value, " ;")
		if pos != -1 {
			value = value[:pos]
		}

		if len(value) == 0 {
			log.Fatalln("\033[1;31m[error]\033[0m: invalid value defined. line:", linenum)
		}

		c[ssip][key] = value
	}
}
