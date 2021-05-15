package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"inet.af/netaddr"
)

func readFile(name string) []netaddr.IPPrefix {
	var prefixes []netaddr.IPPrefix
	var f *os.File
	if name == "-" {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, "/") {
			if strings.Contains(line, ":") {
				line = line + "/128"
			} else {
				line = line + "/32"
			}
		}
		prefix, err := netaddr.ParseIPPrefix(line)
		if err != nil {
			log.Fatal(err)
		}
		prefixes = append(prefixes, prefix)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return prefixes
}

func filterIPs(filter, input []netaddr.IPPrefix, nonmatch bool) []netaddr.IPPrefix {
	var res []netaddr.IPPrefix
	for _, prefix := range input {
		match := false
		isv6 := prefix.IP().Is6()
		for _, cidr := range filter {
			if isv6 != cidr.IP().Is6() {
				continue
			}
			if cidr.Bits() <= prefix.Bits() && cidr.Contains(prefix.IP()) {
				match = true
				break
			}
		}
		if match != nonmatch {
			res = append(res, prefix)
		}
	}
	return res
}

func matchFile(filter []netaddr.IPPrefix, nonmatch bool, inFile string) {
	ips := readFile(inFile)
	for _, ip := range filterIPs(filter, ips, nonmatch) {
		fmt.Println(ip)
	}
}

func main() {
	filterFile := flag.String("F", "", "file with networks to filter for")
	nonmatchFlag := flag.Bool("v", false, "Show non-matching entries instead of matching ones")
	flag.Parse()

	if *filterFile == "" {
		fmt.Println("missing -F parameter with filter file")
		return
	}

	filter := readFile(*filterFile)

	if len(flag.Args()) == 0 {
		matchFile(filter, *nonmatchFlag, "-")
	} else {
		for _, inFile := range flag.Args() {
			matchFile(filter, *nonmatchFlag, inFile)
		}
	}
}
