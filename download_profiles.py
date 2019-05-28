import json
import os
import subprocess

while True:
    output = subprocess.check_output(['./commands/encounters.sh'], stderr=subprocess.DEVNULL)
    data = json.loads(output)
    for entry in data['body'][0]['client_encounters']['results']:
        data = entry['user']
        user_id = data['user_id']
        dest_path = os.path.join('profiles', user_id + '.json')
        if not os.path.exists(dest_path):
            print(user_id)
            with open(dest_path, 'w') as f:
                json.dump(data, f)
        subprocess.check_call(['./commands/dislike.sh', user_id],
                              stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
