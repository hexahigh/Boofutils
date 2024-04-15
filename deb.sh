#!/bin/bash

if [ -z "$1" ]; then
    echo "Please provide a version number as an argument."
    exit 1
fi

VERSION=$1
FILE_TO_MODIFY="./templates/deb/DEBIAN/control"
STRING_TO_FIND="%{VERSION}%"
STRING_TO_REPLACE="${VERSION}"
FILE_TO_COPY="./boofutils"
DEB_DIR="./templates/deb"

sed -i "s/${STRING_TO_FIND}/${STRING_TO_REPLACE}/g" ${FILE_TO_MODIFY}

cp ${FILE_TO_COPY} ${DEB_DIR}/usr/local/bin/

dpkg-deb --build ${DEB_DIR} boofutils_${VERSION}_amd64.deb

# Undo SED changes
sed -i "s/${STRING_TO_REPLACE}/${STRING_TO_FIND}/g" ${FILE_TO_MODIFY}

echo "Deb package created: boofutils_${VERSION}_amd64.deb"