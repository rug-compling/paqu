package main

import (
	"github.com/pebbe/util"

	"fmt"
	"net"
	"regexp"
	"strings"
)

// Voorbereiding voor access-functies, aangeroepen vanuit main()
func accessSetup() {

	for i, view := range Cfg.View {
		Cfg.View[i].ip = make([]net.IP, 0, len(view.Addr))
		Cfg.View[i].ipnet = make([]*net.IPNet, 0, len(view.Addr))
		for _, addr := range view.Addr {
			if addr == "all" {
				Cfg.View[i].all = true
				break
			}
			_, ipnet, err := net.ParseCIDR(addr)
			if err == nil {
				Cfg.View[i].ipnet = append(Cfg.View[i].ipnet, ipnet)
				continue
			}
			ip := net.ParseIP(addr)
			if ip != nil {
				Cfg.View[i].ip = append(Cfg.View[i].ip, ip)
				continue
			}
			util.CheckErr(fmt.Errorf("Ongeldig IP-adres in setup.toml: %s", addr))
		}
	}

	for i, access := range Cfg.Access {
		for _, mail := range access.Mail {
			if mail == "all" {
				Cfg.Access[i].all = true
				break
			}
			re, err := regexp.Compile(mail)
			if err != nil {
				util.CheckErr(fmt.Errorf("Ongeldige reguliere expressie als mailadres in setup.toml: %v", err))
			}
			Cfg.Access[i].re = append(Cfg.Access[i].re, re)
		}
	}
}

// Toegang op basis van ip-adress: algemene toegang tot de site
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

	// haken weghalen bij IPv6
	if l := len(ad); l > 1 {
		if ad[0] == '[' && ad[l-1] == ']' {
			ad = ad[1 : l-1]
		}
	}

	// ip parsen
	ip := net.ParseIP(ad)
	if ip == nil {
		logf("Kan IP-adres niet parsen: %v", addr)
		return true
	}

	for _, a := range Cfg.View {
		if a.all {
			return a.Allow
		}
		for _, ipnet := range a.ipnet {
			if ipnet.Contains(ip) {
				return a.Allow
			}
		}
		for _, aip := range a.ip {
			if aip.Equal(ip) {
				return a.Allow
			}
		}
	}
	return true
}

// Toegang op basis van e-mailadress: wie mag inloggen
func accessLogin(mail string) bool {
	if len(Cfg.Access) == 0 {
		return true
	}
	for _, a := range Cfg.Access {
		if a.all {
			return a.Allow
		}
		for _, re := range a.re {
			if re.MatchString(mail) {
				return a.Allow
			}
		}
	}
	return true
}
