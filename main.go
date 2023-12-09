package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	SeparatorLength    = 20
	ReaderIgnoreFile   = ".readerignore"
	OutputFileName     = "files_structure.txt"
	DefaultPermissions = 0644
)

func main() {
	rootPath := "."
	ensureIgnoreFileExists()
	paths := getPathsToProcess(rootPath, readIgnorePatterns())
	processDirectoryStructure(rootPath, paths)
}

func ensureIgnoreFileExists() {
	if _, err := os.Stat(ReaderIgnoreFile); os.IsNotExist(err) {
		content := getDefaultIgnorePatterns()
		err := ioutil.WriteFile(ReaderIgnoreFile, []byte(content), DefaultPermissions)
		if err != nil {
			fmt.Printf("Error creating %s: %v\n", ReaderIgnoreFile, err)
			os.Exit(1)
		}
	}
}

func getDefaultIgnorePatterns() string {
	patterns := []string{
		"# Version Control", ".git", "",
		"# REPLIT", ".local", ".config", ".cache", "",
		"# Node", "node_modules", "",
		"# Logs", "*.log", "",
		"# IDEs and Editors", ".vscode", ".idea", "*.iml", "*.ipr", "*.iws", "*~", "*.swp", "",
		"# Operating System", ".DS_Store", "Thumbs.db", "",
		"# Reader", ".readerignore", "files_structure.txt", "",
		"# Additional Patterns", ".project-rc", "__pycache__/", "*.py[cod]", "*$py.class", "*.so",
		".Python", "build/", "develop-eggs/", "dist/", "downloads/", "eggs/", ".eggs/", "lib/", "lib64/",
		"parts/", "sdist/", "var/", "wheels/", "*.egg-info/", ".installed.cfg", "*.egg", "MANIFEST",
		"*.manifest", "*.spec", "pip-log.txt", "pip-delete-this-directory.txt", "htmlcov/", ".tox/",
		".coverage", ".coverage.*", ".cache", "nosetests.xml", "coverage.xml", "*.cover", ".hypothesis/",
		".pytest_cache/", "core.*", "*.mo", "*.pot", "*.log", "local_settings.py", "db.sqlite3", "instance/",
		".webassets-cache", ".scrapy", "docs/_build/", "target/", ".ipynb_checkpoints", ".python-version",
		"celerybeat-schedule", "*.sage.py", "/site", ".mypy_cache/", "",
		"# Media files", "*.mp4", "*.jpg", "*.jpeg", "*.png", "*.gif", "*.bmp", "*.tiff", "*.ico",
	}
	return strings.Join(patterns, "\n")
}

func getPathsToProcess(rootPath string, excludePatterns []string) []string {
	var paths []string
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}
		if !shouldExclude(path, excludePatterns) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", rootPath, err)
		os.Exit(1)
	}
	sort.Strings(paths)
	return paths
}

func readIgnorePatterns() []string {
	var patterns []string
	file, err := os.Open(ReaderIgnoreFile)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", ReaderIgnoreFile, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return patterns
}

func processDirectoryStructure(rootPath string, paths []string) {
	outputFile, err := os.Create(OutputFileName)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	for _, path := range paths {
		printFileDetails(rootPath, path, outputFile)
	}
}

func printFileDetails(rootPath, path string, outputFile *os.File) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(outputFile, "Error getting info for %s: %v\n", path, err)
		return
	}

	sizeStr := formatFileSize(info.Size())
	relPath, _ := filepath.Rel(rootPath, path)

	separator := strings.Repeat("-", SeparatorLength)
	fmt.Fprintf(outputFile, "%s\n// The size of (%s): %s\n", separator, relPath, sizeStr)
	fmt.Fprintf(outputFile, "%s\n// The file location of (%s):\n", separator, relPath)
	printPath(outputFile, relPath)
	fmt.Fprintf(outputFile, "%s\n", separator)

	if !info.IsDir() {
		content, _ := ioutil.ReadFile(path)
		fmt.Fprintf(outputFile, "\n// The content of (%s):\n", relPath)
		fmt.Fprintln(outputFile, string(content))
		fmt.Fprintf(outputFile, "%s\n", separator)
	}
}

func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}
	return fmt.Sprintf("%dKB", size/1024)
}

func shouldExclude(path string, patterns []string) bool {
	relPath, err := filepath.Rel(".", path)
	if err != nil {
		return false
	}

	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, relPath)
		if err != nil {
			fmt.Printf("Error matching pattern %s: %v\n", pattern, err)
			continue
		}

		for _, part := range strings.Split(relPath, string(os.PathSeparator)) {
			dirMatched, _ := filepath.Match(pattern, part)
			if dirMatched {
				return true
			}
		}

		if matched {
			return true
		}
	}
	return false
}

func printPath(outputFile *os.File, relPath string) {
	dirs := strings.Split(relPath, string(os.PathSeparator))
	indent := ""
	for i, dir := range dirs {
		if i == len(dirs)-1 {
			fmt.Fprintf(outputFile, "%s└── %s (<-)\n", indent, dir)
		} else {
			fmt.Fprintf(outputFile, "%s├── %s\n", indent, dir)
			indent += "│   "
		}
	}
}
