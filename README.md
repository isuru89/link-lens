
<!-- [![Build](https://github.com/isuru89/link-lens/actions/workflows/ci.yml/badge.svg)](https://github.com/isuru89/link-lens/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/isuru89/link-lens/graph/badge.svg?token=KRTRYTDDK7)](https://codecov.io/gh/isuru89/link-lens) -->


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
  * Is a login form or not?

## How to Run

**Prerequisites:**
  * Should have installed latest [go version](https://go.dev/doc/install) (required minimum v1.21.6)

### Steps

 1. Clone the repository and then move to the folder.
```bash
git clone https://github.com/isuru89/link-lens.git && cd link-lens
```

 2. Execute the below command.
```bash
go run server/...
```

#### Using UI

 3. Open a browser and navigate to `http://localhost:8080`
 4. Type a URL you want to analyze.
 5. Click `Analyze` button.

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

### FAQs

* __How do you find inaccessible links?__

   Inaccessible links are broken or invalid links that no successful HTTP code is returned when navigated. Even if it returns a valid HTML content but contains response status a non-`2xx`, then it considers as an inaccessible link.

* __What protocols do you support?__

    Currently this program supports only `http` or `https` protocols. Any other protocols will be treated as invalid.

* __When finding inaccessible links, does the program check only http/https protocols?__

    Yes. Links with other protocols will be ignored. They will not be treated as inaccessible links btw.