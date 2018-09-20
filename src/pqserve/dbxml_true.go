// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	has_dbxml = true
)

func get_dact(archive, filename string) ([]byte, error) {
	reader, err := dbxml.OpenRead(archive)
	if err != nil {
		return []byte{}, err
	}
	defer reader.Close()
	d, err := reader.Get(filename)
	if err != nil {
		return []byte{}, err
	}
	return []byte(d), nil
}

func makeDact(dact, conllu, xml string, stripchar string, chKill chan bool) error {
	files, err := filenames2(xml, false)
	if err != nil {
		return err
	}

	os.Remove(dact)
	db, err := dbxml.OpenReadWrite(dact)
	if err != nil {
		return err
	}
	defer db.Close()

	var dbx *dbxml.Db
	if Cfg.Dactx {
		os.Remove(dact + "x")
		dbx, err = dbxml.OpenReadWrite(dact + "x")
		if err != nil {
			return err
		}
		defer dbx.Close()
	}

	for _, name := range files {

		select {
		case <-chGlobalExit:
			return errGlobalExit
		case <-chKill:
			return errKilled
		default:
		}

		data, err := ioutil.ReadFile(filepath.Join(xml, name))
		if err != nil {
			return err
		}

		name = decode_filename(name)
		if stripchar != "" {
			name = name[1+strings.Index(name, stripchar):]
		}
		err = db.PutXml(name, string(data), false)
		if err != nil {
			return err
		}

		if Cfg.Dactx {
			content, err := dactExpand(data)
			if err != nil {
				return err
			}
			if content == "" {
				content = string(data)
			}
			err = dbx.PutXml(name, content, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func unpackDact(data, xmldir, dact, conllu, stderr string, chKill chan bool) (tokens, nline int, err error) {

	if Cfg.Conllu {
		cmd := shell(
			`pqudep -p %s/data/ -o %s > /dev/null 2> %s.err`,
			paqudatadir, dact, conllu)
		err = run(cmd, chKill, nil)
		if err != nil {
			return 0, 0, err
		}
		if cu, _ := os.Stat(conllu + ".err"); cu.Size() != 0 {
			sysErr(fmt.Errorf("CONLLU error(s) in %s.err", conllu))
		}
	}

	os.Mkdir(xmldir, 0777)

	dc, err := dbxml.OpenRead(dact)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"Openen mislukt. PaQu kan geen dact-bestanden lezen die gemaakt zijn met DbXml nieuwer dan versie %d",
			dbxml_version_major())
	}
	defer dc.Close()

	var dbx *dbxml.Db
	if Cfg.Dactx {
		os.Remove(dact + "x")
		dbx, err = dbxml.OpenReadWrite(dact + "x")
		if err != nil {
			return 0, 0, err
		}
		defer dbx.Close()
	}

	docs, err := dc.All()
	if err != nil {
		return 0, 0, fmt.Errorf("Lezen dact-bestand: %s", err)
	}
	defer docs.Close()

	fperr, err := os.Create(stderr)
	if err != nil {
		return 0, 0, err
	}
	defer fperr.Close()

	fplines, err := os.Create(data + ".lines")
	if err != nil {
		return 0, 0, err
	}
	defer fplines.Close()

	tokens = 0
	nline = 0
	nd := -1
	sdir := ""
	for docs.Next() {

		select {
		case <-chGlobalExit:
			return 0, 0, errGlobalExit
		case <-chKill:
			return 0, 0, errKilled
		default:
		}

		nline++

		if nline%10000 == 1 {
			nd++
			sdir = fmt.Sprintf("%04d", nd)
			os.Mkdir(filepath.Join(xmldir, sdir), 0777)
		}

		name := docs.Name()
		if strings.HasSuffix(strings.ToLower(name), ".xml") {
			name = name[:len(name)-4]
		}
		encname := encode_filename(name)

		data := []byte(docs.Content())

		fp, err := os.Create(filepath.Join(xmldir, sdir, encname+".xml"))
		if err != nil {
			return 0, 0, err
		}
		_, err = fp.Write(data)
		fp.Close()
		if err != nil {
			return 0, 0, err
		}

		alpino := Alpino_ds_no_node{}
		err = xml.Unmarshal(data, &alpino)
		if err != nil {
			return 0, 0, fmt.Errorf("Parsen van %q uit dact-bestand: %s", docs.Name(), err)
		}
		tokens += len(strings.Fields(alpino.Sentence))
		fmt.Fprintf(fplines, "%s|%s\n", name, strings.TrimSpace(alpino.Sentence))
		if alpino.Comments != nil {
			for _, c := range alpino.Comments {
				if strings.HasPrefix(c, "Q#") {
					a := strings.SplitN(c, "|", 2)
					if len(a) == 2 {
						fmt.Fprintf(fperr, "Q#%s|%s\n", name, strings.TrimSpace(a[1]))
						break
					}
				}
			}
		}

		if Cfg.Dactx {
			content, err := dactExpand(data)
			if err != nil {
				return 0, 0, err
			}
			if content == "" {
				content = string(data)
			}
			err = dbx.PutXml(docs.Name(), content, false)
			if err != nil {
				return 0, 0, err
			}
		}
	}

	return tokens, nline, nil
}

func dbxml_version() string {
	x, y, z := dbxml.Version()
	return fmt.Sprintf("Version %d.%d.%d", x, y, z)
}

func dbxml_version_major() int {
	x, _, _ := dbxml.Version()
	return x
}

func dactExpand(data []byte) (string, error) {
	alpino := Alpino_ds_complete{}
	err := xml.Unmarshal(data, &alpino)
	if err != nil {
		return "", err
	}

	refs := make(map[string]*Node)
	getIndexed(alpino.Node0, refs)
	if len(refs) == 0 {
		return "", nil
	}
	err = expandNode(alpino.Node0, refs)
	if err != nil {
		return "", err
	}

	alpino.Version = "X-" + alpino.Version

	return format(alpino)
}

func getIndexed(node *Node, nodes map[string]*Node) {
	if node.Index != "" && (node.NodeList != nil || node.Word != "") {
		nodes[node.Index] = node
	}
	if node.NodeList != nil {
		for _, n := range node.NodeList {
			getIndexed(n, nodes)
		}
	}
}

func expandNode(n *Node, nodes map[string]*Node) error {
	if n.NodeList != nil {
		for _, node := range n.NodeList {
			err := expandNode(node, nodes)
			if err != nil {
				return err
			}
		}
	}

	if n.Index == "" || n.NodeList != nil || n.Word != "" {
		return nil
	}

	o, ok := nodes[n.Index]
	if !ok {
		return fmt.Errorf("Expanding Dact: Missing index node")
	}

	n.OtherId = o.Id

	copyNodeOnEmpty(n, o)

	return nil
}
