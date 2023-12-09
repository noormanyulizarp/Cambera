// Root directory structure (./main.go):
// .
// ├── main.go (<-)

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	separator := createSeparator(20)
	initializeReaderIgnore()
	processDirectoryStructure(".", "files_structure.txt", separator)
}

func createSeparator(length int) string {
	return strings.Repeat("-", length)
}

func initializeReaderIgnore() {
	if _, err := os.Stat(".readerignore"); os.IsNotExist(err) {
		defaultPatterns := []string{
			"# Version Control",
			".git",
			"",
			"# REPLIT",
			".local",
			".config",
			".cache",
			"",
			"# Node",
			"node_modules",
			"",
			"# Logs",
			"*.log",
			"",
			"# IDEs and Editors",
			".vscode",
			".idea",
			"*.iml",
			"*.ipr",
			"*.iws",
			"*~",
			"*.swp",
			"",
			"# Operating System",
			".DS_Store",
			"Thumbs.db",
			"",
			"# Reader",
			".readerignore",
			"files_structure.txt",
		}
		writeFile(".readerignore", strings.Join(defaultPatterns, "\n"))
	}
}

func writeFile(filename, content string) {
	if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
		panic(err)
	}
}

func processDirectoryStructure(startPath, outputFilePath, separator string) {
	excludePatterns := readReaderIgnore()
	paths := walkDirectory(startPath, excludePatterns, startPath)

	outputFile := createFile(outputFilePath)
	defer outputFile.Close()

	printContents(paths, outputFile, separator)
}

func readReaderIgnore() []string {
	var patterns []string
	file, err := os.Open(".readerignore")
	if err != nil {
		return patterns
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && line != "" {
			patterns = append(patterns, line)
		}
	}
	return patterns
}

func walkDirectory(startPath string, excludePatterns []string, relativePath string) []string {
	var paths []string
	filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && path != startPath && !shouldExclude(path, excludePatterns, relativePath) {
			paths = append(paths, path)
		}
		return nil
	})
	sort.Strings(paths)
	return paths
}

func shouldExclude(path string, patterns []string, relativePath string) bool {
	rel, err := filepath.Rel(relativePath, path)
	if err != nil {
		return false
	}
	for _, pattern := range patterns {
		if strings.HasPrefix(rel, pattern) {
			return true
		}
	}
	return false
}

func createFile(filename string) *os.File {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func printContents(paths []string, outputFile *os.File, separator string) {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		printFileDetails(path, info, outputFile, separator)
	}
}

func printFileDetails(path string, info os.FileInfo, outputFile *os.File, separator string) {
	sizeStr := formatFileSize(info.Size())
	fmt.Fprintf(outputFile, "// The size of (%s): %s\n", path, sizeStr)
	fmt.Fprintln(outputFile, separator)
	fmt.Fprintf(outputFile, "// The file location of (%s):\n", path)
	fmt.Fprintln(outputFile, "// .")
	fmt.Fprintf(outputFile, "// ├── %s (<-)\n", filepath.Base(path))
	fmt.Fprintln(outputFile, separator)
	if !info.IsDir() {
		content, _ := ioutil.ReadFile(path)
		fmt.Fprintf(outputFile, "// The content of (%s):\n", path)
		fmt.Fprintln(outputFile, string(content))
		fmt.Fprintln(outputFile, separator)
	}
}

func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}
	return fmt.Sprintf("%dKB", size/1024)
}
