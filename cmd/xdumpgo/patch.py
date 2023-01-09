import re
import string
import random

print('Loading binary..')
f = open('xdumpgo.exe', 'rb')
s = f.read()
f.close()

re.sub(b'^((\*|\#)|)(git.zertex\.space|github\.com)', ''.join(random.choices(string.ascii_uppercase + string.ascii_lowercase + string.digits, k = 12)), s)

f = open('xdumpgo.exe', 'w')

f.write(s.decode())