package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"
)

func sanitizeFileName(input string) string {
	re := regexp.MustCompile(`[^\w]+`)
	return re.ReplaceAllString(input, "_")
}

func captureScreenshot(ctx context.Context, targetURL, outputDir, proxy string, timeout time.Duration, headers map[string]string) error {
	// Configure the Chrome context with proxy support if provided
	opts := chromedp.DefaultExecAllocatorOptions[:]
	if proxy != "" {
		opts = append(opts, chromedp.ProxyServer(proxy))
	}
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Apply timeout if specified
	ctxWithTimeout, cancel := context.WithTimeout(browserCtx, timeout)
	defer cancel()

	// Enable the network and set custom headers
	err := chromedp.Run(ctxWithTimeout,
		network.Enable(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			headerMap := network.Headers{}
			for key, value := range headers {
				headerMap[key] = value
			}
			return network.SetExtraHTTPHeaders(headerMap).Do(ctx)
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to enable network and set headers: %v", err)
	}

	// Add tasks for navigation and taking a screenshot
	var screenshot []byte
	tasks := []chromedp.Action{
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.CaptureScreenshot(&screenshot),
	}

	if err := chromedp.Run(ctxWithTimeout, tasks...); err != nil {
		if ctxWithTimeout.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout: page load exceeded %v for %s", timeout, targetURL)
		}
		return fmt.Errorf("failed to capture screenshot for %s: %v", targetURL, err)
	}

	// Generate filename
	fileName := fmt.Sprintf(
		"%s_%s.png",
		sanitizeFileName(targetURL),
		time.Now().Format("20060102_150405"),
	)
	filePath := filepath.Join(outputDir, fileName)

	// Save the screenshot
	if err := os.WriteFile(filePath, screenshot, 0644); err != nil {
		return fmt.Errorf("failed to save screenshot for %s: %v", targetURL, err)
	}

	log.Printf("Screenshot saved: %s\n", filePath)
	return nil
}

func processURLs(urls []string, outputDir, proxy string, threads int, timeout time.Duration, headers map[string]string) {
	var wg sync.WaitGroup
	urlChan := make(chan string, threads)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for targetURL := range urlChan {
				log.Printf("Processing URL: %s\n", targetURL)
				if err := captureScreenshot(context.Background(), targetURL, outputDir, proxy, timeout, headers); err != nil {
					log.Printf("Error processing %s: %v\n", targetURL, err)
				}
			}
		}()
	}

	for _, targetURL := range urls {
		urlChan <- targetURL
	}
	close(urlChan)

	wg.Wait()
}

func parseHeaders(headerFlag string) map[string]string {
	headers := make(map[string]string)
	if headerFlag != "" {
		headerPairs := strings.Split(headerFlag, ",")
		for _, pair := range headerPairs {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) == 2 {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}
	return headers
}

func main() {
	// Flags
	urlFlag := flag.String("u", "", "Target URL (required unless -f is provided)")
	fileFlag := flag.String("f", "", "File containing list of URLs (one per line)")
	outputDirFlag := flag.String("o", ".", "Output directory for screenshots")
	proxyFlag := flag.String("proxy", "", "Proxy server to use (e.g., http://127.0.0.1:8080)")
	threadsFlag := flag.Int("threads", 4, "Number of goroutines for processing URLs")
	timeoutFlag := flag.Int("t", 0, "Number of seconds to wait for page load (0 for unlimited)")
	headerFlag := flag.String("H", "", "Custom headers to add to the browser requests (comma-separated key:value)")

	flag.Parse()

	// Validate flags
	if *urlFlag == "" && *fileFlag == "" {
		log.Fatal("Error: You must specify either -u or -f flag")
	}

	// Create output directory if not exists
	if _, err := os.Stat(*outputDirFlag); os.IsNotExist(err) {
		if err := os.MkdirAll(*outputDirFlag, os.ModePerm); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	var urls []string
	if *urlFlag != "" {
		urls = append(urls, *urlFlag)
	}

	if *fileFlag != "" {
		file, err := os.Open(*fileFlag)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				urls = append(urls, line)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
	}

	// Convert timeout to duration
	var timeout time.Duration
	if *timeoutFlag > 0 {
		timeout = time.Duration(*timeoutFlag) * time.Second
	} else {
		timeout = time.Hour * 24 * 365 * 100 // Effectively unlimited timeout
	}

	headers := parseHeaders(*headerFlag)
	processURLs(urls, *outputDirFlag, *proxyFlag, *threadsFlag, timeout, headers)
}
