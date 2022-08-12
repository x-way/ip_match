package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strings"
)

func readFile(name string) []netip.Prefix {
	var prefixes []netip.Prefix
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
		prefix, err := netip.ParsePrefix(line)
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

func filterIPs(filter, input []netip.Prefix, nonmatch bool) []netip.Prefix {
	var res []netip.Prefix
	for _, prefix := range input {
		match := false
		isv6 := prefix.Addr().Is6()
		for _, cidr := range filter {
			if isv6 != cidr.Addr().Is6() {
				continue
			}
			if cidr.Bits() <= prefix.Bits() && cidr.Contains(prefix.Addr()) {
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

func matchFile(filter []netip.Prefix, nonmatch bool, inFile string) {
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
