package main

import (
	"flag"
	"fmt"
	"os"
	"encoding/json"
    "strings"
	"os/exec"
	"net/url"
	"regexp"
)

const (
    red     = "\033[31m"
    green   = "\033[32m"
    yellow  = "\033[33m"
    magenta = "\033[35m"
    cyan    = "\033[36m"
    reset   = "\033[0m"
)

func flow(urlFlag string, args string) {
	dockerInstalled := checkDockerInstalled()
	if !dockerInstalled {
		fmt.Println(red + "[!] Docker is not installed" + reset )
	}

	uCVE_path, err := verify_install_uCVE()
	if err != nil {
		fmt.Println(red + "[!] " + reset + "Error: ", err)
		return
	}
	
	err = checkDockerImage()
    if err != nil {
        fmt.Println(red + "[!] " + reset + "Error: ",err)
        return
    }

	fmt.Printf(yellow + "[+] " + reset + "Running wappalyzer...\n")

	out, err := exec.Command("docker", "run", "--rm", "shockzoffsec/wappalyzer", urlFlag).Output()
	if err != nil {
		errMsg := fmt.Sprintf("%s[!] %sError: %v%s", red, reset, err, reset)
		panic(errMsg)
	}

	fmt.Printf(yellow + "[+] " + reset + "Processing results...\n")
	versions, err := parseTechnologiesJSON(out)
	if err != nil {
		errMsg := fmt.Sprintf("%s[!] %sError: %v%s", red, reset, err, reset)
		panic(errMsg)
	}

	fmt.Printf(green + "[+] " + reset + "Software with versions detected:\n")
	for product, version := range versions {
        fmt.Printf("\n%s %s", product, version)
    }

	// Call writeVersionsToJSONFile to write the parsed technology versions to a file
	err = writeVersionsToJSONFile(versions, "technologies.json")
	if err != nil {
		errMsg := fmt.Sprintf("%s[!] %sError: %v%s", red, reset, err, reset)
		panic(errMsg)
	}

	fmt.Printf(yellow + "\n\n[+] " + reset + "Obtaining CVE's\n")
	uCVE(uCVE_path, args, "technologies.json")
}

func parseTechnologiesJSON(jsonBytes []byte) (map[string]string, error) {

    var data map[string]interface{}
    if err := json.Unmarshal(jsonBytes, &data); err != nil {
        return nil, err
    }

    versions := make(map[string]string)
    for _, tech := range data["technologies"].([]interface{}) {
        techMap := tech.(map[string]interface{})
        var cpe, name, version string
        if cpeVal, ok := techMap["cpe"].(string); ok {
            cpe = cpeVal
            if versionVal, ok := techMap["version"].(string); ok {
                version = versionVal
            }
            name = strings.ReplaceAll(strings.Split(cpe, ":")[4], "-", "_")
        } else if slugVal, ok := techMap["slug"].(string); ok {
            slug := slugVal
            if versionVal, ok := techMap["version"].(string); ok {
                version = versionVal
            }
            name = strings.ReplaceAll(slug, "-", "_")
        }
        if len(version) > 0 && regexp.MustCompile(`^\d+[a-zA-Z]?(\.\d+[a-zA-Z]?)*$`).MatchString(version) {
            versions[name] = version
        }
    }

    return versions, nil
}

func writeVersionsToJSONFile(versions map[string]string, filename string) error {
    outputFile, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    jsonStr, err := json.Marshal(versions)
    if err != nil {
        return err
    }

    if _, err := outputFile.Write(jsonStr); err != nil {
        return err
    }

    return nil
}


func uCVE(uCVE_path string, args string, filename string) {
	uCVEArgs := []string{"-iJSON", "technologies.json"}
    if args != "" {
        uCVEArgs = append(uCVEArgs, strings.Split(args, " ")...)
    }

	uCVECmd := exec.Command(uCVE_path, uCVEArgs...)
	uCVECmd.Stdout = os.Stdout
	uCVECmd.Stderr = os.Stderr

	if err := uCVECmd.Run(); err != nil {
		errMsg := fmt.Sprintf("%s[!] %sError: %v%s", red, reset, err, reset)
		panic(errMsg)
	}
}

