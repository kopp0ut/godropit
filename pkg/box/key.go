package box

const DomKey = `func ChkDom() {
	var hope []byte

	hope = {{.Domain}}
	host, _ := os.Hostname()
	dnsenv := os.Getenv("USERDNSDOMAIN")
	if !strings.Contains(host, string(hope)) && !strings.Contains(dnsenv, string(hope)) {
		time.Sleep({{.Delay}} * time.Seconds)
		os.Exit(0)
	}

}`
