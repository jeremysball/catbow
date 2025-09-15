#!/usr/bin/env fish
set cwd "$(pwd)"
pushd /tmp
if not test -d cb-e2e
    mkdir cb-e2e
end
pushd cb-e2e
command $cwd/generate_text.py --line-width 4 --num-lines 1 1>data.txt
cat data.txt | catbow -freq .31 -spread 3 -seed 420 1>cb.txt
cat data.txt | lolcat -S 420 --force 1>lc.txt
diff cb.txt lc.txt
popd
popd