func checkDockerInstalled() bool {
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func checkDockerImage() error {
    cmd := exec.Command("docker", "images")
    out, err := cmd.Output()
    if err != nil {
        return fmt.Errorf(red + "[!] " + reset + "Error when executing the command 'docker images'.: %v", err)
    }

    if !strings.Contains(string(out), "wappalyzer") {
		fmt.Println(yellow + "[+] " + reset + "Image 'shockzoffsec/wappalyzer' not found. Building...")
        cmd := exec.Command("docker", "pull", "shockzoffsec/wappalyzer:latest")
        err := cmd.Run()
        if err != nil {
            return fmt.Errorf(red + "[!] " + reset + "Error pulling the image 'wappalyzer'.: %v", err)
        }
        fmt.Println(green + "[+] " + reset + "wappalyzer' image built successfully")
    } else {
        fmt.Println(green + "[+] " + reset + "Image 'wappalyzer' is already built.")
    }

    return nil
}

func verify_install_uCVE() (string, error) {
	tool := "uCVE"

	path, err := exec.LookPath(tool)
	if err != nil {
		fmt.Printf(yellow + "[+] " + reset + "The tool %s is not installed. Installing...\n", tool)
		cmd := exec.Command("go", "install", "github.com/m3n0sd0n4ld/uCVE@latest")
		if err := cmd.Run(); err != nil {
			fmt.Printf(red + "[!]" + reset + "Error: Could not install the tool %s.\n", tool)
			return "", err
		}
		path, err = exec.LookPath(tool)
		if err != nil {
			fmt.Printf(red + "[!] " + reset + "Error: Tool %s was not found after installation.\n", tool)
			return "", err
		}
		fmt.Printf(green + "[+] " + reset + "The %s tool was successfully installed on: %s\n", tool, path)
	} else {
		fmt.Printf(green + "[+] " + reset + "The tool %s is already installed in: %s\n", tool, path)
	}
	return path, nil
}

func validateFlags(cvss string, language string, htmlFile string, csvFile string) bool {
    validCVSS := false
    if cvss == "" {
        validCVSS = true
    } else {
        cvssValues := strings.Split(cvss, ",")
        for _, v := range cvssValues {
            switch v {
            case "critical", "high", "medium", "low", "none":
                validCVSS = true
            default:
                validCVSS = false
                break
            }
        }
    }
    if !validCVSS {
        fmt.Println(red + "[!] Invalid CVSS value\n" + reset)
        flag.Usage()
        os.Exit(1)
    }

    validLanguage := false
    if language == "" {
        validLanguage = true
    } else {
        for _, v := range []string{"en", "es"} {
            if v == language {
                validLanguage = true
                break
            }
        }
    }
    if !validLanguage {
        fmt.Println(red + "[!] Invalid language value\n" + reset)
        flag.Usage()
        os.Exit(1)
    }

	if htmlFile != "" {
        validFilename := regexp.MustCompile(`^[a-zA-Z0-9_]+\.[hH][tT][mM][lL]$`).MatchString(htmlFile)
        if !validFilename {
            fmt.Println(red + "[!] Invalid HTML filename. Must have .html extension\n" + reset)
            flag.Usage()
            os.Exit(1)
        }
    }

    if csvFile != "" {
        validFilename := regexp.MustCompile(`^[a-zA-Z0-9_]+\.[cC][sS][vV]$`).MatchString(csvFile)
        if !validFilename {
            fmt.Println(red + "[!] Invalid CSV filename. Must have .csv extension\n" + reset)
            flag.Usage()
            os.Exit(1)
        }
    }
    return true
}


func main() {
	urlFlag := flag.String("u", "", "URL to scan (mandatory)")
    cvssFlag := flag.String("cvss", "", "Filter vulnerabilities by CVSS [critical,high,medium,low,none]")
    languageFlag := flag.String("lg", "", "Set language of information [en,es]")
    htmlFileFlag := flag.String("oHTML", "", "Save CVEs list in HTML file [filename]")
    csvFileFlag := flag.String("oCSV", "", "Save CVEs list in CSV file [filename]")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage:\n")
        fmt.Fprintf(os.Stderr, "\twappacvez -u [URL] [OPTIONS]\n\n")
        fmt.Fprintf(os.Stderr, "Options:\n")
        flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Examples:\n")
        fmt.Fprintf(os.Stderr, "\twappacvez -u \"https://www.nasa.gov\" -oHTML \"nasa.html\" -cvss critical,high\n\n")
    }

    flag.Parse()

    if flag.NFlag() == 0 || *urlFlag == "" {
        flag.Usage()
        os.Exit(0)
    }

    if _, err := url.ParseRequestURI(*urlFlag); err != nil {
        fmt.Println(red + "[!] Invalid URL\n" + reset )
        flag.Usage()
        os.Exit(1)
    }

    validateFlags(*cvssFlag, *languageFlag, *htmlFileFlag, *csvFileFlag)

    args := ""
    if *cvssFlag != "" {
        args += fmt.Sprintf("-cvss %s ", *cvssFlag)
    }
    if *languageFlag != "" {
        args += fmt.Sprintf("-lg %s ", *languageFlag)
    }
    if *htmlFileFlag != "" {
        args += fmt.Sprintf("-oHTML %s ", *htmlFileFlag)
    }else if *csvFileFlag != "" {
        args += fmt.Sprintf("-oCSV %s ", *csvFileFlag)
    }else{
        args += fmt.Sprintf("-oHTML report.html")
    }


	var BANNER =`
_       _____    ____  ____  ___   _______    _____________
| |     / /   |  / __ \/ __ \/   | / ____/ |  / / ____/__  /
| | /| / / /| | / /_/ / /_/ / /| |/ /    | | / / __/    / / 
| |/ |/ / ___ |/ ____/ ____/ ___ / /___  | |/ / /___   / /__
|__/|__/_/  |_/_/   /_/   /_/  |_\____/  |___/_____/  /____/
                                                  by Shockz
`

	fmt.Println(BANNER)

	flow(*urlFlag, args)
}