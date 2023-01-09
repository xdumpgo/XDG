#!/bin/bash

echo  "Dumping ${1} to ascii hex..."
hexdump -ve '1/1 "%.2X"' $1 > tmp.bin

echo "Grabbing strings from ${1}..."
for NAME in $(strings $1 | egrep -i "^((\*|\#)|)(git\.zertex\.space|github\.com)" | sort -u -r)
do
	NEW_STRING="$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 20 ; echo '')"
	OLD_STRING_HEX="$(echo -n ${NAME} | xxd -g 0 -u -ps -c 256)00"
	NEW_STRING_HEX="$(echo -n ${NEW_STRING} | xxd -g 0 -u -ps -c 256)00"

	if [ ${#NEW_STRING_HEX} -le ${#OLD_STRING_HEX} ] ; then
		while [ ${#NEW_STRING_HEX} -lt ${#OLD_STRING_HEX} ] ; do
			NEW_STRING_HEX="${NEW_STRING_HEX}00"
		done
		
		echo -n "Replacing ${NAME} with ${NEW_STRING}"
		sed -i "s/${OLD_STRING_HEX}/${NEW_STRING_HEX}/g" tmp.bin

		echo " Done!"
	else
		echo "New string '${NEW_STRING}' is longer than old '${OLD_STRING}', skipping."
	fi
done

xxd -r -p tmp.bin patched-$1