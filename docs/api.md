# Scoreboard API

In order to decouple the scoreboard application's backend functionality from the user interface, a
REST-like API is used to serve and modify data.

**Contents**

* Notation and Conventions
* Event Endpoints
  * Information about the event
* Team Endpoints
  * List of teams and their scores
* Submission Endpoints
* Administration Endpoints


## Notation and Conventions

The scoreboard API's endpoints will each be documented with the expected inputs and outputs, and
each value will have an accompanying type, inspired by
[TypeScript](https://www.typescriptlang.org/docs/home.html#toc-handbook)'s notation.

All of the endpoints implemented by this API accept (only) JSON when a body is present and will always
respond with JSON. Thus, a response with an `error: string | null` (an error that can be either
a string or `null`) and a `newScore: number` field might look like the following.

```json
{
    "error": "Invalid flag",
    "newScore": 0
}
```

The major exception to this rule is, of course, the case of GET requests. Since GET requests should
not contain data in their body, input parmeters to such requests are expected to be present in the
request URL's [query string](https://en.wikipedia.org/wiki/Query_string) using the standard
`<url>?paramName=value&otherParma=otherValue` format. Inputs will, even in this case, still be
described as a JSON object. Just replace the JSON keys with query string parameters.

## Event Endpoints

### Information about hte event

    GET /event

Obtain descriptive information about the event.

#### Parameters

None

#### Response

```js
{
    "name": string,
}
```

## Team Endpoitns

### List of teams and their scores

    GET /teams/scoreboard

Obtain a list of teams and all of the information needed to populate a scoreboard.

#### Parameters

None

#### Response

```js
{
    "teams": [
        {
            "name": string,
            "score": number,
            "position": number,
            "members": string,
            "lastSubmission": string
        }
    ]
}
```