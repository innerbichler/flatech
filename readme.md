<div align="center">
<img src="./mountain_fuzzy.png" width="350">
</div>

# Flatech 

is a project to perform various task in the Flatex Webportal with go.
The ./webWorker does all of the browser work.

You can find some binaries in ./binaries but they are just artifacts

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
Build with earthly

To build the scraper component do:
~~~
earthly +scraper
~~~
and then to start it just do
~~~
docker compose -f start-scraper.yaml up -d
~~~

