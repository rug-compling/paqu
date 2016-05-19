/*

Dit is een voorbeeldprogramma dat laat zien hoe je via een externe applicatie
kunt inloggen in PaQu. Dit wordt aangeroepen als cgi-script zonder parameters
vanuit een reguliere webserver.

Om dit te gebruiken in plaats van het in PaQu ingebouwde mechanisme om in te loggen,
definieer 'loginurl' in de setup van PaQu.

*/

package main

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"

	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Config struct {
	Url    string
	Login  string
	Prefix string
	Maxwrd int
	Access []AccessType
}

type AccessType struct {
	Allow bool
	Mail  []string
	all   bool
	re    []*regexp.Regexp
}

var (
	Cfg Config
)

func main() {

	////////////////////////////////////////////////////////////////
	// AANPASSEN

	// Het complete path naar het setup-bestand van PaQu
	setup := "/path/naar/.paqu/setup.toml"

	// Plaats hier code voor het inloggen.
	// Als uitkomst moet het e-mailadres van de gebruiker
	// in de variabele 'mail' komen.
	mail := "voorbeeld@voorbeeld.nl"

	// Laat de rest van de code zoals het is.
	////////////////////////////////////////////////////////////////

	// Het mailadres moet in kleine letters zijn.
	// Wanneer er meerdere adressen zijn wordt het adres gebruikt dat
	// alfabetisch eerst komt. Zo wordt verzekerd dat altijd hetzelfde adres
	// gebruikt wordt, ongeacht in welke volgorde de adressen worden ontvangen.
	mails := strings.Fields(strings.ToLower(mail))
	if len(mails) < 1 {
		x(errors.New("Missing e-mail"))
		return
	}
	sort.Strings(mails)
	mail = mails[0]

	// Inlezen van de configuratie van PaQu
	_, err := toml.DecodeFile(setup, &Cfg)
	if x(err) {
		return
	}

	// Mag deze gebruiker inloggen?
	if !accessLogin(mail) {
		x(errors.New("Verboden toegang"))
		return
	}

	// Twee strings die gebruikt worden bij het inloggen in PaQu
	auth := rand16()
	sec := rand16()

	// Verbinding maken met MySQL
	db, err := sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
	if x(err) {
		return
	}
	defer db.Close()

	// Kijk of er al een record aanwezig is bij dit mailadres.
	rows, err := db.Query(fmt.Sprintf("SELECT * from `%s_users` WHERE `mail` = %q", Cfg.Prefix, mail))
	if x(err) {
		return
	}
	if rows.Next() {
		rows.Close()
		// Als het mailadres bekend is, alleen de waardes van 'sec' en 'pw' bijwerken.
		_, err = db.Exec(fmt.Sprintf(
			"UPDATE `%s_users` SET `sec` = %q, `pw` = %q WHERE `mail` = %q",
			Cfg.Prefix, sec, auth, mail))
	} else {
		// Als het mailadres niet bekend is, een nieuw record aanmaken, met 'sec' en 'pw', en de standaardwaarde voor 'quotum'.
		_, err = db.Exec(fmt.Sprintf(
			"INSERT INTO `%s_users` (`mail`, `sec`, `pw`, `quotum`) VALUES (%q, %q, %q, %d)",
			Cfg.Prefix, mail, sec, auth, Cfg.Maxwrd))
	}
	if x(err) {
		return
	}

	// Redirect de browser terug naar de pagina van PaQu die de login verder afhandelt.
	fmt.Printf("Location: %s?pw=%s\n\n", urlJoin(Cfg.Url, "login"), urlencode(auth))
}

func rand16() string {
	a := make([]byte, 16)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		a[i] = byte(97 + rnd.Intn(26))
	}
	return string(a)
}

func urlencode(s string) string {
	var buf bytes.Buffer
	for _, b := range []byte(s) {
		if b >= 'a' && b <= 'z' ||
			b >= 'A' && b <= 'Z' ||
			b >= '0' && b <= '9' {
			buf.WriteByte(b)
		} else {
			buf.WriteString(fmt.Sprintf("%%%02x", b))
		}
	}
	return buf.String()
}

func urlJoin(elem ...string) string {
	p := elem[0]
	for _, e := range elem[1:] {
		if strings.HasSuffix(p, "/") {
			if strings.HasPrefix(e, "/") {
				p += e[1:]
			} else {
				p += e
			}
		} else {
			if strings.HasPrefix(e, "/") {
				p += e
			} else {
				p += "/" + e
			}
		}
	}
	return p
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

func x(e error) bool {
	if e == nil {
		return false
	}
	fmt.Printf(`Content-type: text/plain
Cache-Control: no-cache
Pragma: no-cache

Error: %v
`, e)
	return true
}
