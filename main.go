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

const separatorLength = 20

func main() {
	initializeReaderIgnore()
	processDirectoryStructure(".", "files_structure.txt")
}

func initializeReaderIgnore() {
	defaultIgnorePatterns := getDefaultIgnorePatterns()
	createFileWithContent(".readerignore", defaultIgnorePatterns)
}

func getDefaultIgnorePatterns() string {
	patterns := []string{
		"# Version Control", ".git", "",
		"# REPLIT", ".local", ".config", ".cache", "",
		"# Node", "node_modules", "",
		"# Logs", "*.log", "",
		"# IDEs and Editors", ".vscode", ".idea", "*.iml", "*.ipr", "*.iws", "*~", "*.swp", "",
		"# Operating System", ".DS_Store", "Thumbs.db", "",
		"# Reader", ".readerignore", "files_structure.txt",
	}
	return strings.Join(patterns, "\n")
}

func createFileWithContent(filename, content string) {
	if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
		panic(err)
	}
}

func processDirectoryStructure(startPath, outputFilePath string) {
	excludePatterns := readIgnorePatterns(".readerignore")
	paths := getDirectoryPaths(startPath, excludePatterns)

	outputFile := createFile(outputFilePath)
	defer outputFile.Close()

	printDirectoryDetails(paths, outputFile)
}

func readIgnorePatterns(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && line != "" {
			patterns = append(patterns, line)
		}
	}
	return patterns
}

func getDirectoryPaths(startPath string, excludePatterns []string) []string {
	var paths []string
	filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
		if path != startPath && !shouldExclude(path, excludePatterns, startPath) {
			paths = append(paths, path)
		}
		return nil
	})
	sort.Strings(paths)
	return paths
}

func shouldExclude(path string, patterns []string, relativePath string) bool {
	rel, _ := filepath.Rel(relativePath, path)
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

func printDirectoryDetails(paths []string, outputFile *os.File) {
	separator := strings.Repeat("-", separatorLength)
	for _, path := range paths {
		info, _ := os.Stat(path)
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
