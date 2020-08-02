package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/mikioh/ipaddr"
)

func readFile(name string) ([]ipaddr.Prefix, []ipaddr.Prefix) {
	var prefixesv6 []ipaddr.Prefix
	var prefixesv4 []ipaddr.Prefix
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
		_, ipNet, err := net.ParseCIDR(line)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(line, ":") {
			prefixesv6 = append(prefixesv6, *(ipaddr.NewPrefix(ipNet)))
		} else {
			prefixesv4 = append(prefixesv4, *(ipaddr.NewPrefix(ipNet)))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return prefixesv6, prefixesv4
}

func filterIPs(filter, input []ipaddr.Prefix, nonmatch bool) []ipaddr.Prefix {
	var res []ipaddr.Prefix
	for _, ip := range input {
		match := false
		for _, net := range filter {
			if net.Contains(&ip) {
				match = true
				break
			}
		}
		if match != nonmatch {
			res = append(res, ip)
		}
	}
	return res
}

func matchFile(filter6, filter4 []ipaddr.Prefix, nonmatch bool, inFile string) {
	ips6, ips4 := readFile(inFile)
	for _, ip := range filterIPs(filter6, ips6, nonmatch) {
		fmt.Println(ip)
	}
	for _, ip := range filterIPs(filter4, ips4, nonmatch) {
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

	filter6, filter4 := readFile(*filterFile)

	if len(flag.Args()) == 0 {
		matchFile(filter6, filter4, *nonmatchFlag, "-")
	} else {
		for _, inFile := range flag.Args() {
			matchFile(filter6, filter4, *nonmatchFlag, inFile)
		}
	}
}
