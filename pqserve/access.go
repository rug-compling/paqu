package main

import (
	"github.com/pebbe/util"

	"fmt"
	"net"
	"regexp"
	"strings"
)

func accessSetup() {

	var err error

	for i, view := range Cfg.View {
		if view.Addr == "all" {
			Cfg.View[i].all = true
			continue
		}
		_, ipnet, err := net.ParseCIDR(view.Addr)
		if err == nil {
			Cfg.View[i].ipnet = ipnet
			continue
		}
		ip := net.ParseIP(view.Addr)
		if ip != nil {
			Cfg.View[i].ip = ip
			continue
		}
		util.CheckErr(fmt.Errorf("Ongeldig IP-adres in setup.toml: %s", view.Addr))
	}

	for i, access := range Cfg.Access {
		if access.Mail == "all" {
			Cfg.Access[i].all = true
			continue
		}
		Cfg.Access[i].re, err = regexp.Compile(access.Mail)
		if err != nil {
			util.CheckErr(fmt.Errorf("Ongeldige reguliere expressie als mailadres in setup.toml: %v", err))
		}
	}

}

func accessView(addr string) bool {
	if len(Cfg.View) == 0 {
		return true
	}

	// poortnummer verwijderen
	ad := addr
	i := strings.LastIndex(addr, ":")
	if i > 0 {
		ad = addr[:i]
	}
	// ip parsen
	ip := net.ParseIP(ad)
	if ip == nil {
		logf("Kan IP-adres niet parsen: %v", addr)
		return true
	}

	access := true
	for _, a := range Cfg.View {
		if a.all {
			access = a.Allow
			continue
		}
		if a.ipnet != nil {
			if a.ipnet.Contains(ip) {
				access = a.Allow
			}
			continue
		}
		if a.ip.Equal(ip) {
			access = a.Allow
		}
	}
	return access
}

func accessLogin(mail string) bool {
	if len(Cfg.Access) == 0 {
		return true
	}
	access := true
	for _, a := range Cfg.Access {
		if a.all {
			access = a.Allow
			continue
		}
		if a.re.MatchString(mail) {
			access = a.Allow
		}
	}
	return access
}
