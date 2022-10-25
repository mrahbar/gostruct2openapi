package main

import (
	"flag"
	"fmt"
	"github.com/mrahbar/gostruct2openapi/doc"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	packagesFlag := flag.String("packages", "", "comma separated package to scan")
	filterFlag := flag.String("filter", ".*", "regular expression used to filter struct names")
	flag.Parse()

	filter := regexp.MustCompile(*filterFlag)
	packages := parsePackages(packagesFlag)

	if len(packages) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	generator := doc.NewOpenapiGenerator(filter, "json")
	specs, err := generator.DocumentStruct(packages...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(specs)
}

func parsePackages(packagesFlag *string) (res []string) {
	for _, s := range strings.Split(*packagesFlag, ",") {
		if len(s) > 0 {
			res = append(res, s)
		}
	}

	return
}
