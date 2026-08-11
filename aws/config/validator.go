package awsconfig

import (
	"bufio"
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

// ErrPromptAborted means the user ended an interactive prompt with Ctrl-C or EOF, or
// that stdin was never interactive. Callers should treat it as a deliberate abort and
// stop rather than re-prompt: a non-interactive stdin returns EOF immediately, so
// retrying spins.
var ErrPromptAborted = errors.New("prompt aborted")

// promptStdin is the source the interactive selectors read from. Nil in production,
// where readline and the scanners fall back to os.Stdin; tests set it to drive the
// prompts without a terminal.
var promptStdin io.ReadCloser

func promptReader() io.Reader {
	if promptStdin != nil {
		return promptStdin
	}
	return os.Stdin
}

func StdinRegionSelector() (string, error) {
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
		Stdin:        promptStdin,
	})
	if err != nil {
		return "", fmt.Errorf("selecting region: %w", err)
	}
	defer rl.Close()

	for !IsValidRegion(region) {
		line, err := rl.Readline()
		switch {
		case errors.Is(err, readline.ErrInterrupt), errors.Is(err, io.EOF):
			// Ctrl-C or EOF while choosing a region aborts. This used to be one of
			// two os.Exit calls left outside main, because the enclosing signature
			// had no error to return; it now does.
			return "", ErrPromptAborted
		case err != nil:
			return "", fmt.Errorf("selecting region: %w", err)
		}

		region = strings.TrimSpace(line)
		if !IsValidRegion(region) {
			fmt.Fprintf(os.Stderr, "'%s' is not a valid region\n", region)
		}
	}

	return region, nil
}

func StdinInstanceTypeSelector() (string, error) {
	fmt.Println("Please choose one instance type")
	fmt.Println()
	fmt.Println("Here are few examples:")

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

	// Read with a scanner rather than fmt.Scan. fmt.Scan returns its error without
	// assigning, so the previous `for !isValid(...)` loop spun at full speed forever
	// on a closed or piped stdin instead of giving up.
	scanner := bufio.NewScanner(promptReader())
	for {
		fmt.Print("Value ? > ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", fmt.Errorf("selecting instance type: %w", err)
			}
			return "", ErrPromptAborted
		}
		instanceType := strings.TrimSpace(scanner.Text())
		if isValidInstanceType(instanceType) {
			return instanceType, nil
		}
		fmt.Printf("'%s' is not a valid instance type\n", instanceType)
	}
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
