# bumble-dump

This is a comprehensive set of tools for dumping Bumble dating profiles into a database and downloading accompanying image data.

# Configuration

The configuration is specified via a few environment variables. Here are the variables:

 * `BUMBLE_DB`: a MongoDB database URI. **Default:** `mongodb://localhost:27017`.
 * `BUMBLE_IMAGES`: the directory path for storing profile images. **Default:** `./photos`.
