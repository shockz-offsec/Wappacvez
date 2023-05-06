<div align="center">
  <h1>Wappacvez</h1>
</div>
<div align="center">
  <img src="" align="center">
</div>

Wappacvez is a command-line tool that analyzes a web application by using a dockerized Wappalyzer. It then extracts the software for which a version is detected, and finally employs the uCVE tool to search for associated CVEs. The output can be exported in HTML or CSV format, depending on the user's preference.

<div align="center">
  <h2>Requirementes</h2>
</div>

* Linux or Mac
* Go (version 1.16+)

<div align="center">
  <h2>Installation</h2>
</div>

To install Wappacvez, run the following command:

```
go get github.com/shockz-offsec/wappacvez
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
* `-cvss`: Filter vulnerabilities by CVSS [critical,high,medium,low,none] (default:all)
* `-lg`: Set language of information [en,es] (default:en)
* `-oHTML`: Save CVEs list in HTML file [filename] (default:report.html)
* `-oCSV`: Save CVEs list in CSV file [filename]

| The only mandatory argument is the url

<div align="center">
  <h2>Examples</h2>
</div>

```bash
wappacvez -u "https://www.nasa.gov" -oHTML "nasa.html" -cvss critical,high
```

<div align="center">
  <h2>Details</h2>
</div>

Wappacvez will proceed to install Docker and build my Wappalyzer image and install uCVE on the system.

<div align="center">
  <h2>Credits</h2>
</div>

[Wappalyzer](https://github.com/wappalyzer/wappalyzer)
[uCVE](https://github.com/m3n0sd0n4ld/uCVE)


<div align="center">
  <h2>License</h2>
</div>

This tool is licensed under the GPL-3.0 License.
