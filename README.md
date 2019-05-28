# bumble-dump

This is a small set of tools for dumping all the Bumble profiles in your area. It requires that you look at the requests made from your browser on Bumble's webapp, and convert those requests into small shell scripts. This can be done easily using Chrome's "Copy as cURL" tool.

# Usage

Create a `commands` directory with two files:

 * dislike.sh - a script that takes one argument (user ID) and dislikes them. Use endpoint `https://bumble.com/unified-api.phtml?SERVER_ENCOUNTERS_VOTE`.
 * encounters.sh - a script that takes no arguments and dumps encounters JSON. Use endpoint `https://bumble.com/unified-api.phtml?SERVER_GET_ENCOUNTERS`.

Now create a `profiles` directory. Next run the `download_profiles.py` script. At any time you can create a `photos` directory run `download_photos.py` to download photos as well.
