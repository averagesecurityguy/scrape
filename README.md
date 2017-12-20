# Scrape
Scrape is extensible Pastebin scraper written in Golang. To use scrape without getting blacklisted at Pastebin.com you will need to get a Lifetime Pro membership and whitelist your IP address. Scrape implements Pastebin's recommended scraping logic, which is defined at https://pastebin.com/api_scraping_faq.

## Files

    * config.go - Configure Scrape including any keywords for which you want to search.
    * paste.go - Defines the paste struct.
    * process.go - Handles all of the paste processing logic. Add new processing logic here.
    * scrape.go - The main function.
