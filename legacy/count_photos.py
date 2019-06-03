import json
import os
import subprocess

photo_set = set()

for profile_name in os.listdir('profiles'):
    with open(os.path.join('profiles', profile_name), 'r') as f:
        data = json.load(f)
    for album in data['albums']:
        if 'photos' not in album:
            continue
        for photo in album['photos']:
            photo_set.add(photo['id'])
            if not len(photo_set) % 1000:
                print('found %d photos so far' % len(photo_set))
print('found %d photos' % len(photo_set))
