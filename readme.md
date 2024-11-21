<div align="center">
<img src="./mountain_fuzzy.png" width="350">
</div>

# Flatech 

is a webWorker using selenium to perform various task in the Flatex Webportal in the CLI

## Getting Started
add your credentials into a .env file in the source directory
~~~
USERID=
PASSWORD=
~~~

- make sure you have firefox installed 
~~~
sudo apt install firefox
~~~
this should automatically install the geckodriver

### Build
Build uses earthly.

To build the scraper component do:
~~~
earthly +scraper
~~~
and then to start it just do
~~~
docker compose up -d -f start-scraper.yaml
~~~

