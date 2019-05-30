import os
import sys

from PIL import Image

for name in os.listdir('photos'):
    try:
        path = os.path.join('photos', name)
        img = Image.open(path)
        img.thumbnail(512)
        img.save(path)
    except:  # noqa
        sys.stderr.write('failed for %s\n' % name)
