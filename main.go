package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ConfigFile string
	JdFile     string
	CvSkills   []string `yaml:"skills"`
	JdSkills   []string
}

type Content struct {
	Name         string
	Email        string
	Phone        string
	Residence    string
	Header       string
	Education    string
	Certificates string
	Experience   string
	LinkedIn     string
	CvSkills     []string
	CommonSkills []string
}

// readJdSkillsFromFile reads a file and extracts skills from it.
// It takes a filename as input and returns a slice of strings containing the extracted skills.
// If there is an error while reading the file, it returns nil and the error.
func readJdSkillsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var skills []string
	for scanner.Scan() {
		skill := strings.ToLower(scanner.Text())
		skills = append(skills, strings.Split(skill, ",")...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}

func loadConfigFromFile(filename string, config *Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}

func loadContentFromConfig(config *Config) (*Content, error) {
	data, err := os.ReadFile(config.ConfigFile)
	if err != nil {
		return nil, err
	}
	var Content Content
	err = yaml.Unmarshal(data, &Content)
	if err != nil {
		return nil, err
	}
	return &Content, nil
}

// runHTML2PDFCommand executes the HTML to PDF conversion command.
// It runs the "node js/html2pdf.js" command and prints the output to the console.
// Returns an error if the command execution fails.
func runHTML2PDFCommand() error {
	cmd := exec.Command("node", "js/html2pdf.js")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// generateHTMLFile generates an HTML file based on the provided content.
// It takes a pointer to a Content struct as input and returns an error if any.
// The generated HTML file includes personal information, experience, education,
// certificates, and skills. It uses a predefined HTML template to structure the content.
// The generated file is saved as "output.html" in the current directory.
func generateHTMLFile(content *Content) error {
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>CV</title>
		<link rel="stylesheet" type="text/css" href="style.css">
	</head>
	<body>
		<h1>Curriculum Vitae</h1>
		<p>{{.Header}}</p>
		<h2>Personal Information</h2>
		<p>Name: {{.Name}}</p>
		<p>Email: {{.Email}}</p>
		<p>Phone: {{.Phone}}</p>
		<p>LinkedIn: <a href={{.LinkedIn}}>linkedin.com/mikko-turpeinen</a></p>
		<p>Residence: {{.Residence}}</p>
		<h2>Experience</h2>
		<p>{{.Experience}}</p>
		<h2>Education</h2>
		<p>{{.Education}}</p>
		<h2>Certificates</h2>
		<p>{{.Certificates}}</p>
		<h2>Skills</h2>
		<ul>
			{{range .CommonSkills}}
			<li>{{.}}</li>
			{{end}}
		</ul>
	</body>
	</html>
	`

	htmlFile, err := os.Create("output.html")
	if err != nil {
		return err
	}

	defer htmlFile.Close()

	t := template.Must(template.New("cvTemplate").Parse(htmlTemplate))

	err = t.Execute(htmlFile, content)
	if err != nil {
		return err
	}

	return nil
}

// findCommonSkills finds the common skills between the CV skills and the job description (JD) skills.
// It takes two slices of strings, cvSkills and jdSkills, and returns a new slice containing the common skills.
func findCommonSkills(cvSkills []string, jdSkills []string) []string {
	commonSkills := make([]string, 0)
	for _, cvSkill := range cvSkills {
		for _, jdSkill := range jdSkills {
			if cvSkill == jdSkill {
				commonSkills = append(commonSkills, cvSkill)
				break
			}
		}
	}
	return commonSkills
}

func main() {
	var config Config

	config.ConfigFile = "config.yaml"
	config.JdFile = "jdskills.txt"

	err := loadConfigFromFile(config.ConfigFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	content, err := loadContentFromConfig(&config)
	if err != nil {
		log.Fatal(err)
	}
	content.CvSkills = config.CvSkills

	jdSkills, err := readJdSkillsFromFile(config.JdFile)
	if err != nil {
		log.Fatal(err)
	}

	commonSkills := findCommonSkills(content.CvSkills, jdSkills)
	content.CommonSkills = commonSkills
	generateHTMLFile(content)
	runHTML2PDFCommand()
}
