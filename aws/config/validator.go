package awsconfig

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/chzyer/readline"
)

var AWSHomeDir = func() string {
	var home string
	if runtime.GOOS == "windows" { // Windows
		home = os.Getenv("USERPROFILE")
	} else {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".aws")
}

func ParseRegion(i string) (any, error) {
	if !IsValidRegion(i) {
		return i, fmt.Errorf("'%s' is not a valid region", i)
	}
	return i, nil
}

func ParseInstanceType(i string) (any, error) {
	if !isValidInstanceType(i) {
		return i, fmt.Errorf("'%s' is not a valid instance type", i)
	}
	return i, nil
}

func StdinRegionSelector() string {
	var regionItems []readline.PrefixCompleterInterface
	for _, r := range allRegions() {
		regionItems = append(regionItems, readline.PcItem(r))
	}
	var regionCompleter = readline.NewPrefixCompleter(regionItems...)

	fmt.Println("Please enter one region: (Ctrl+C to quit, Tab for completion)")
	var region string
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		AutoComplete: regionCompleter,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while selecting region: %s", err)
		return ""
	}
	defer rl.Close()

	for !IsValidRegion(region) {
		line, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) || errors.Is(err, io.EOF) {
			os.Exit(1)
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error while selecting region: %s", err)
			return ""
		}

		region = strings.TrimSpace(line)
		if !IsValidRegion(region) {
			fmt.Fprintf(os.Stderr, "'%s' is not a valid region\n", region)
		}
	}

	return region
}

func StdinInstanceTypeSelector() string {
	fmt.Println("Please choose one instance type")
	fmt.Println()
	fmt.Println("Here are few examples:")

	var instanceType string
	t := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(t, "\tinstance type\tvCPU\tMemory (GiB)")
	fmt.Fprintln(t, "\tt2.nano\t1\t0.5")
	fmt.Fprintln(t, "\tt2.micro\t1\t1")
	fmt.Fprintln(t, "\tt2.small\t1\t2")
	fmt.Fprintln(t, "\tt2.medium\t2\t4")
	fmt.Fprintln(t, "\tt2.large\t2\t8")
	fmt.Fprintln(t, "\tt2.xlarge\t4\t16")
	fmt.Fprintln(t, "\tt2.2xlarge\t8\t32")
	fmt.Fprintln(t, "\tm4.large\t2\t8")
	fmt.Fprintln(t, "\tm4.xlarge\t4\t16")
	fmt.Fprintln(t, "\tc4.large\t2\t3.75")
	fmt.Fprintln(t, "\tc4.xlarge\t4\t7.5")
	fmt.Fprintln(t, "\t...")
	t.Flush()

	fmt.Println()
	fmt.Print("Value ? > ")
	fmt.Scan(&instanceType)
	for !isValidInstanceType(instanceType) {
		fmt.Printf("'%s' is not a valid instance type\n", instanceType)
		fmt.Print("Value ? > ")
		fmt.Scan(&instanceType)
	}
	return instanceType
}

func IsValidRegion(given string) bool {
	reg, _ := regexp.Compile(`^(us|eu|ap|sa|ca|af|me|il|cn)\-\w+\-\d+$`)
	regCompound, _ := regexp.Compile(`^(us\-gov|us\-iso|us\-isob|eu\-isoe)\-\w+\-\d+$`)

	return reg.MatchString(given) || regCompound.MatchString(given)
}

func isValidInstanceType(given string) bool {
	return regexp.MustCompile(`\w+\.\w+`).MatchString(given)
}

func allRegions() []string {
	regions := sort.StringSlice{
		"af-south-1",
		"ap-east-1",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-south-1",
		"ap-south-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-southeast-3",
		"ap-southeast-4",
		"ap-southeast-5",
		"ca-central-1",
		"ca-west-1",
		"cn-north-1",
		"cn-northwest-1",
		"eu-central-1",
		"eu-central-2",
		"eu-north-1",
		"eu-south-1",
		"eu-south-2",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"il-central-1",
		"me-central-1",
		"me-south-1",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-gov-east-1",
		"us-gov-west-1",
		"us-west-1",
		"us-west-2",
	}
	sort.Sort(regions)
	return regions
}

func IsValidProfile(given string) bool {
	return stringInSlice(given, AllProfiles())
}

var awsHomeFunc func() string = AWSHomeDir

var profileNameRegex = regexp.MustCompile(`\[(.*)\]`)

func AllProfiles() (profiles []string) {
	awsHome := awsHomeFunc()
	files := []string{filepath.Join(awsHome, "config"), filepath.Join(awsHome, "credentials")}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			continue
		}
		out, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		matches := profileNameRegex.FindAllSubmatch(out, -1)
		for _, match := range matches {
			profile := string(match[1])
			profile = strings.TrimSpace(profile)
			profile = strings.TrimPrefix(profile, "profile ")
			profile = strings.TrimSpace(profile)
			if profile != "" {
				profiles = append(profiles, profile)
			}
		}
	}
	return profiles
}

func stringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
