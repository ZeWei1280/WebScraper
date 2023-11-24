# Web Scraper
## About
This go program is designed to scrape table of test results from the web pages and save it in structured CSV format.

## Installation
To install and run this program, please ensure that GO has already been installed  in your system. And follow the steps:
1. Clone this repository to your local machine:
    ```shell
    git clone [THIS_PEPOSITORY]
    ```
2. Navigate to the cloned directory in the terminal and build this directory
    ``` shell
    cd PATH/TO/YOUR/CLONED/REPOSITOAY
    ```
    ```shell
    go build -o [PROGRAM_NAME]
    ```
3. Run the executable file that you just built: 
    ```shell
    ./[PROGRAM_NAME]
    ```
## Usage
### Flags

When you're executing this program, you can run it with command-line arguments to specify various parameters:
* `-dir` or `-d`:
  * set the working directory
  * default is `"./test data"`
* `-outputDir` or `-o`
  * set the directory where the output csv files will be saved
  * default is `"./results"`
* `-concurrency` or `-c`:
  * to set the number of go routines process at the same time
  * default is `1`
### Example
```shell
./[PROGRAM_NAME] -d="test data" -o="results" -c=1
```
## Output
The program will generate CSV files in the specified output directory. These files contain the data extracted from the web pages.

## Program Structure
The program consists of several parts:

* `main.go`: The entry point of the program which sets up the environment and starts the web scraping process.
* `scraper.go`: Contains the `VisitAndScrapePage()` function for visiting web pages and the `scrapeSubpage()` function for scraping individual pages.
* `csv_builder.go`: Defines the CSVBuilder type, which constructs the CSV files from the scraped data.
* /utils:
    * `server.go`: Start a local HTTP which servers the data for scraper testing.
    * `flags.go`: Handles the parsing of command-line flags.
    * `logger.go`: Log setting and management.