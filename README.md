# bumble-dump

This is a comprehensive set of tools for dumping Bumble dating profiles into a database and downloading accompanying image data.

# Usage

## Generating a Bumble API

The first step is to generate an "API" that allows the scanner to access a bumble account. In particular, the API contains URL endpoints for disliking users, changing your location, and listing potential matches. It also includes cookies and headers necessary to authenticate these requests.

To generate an API, you must use the bumble website in a browser with developer tools enabled, then export all the requests as a HAR document, and run this HAR document through the [generate_api](generate_api/generate_api.go) command.

The resulting API is a JSON file that can be fed to the [scan](scan/) command.

## Configuration

The configuration is specified via a few environment variables. Here are the variables:

 * `BUMBLE_DB`: a MongoDB database URI. **Default:** `mongodb://localhost:27017`.
 * `BUMBLE_PHOTOS`: the directory path for storing profile photos. **Default:** `./photos`.

## Scanning

Simply run the `scan` command and pipe it into `scan_dump`:

```
go run scan/*.go | go run scan_dump/*.go
```

## Finding geocoordinates

The `scan` command dumps raw user profiles, and a user profile doesn't come with an exact set of geocoordinates. Instead, it comes with a string such as `Philadelphia, PA`.

The `find_locations` command populates a collection in the database mapping location strings to geocoordinates. Once the location collection is populated, you can use the database to search for users within a certain distance of a given location.

## Indexes

At some point, you may want to setup indexes on the database so that users can be found faster. This can be done with the `setup_indexes` command.
