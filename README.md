### 3s - Save the ScreenShot ###

**3s (Save the ScreenShot)** is a command-line tool written in Go that automates the process of capturing website screenshots using a headless Chrome browser. It supports custom headers, proxy settings, multithreading, and timeouts.

### Features ###
- 🚀 **Capture website screenshots** using headless Chrome.
- 🔥 **Custom headers support** (e.g., User-Agent, Authorization).
- 🌍 **Proxy support** to bypass restrictions.
- ⚡ **Multithreading** for faster execution.
- ⏳ **Timeout management** to prevent infinite page loads.
- 📁 **Save screenshots** with sanitized filenames including timestamps.

### Installation ###
```bash
git clone <repository-url>
cd <repository-folder>
go build -o 3s
```

### Usage ###
```bash
./3s -u "<URL>" [options]
```

### Command-Line Options ###
| Flag | Description | Example |
|------|-------------|---------|
| `-u` | Target URL to capture | `-u "https://example.com"` |
| `-f` | File containing multiple URLs (one per line) | `-f "urls.txt"` |
| `-o` | Output directory for screenshots | `-o "screenshots"` |
| `-H` | Custom headers (comma-separated key:value pairs) | `-H "User-Agent:Custom,Authorization:Token123"` |
| `-proxy` | Proxy server to use | `-proxy "http://127.0.0.1:8080"` |
| `-threads` | Number of concurrent browser instances | `-threads 4` |
| `-t` | Timeout in seconds for page loading | `-t 15` |

### Examples ###

**1. Capture a screenshot of a single URL**
```bash
./3s -u "https://example.com"
```

**2. Capture multiple URLs from a file**
```bash
./3s -f urls.txt -o screenshots
```

**3. Capture a URL with custom headers**
```bash
./3s -u "https://example.com" -H "User-Agent:CustomAgent,Authorization:Bearer 123"
```

**4. Use a proxy server**
```bash
./3s -u "https://example.com" -proxy "http://127.0.0.1:8080"
```

**5. Capture screenshots with 8 threads**
```bash
./3s -f urls.txt -threads 8
```

### Note on Performance ###
**Warning:** Using too many threads can make your computer slow because each thread launches a headless browser instance. Adjust the `-threads` value based on your system’s capabilities.

### Dependecies ###
	github.com/chromedp/cdproto/network
	github.com/chromedp/chromedp
	golang.org/x/net/context

### License ###
3s is licensed under the SushiWare license. For more information, check [docs/license.txt](docs/license.txt).
