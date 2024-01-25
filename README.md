
[![Build](https://github.com/isuru89/link-lens/actions/workflows/ci.yml/badge.svg)](https://github.com/isuru89/link-lens/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/isuru89/link-lens/graph/badge.svg?token=KRTRYTDDK7)](https://codecov.io/gh/isuru89/link-lens)


# Link-Lens

Analyze a link and reports basic information.


Following information will be reported:
  * HTML version
  * Title of the site
  * Heading level counts
  * Link statistics
     * No of external links
     * No of internal links
     * No of inaccessible links (invalid/broken links)
     * Full URLs of inaccessible links
  * Is a login form or not?

## How to Run

**Prerequisites:**
  * Should have installed latest [go version](https://go.dev/doc/install) (required minimum v1.21.6)
  * Latest node.js (at least node v1.18.x or above) (*This is required only if you are expecting to use UI*)

### Steps

 1. Clone the repository and then move to the folder.
```
git clone https://github.com/isuru89/link-lens.git && cd link-lens
```

 2. Execute the below command to get dependencies and execute the server.
```
go get ./...
go run linklens/...
```

  By default, the server will be running on port `8080`. If you want to change this behaviour,
  you can use `--port` argument to pass another port.

```
go run linklens/... --port 8070
```

**Note:** If you want to build an executable instead of running directly with source code, execute the below commands to get an executable.

```
cd main
go build -o linklens
```

Above command will create a binary file called `linklens` inside the `main` folder and you can execute it using below command.

```
./linklens
```

If you want to run with web UI, then you need to pass the path of the built web artifacts using `webDir` argument. The path must be relative to the executable location.

```
./linklens -webDir=../web/build
```

#### Using UI

Make sure you are running the program using any of the above methods. And also, at least node v1.18.x installed.

 * Install and build the web page using below set of commands.
```
cd web && npm i && npm run build
```

 * This will build the required web artifacts and will be served through the application server.
 * Open a browser and navigate to `http://localhost:8080`. If you have changed the port, use it instead of 8080.
 * Type a URL you want to analyze.
 * Click `Analyze` button.

#### Using API

We have also exposed a API end point directly to analyze a URL. You can examine it by invoking below cUrl command.

```
curl --request POST \
  --url http://localhost:8080/api/analyze \
  --header 'Content-Type: application/json' \
  --data '{
	"url": "https://github.com"
}'
```

It will return a response similar to what shown below, if the operation was successful with a 200 HTTP status code.

```json
{
   "SourceUrl": "https://github.com",
   "HtmlVersion": "5",
   "Title": "YouTube",
   "HeadingsCount": {
      "H1": 1,
      "H4": 2
   },
   "LinkStats": {
      "InternalLinkCount": 5,
      "ExternalLinkCount": 8,
      "InvalidLinkCount": 1,
      "InvalidLinks": [
         "https://non-existence.com/url"
      ]
   },
   "PageType": "Unknown"
}
```

### Configurations

The link-lens program will accept below configurations via command line arguments.

  * `-port`: Port of the server. (*Default port is 8080*)
  * `-ui`: Whether to serve UI or not (*Default is yes*)
  * `-webDir`: Directory to the web portal artifacts (*Default is ./web/build*)

At anytime, it is possible to know about accepting arguments by invoking help command.

```
./linklens -h
```

### Improvements

  * Caching: We could implement a caching for serving same two URLs in a short time instead of analyzing twice.
  * More Information: Like broken image links, identify sign-up form or different page types
  * Automatic continuous deployment of this project to a hosting site using Github Actions.

### Limitations

  * Single-Page Applications (SPA) will not report the accurate statistics, because of the unavailability of page structure.
  * Links will be verified only in the given anlysis URL. No recursive crawling is supported to identify dead-links further deep under the given website.

### FAQs

* __How do you find inaccessible links?__

   Inaccessible links are broken or invalid links that no successful HTTP code is returned when navigated. Even if it returns a valid HTML content but contains response status a non-`2xx`, then it considers as an inaccessible link.

* __What protocols do you support?__

    Currently this program supports only `http` or `https` protocols. Any other protocols will be treated as invalid.

* __When finding inaccessible links, does the program check only http/https protocols?__

    Yes. Links with other protocols will be ignored. They will not be treated as inaccessible links btw.