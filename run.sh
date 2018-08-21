git fetch --all
git reset --hard origin/master
git fetch

go build -o ssManager .
./ssManager
