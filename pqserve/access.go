package main

import (
	"github.com/pebbe/util"

	"fmt"
	"net"
	"regexp"
	"strings"
)

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
VIEW:
	for _, a := range Cfg.View {
		if a.all {
			access = a.Allow
			continue VIEW
		}
		for _, ipnet := range a.ipnet {
			if ipnet.Contains(ip) {
				access = a.Allow
				continue VIEW
			}
		}
		for _, aip := range a.ip {
			if aip.Equal(ip) {
				access = a.Allow
				continue VIEW
			}
		}
	}
	return access
}

func accessLogin(mail string) bool {
	if len(Cfg.Access) == 0 {
		return true
	}
	access := true
ACCESS:
	for _, a := range Cfg.Access {
		if a.all {
			access = a.Allow
			continue ACCESS
		}
		for _, re := range a.re {
			if re.MatchString(mail) {
				access = a.Allow
				continue ACCESS
			}
		}
	}
	return access
}
