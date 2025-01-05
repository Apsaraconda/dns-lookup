# dns-lookup
The program runs the **nslookup** utility, using the domain specified at startup as an argument to the program, against the list of public DNS servers from the site [public-dns.info](https://public-dns.info)

**Launch:**

Launch the program from the terminal. The parameters are passed as command line arguments. For example, in Ubuntu:
```
./dns-lookup [flag] [domain]
```
or
```
./dns-lookup [domain]
```
If no flag is selected, the program is executed by default against DNS servers in Hong Kong.

*Note:* On Windows, the terminal may use CP866 encoding by default, so before launching it is recommended to set the terminal encoding to UTF-8 by entering the command:
```
chcp 65001
```
**Parameters (flags):**

| Flag | Description |
| --- | --- |
| -k | Use DNS servers in Hong Kong |
| -u | Use DNS servers in the USA |
| -d | Use DNS servers in Germany |
| -c [number] | Specifies the number of IP addresses to run **nslookup** on. Only results without errors are counted. The allowed value is from 1 to 1000000. The value 0 means no flag. |
| -h, --help | Show "Help and information about the program" and exit |

**Usage example:**
```
./dns-lookup -d -k -c 5 example.com
```
The program will run **nslookup** on five IP addresses of DNS servers from the Hong Kong list, then Germany, counting from the beginning of each list. The established run order is: Hong Kong, USA, Germany.

**Compilation requirements:**

golang version 1.18 and higher. More information about installation: [go.dev/doc/install](https://go.dev/doc/install)

**Compilation:**
```
go build dns-lookup.go
```
