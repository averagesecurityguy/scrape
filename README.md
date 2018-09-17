# Scrape
Scrape finds interesting data in text files using keyword searches and regular expressions. Scrape pulls text files from Pastebin and Github Gists. In addition, Scrape can parse text files in a local directory. The search terms are user configurable and are stored in the config.json file. Scrape can run in the background as a service or it can run on demand.

## Sources

### Pastebin
To use scrape without getting blacklisted at Pastebin.com you will need to get a Lifetime Pro membership and whitelist your IP address. Scrape implements Pastebin's recommended scraping logic, which is defined at https://pastebin.com/api_scraping_faq.

### Gists
To use scrape with Github Gists, you will need to create a read-only Github API key. Scrape gets the 100 most recent gists using the API endpoint described at: https://developer.github.com/v3/gists/#list-all-public-gists. At this time, no attempt is made to download truncated files or truncated content.

### Local Files
To use scrape to parse files in a local directory, define the directory in the config.json file. Scrape will parse the files in batches of 100 by default. The batch size is configurable in the config.json file. Keep in mind, that after a file is processed it will be deleted from the directory.

## Installation
You will first need to clone the Git repository with `git clone https://github.com/averagesecurityguy/scrape`. Once you have downloaded the repository, run the setup.sh script from the repository with sudo permissions. This will generate a new user called scrape and install the service.sh init script. If you already have a service account you want to use on your machine, then modify the setup.sh script to use the account name you want.

## Viewing Gathered Data
While scrape is running you can visit https://127.0.0.1:5000 to view the data that has been gathered. You will need to
create a TLS certificate and key and define their locations in the config.json file. When scrape is not runnig you can use the view tool in the install directory to view scrape data.

### View Command Usage
```
Usage:
    view filename action [arguments]

Actions:
    buckets                       Get a list of buckets.
    read <bucketname> <key>       Get the value of the key in the bucket.
    keys <bucketname>             Get a list of keys in a bucket.
    vals <bucketname>             Get a list of values in a bucket.
    search <bucketname> <string>  Get a list of keys from the bucket where the
                                  value contains the given string.
```
