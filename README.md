<div align="center">
  <h1>Wappacvez</h1>
</div>
<div align="center">
  <img src="https://user-images.githubusercontent.com/67438760/236649071-9bc6c030-7ff0-40c9-8dfd-e4e096f0b30a.png" align="center">
</div>


Wappacvez is a command-line tool that analyzes a web application by using a dockerized Wappalyzer. It then extracts the software for which a version is detected, and finally employs the uCVE tool to search for associated CVEs. The output can be exported in HTML or CSV format, depending on the user's preference.

<div align="center">
  <h2>Requirementes</h2>
</div>

* Linux or Mac
* Go (version 1.16+)
* Docker

<div align="center">
  <h2>Installation</h2>
</div>

To install Wappacvez, run the following command:

```
go install -v github.com/shockz-offsec/wappacvez@latest
```
or via building via repository
```
git clone https://github.com/shockz-offsec/Wappacvez.git
cd Wappacvez
go build -o wappacvez wappacvez.go
```

<div align="center">
  <h2>Download the compiled binary for Linux or MacOS</h2>
</div>

[Download the latest version](https://github.com/shockz-offsec/Wappacvez/releases)


<div align="center">
  <h2>Usage</h2>
</div>

```bash
wappacvez -u <url> [-cvss value] [-lg value] [-oHTML value.html] [-oCSV value.csv]
```

* `-u`: URL to scan (mandatory)
* `-cvss`: Filter vulnerabilities by CVSS [critical,high,medium,low,none] (default: all)
* `-lg`: Set language of information [en,es] (default: en)
* `-oHTML`: Save CVEs list in HTML file [filename] (default: report.html)
* `-oCSV`: Save CVEs list in CSV file [filename]

| The only mandatory argument is the url

<div align="center">
  <h2>Examples</h2>
</div>

```bash
wappacvez -u "https://www.nasa.gov" -oHTML "nasa.html" -cvss critical,high
```

Output

<div align="center">
  <img src="https://user-images.githubusercontent.com/67438760/236649945-fc9ab712-489f-47fb-89d7-f4cd8b8c705f.png" align="center">
</div>

<div align="center">
  <h2>Details</h2>
</div>

Wappacvez will proceed to install Docker and build my Wappalyzer image and install uCVE on the system.

| Due to the limitations of using the Wappalyzer core versus the extension, it is possible that some websites may not detect all software versions compared to the extension.
| We considered using the official API, but this free API has more limitations in terms of queries and results.


<div align="center">
  <h2>Dockerized Wappalyzer</h2>
</div>

Dockerized version of Wappalyzer developed for this tool.

https://hub.docker.com/r/shockzoffsec/wappalyzer

With the following command the latest available version will be installed and executed.

```bash
docker run --rm shockzoffsec/wappalyzer:latest <url> [arguments]
```

All Wappalyzer options are allowed.

```
Usage:
  wappalyzer <url> [options]

Examples:
  wappalyzer https://www.example.com
  node cli.js https://www.example.com -r -D 3 -m 50 -H "Cookie: username=admin"
  docker wappalyzer/cli https://www.example.com --pretty

Options:
  -b, --batch-size=...       Process links in batches
  -d, --debug                Output debug messages
  -t, --delay=ms             Wait for ms milliseconds between requests
  -h, --help                 This text
  -H, --header               Extra header to send with requests
  --html-max-cols=...        Limit the number of HTML characters per line processed
  --html-max-rows=...        Limit the number of HTML lines processed
  -D, --max-depth=...        Don't analyse pages more than num levels deep
  -m, --max-urls=...         Exit when num URLs have been analysed
  -w, --max-wait=...         Wait no more than ms milliseconds for page resources to load
  -p, --probe=[basic|full]   Perform a deeper scan by performing additional requests and inspecting DNS records
  -P, --pretty               Pretty-print JSON output
  --proxy=...                Proxy URL, e.g. 'http://user:pass@proxy:8080'
  -r, --recursive            Follow links on pages (crawler)
  -a, --user-agent=...       Set the user agent string
  -n, --no-scripts           Disabled JavaScript on web pages
  -N, --no-redirect          Disable cross-domain redirects
  -e, --extended             Output additional information
  --local-storage=...        JSON object to use as local storage
  --session-storage=...      JSON object to use as session storage
  --defer=ms                 Defer scan for ms milliseconds after page load
```


<div align="center">
  <h2>Credits</h2>
</div>

[Wappalyzer](https://github.com/wappalyzer/wappalyzer)

[uCVE](https://github.com/m3n0sd0n4ld/uCVE)


<div align="center">
  <h2>License</h2>
</div>

This tool is licensed under the GPL-3.0 License.
