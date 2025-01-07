package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

const UTF8BOM = "\ufeff"

type DNSRecord struct {
	IP string `json:"ip"`
}

func addBOM() {
	// Добавление BOM в начало вывода
	fmt.Print(UTF8BOM)
}

type Flags struct {
	flagK    bool
	flagU    bool
	flagD    bool
	flagH    bool
	flagC    int
	flagHelp bool
}

var HelpTxt = "[Помощь и информация о программе]\n" +
	"Программа запускает утилиту nslookup, используя указанный при запуске домен " +
	"в виде аргумента программы, по списку публичных DNS серверов с сайта public-dns.info\n" +
	"Использование:\n " +
	"dns-lookup [флаг] [домен]\n" +
	"   или\n dns-lookup [домен]\n" +
	"Если не выбран ни один флаг, по умолчанию программа выполняется по DNS серверам в Hong Kong.\n" +
	"Параметры:\n" +
	" -k 		Использовать DNS серверы в Hong Kong\n" +
	" -u 		Использовать DNS серверы в USA\n" +
	" -d 		Использовать DNS серверы в Germany\n" +
	" -c [число]    	Задает количество IP адресов для выполнения nslookup.\n" +
	"		Считаются результаты без ошибок. Допускается значение от 1 до 1000000.\n" +
	"		Значение 0 означает отсутствие флага.\n" +
	" -h, --help 	Показать \"Помощь и информацию о программе\" и завершить работу"

// Функция для обработки флагов и присвоения значений структуре
func parseFlags() (Flags, string) {
	// Флаги
	flagK := flag.Bool("k", false, "Использовать DNS серверы в Hong Kong")
	flagU := flag.Bool("u", false, "Использовать DNS серверы в USA")
	flagD := flag.Bool("d", false, "Использовать DNS серверы в Germany")
	flagC := flag.Int("c", 0, "Задает количество IP адресов для выполнения nslookup. Допускается значение от 1 до 1000000.\nЗначение 0 означает отсутствие флага.")
	flagH := flag.Bool("h", false, "Помощь и информация")
	flagHelp := flag.Bool("help", false, "Помощь и информация")

	flag.Parse()

	// Если выбраны флаги помощи
	if *flagH || *flagHelp {
		fmt.Println(HelpTxt)
		os.Exit(0)
	}

	// Проверка наличия домена

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Ошибка: Укажите ровно один домен.")
		fmt.Println(HelpTxt)
		os.Exit(1)
	}

	domain := args[0]
	if !validateDomain(domain) {
		fmt.Println("Ошибка: Неверный формат домена.")
		os.Exit(1)
	}

	// Проверка значения флага -c
	if *flagC < 0 || *flagC > 1000000 {
		fmt.Println("Ошибка: Значение флага -c должно быть от 1 до 1000000.")
		os.Exit(1)
	}

	// Обработка флагов
	if !*flagK && !*flagU && !*flagD {
		*flagK = true // По умолчанию Hong Kong
	}
	flags := Flags{
		flagK:    *flagK,
		flagU:    *flagU,
		flagD:    *flagD,
		flagH:    *flagH,
		flagHelp: *flagHelp,
		flagC:    *flagC,
	}
	return flags, domain
}

func lookup(domain string, url string, count int) {
	// Выполнение HTTP-запроса
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Ошибка: %v", err)
		}
	}(resp.Body)

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка: сервер вернул статус %d", resp.StatusCode)
	}

	// Чтение и декодирование JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	var records []DNSRecord
	if err := json.Unmarshal(body, &records); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	// Выполнение команды nslookup для каждого IP
	iteration := 0
	for _, record := range records {
		if record.IP == "" {
			continue
		}
		// Прерываем цикл, если достигнуто заданное количество итераций
		if count > 1 && iteration >= count {
			break
		}
		// Формируем команду nslookup
		cmd := exec.Command("nslookup", domain, record.IP)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Ошибка выполнения nslookup для IP %s: %v\n", record.IP, err)
			continue
		}

		// Вывод результата
		fmt.Printf("Результат nslookup для IP %s:\n%s\n", record.IP, output)
		if count != 0 {
			iteration++
		}
	}
}

func validateDomain(domain string) bool {
	regex := `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, domain)
	return match
}

func main() {
	// Добавляем BOM, чтобы текстовые редакторы корректно распознавали UTF-8
	addBOM()
	flags, domain := parseFlags()

	if flags.flagK {
		fmt.Printf("Выполняется запрос по списку DNS серверов Hong Kong для домена %s\n", domain)
		lookup(domain, "https://public-dns.info/nameserver/hk.json", flags.flagC)
	}
	if flags.flagU {
		fmt.Printf("Выполняется запрос по списку DNS серверов USA для домена %s\n", domain)
		lookup(domain, "https://public-dns.info/nameserver/us.json", flags.flagC)
	}
	if flags.flagD {
		fmt.Printf("Выполняется запрос по списку DNS серверов Germany для домена %s\n", domain)
		lookup(domain, "https://public-dns.info/nameserver/de.json", flags.flagC)
	}
}
