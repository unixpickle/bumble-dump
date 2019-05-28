import json
import os
import subprocess

for profile_name in os.listdir('profiles'):
    with open(os.path.join('profiles', profile_name), 'r') as f:
        data = json.load(f)
    for album in data['albums']:
        if 'photos' not in album:
            continue
        for photo in album['photos']:
            img_path = os.path.join('photos/' + photo['id'] + '.jpg')
            if os.path.exists(img_path):
                continue
            print(photo['id'])
            img_url = 'https:' + photo['large_url']
            data = subprocess.check_output(['curl', img_url], stderr=subprocess.DEVNULL)
            with open(img_path, 'wb+') as f:
                f.write(data)
