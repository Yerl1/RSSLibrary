package loadenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnv читает .env файл и загружает пары KEY=VALUE в окружение процесса.
// Если файл не найден — не считается ошибкой.
func LoadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		// если файла нет, просто игнорируем
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// поддержка export VAR=VALUE
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(line[len("export "):])
		}

		// ищем "="
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "ignoring malformed line %d: %q\n", lineNo, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// уберём кавычки, если есть
		if len(val) >= 2 {
			if val[0] == '"' && val[len(val)-1] == '"' {
				val = strings.ReplaceAll(val[1:len(val)-1], `\"`, `"`)
			} else if val[0] == '\'' && val[len(val)-1] == '\'' {
				val = val[1 : len(val)-1]
			}
		}

		// устанавливаем в окружение
		_ = os.Setenv(key, val)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
