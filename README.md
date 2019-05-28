# Usage

Create a commands/ directory with two files:

 * dislike.sh - a script that takes one argument (user ID) and dislikes them. Use endpoint `https://bumble.com/unified-api.phtml?SERVER_ENCOUNTERS_VOTE`.
 * encounters.sh - a script that takes no arguments and dumps encounters JSON. Use endpoint `https://bumble.com/unified-api.phtml?SERVER_GET_ENCOUNTERS`.

Now create two directories: `profiles` and `photos`. Next run the `download_profiles.py` script.

At any time you can create a `photos` directory run `download_photos.py` to download photos as well.
