#!/bin/bash

# Собираем jpeg-файлы карточек в архив, чтобы создать резервную копию данных
cd /var/artchitect/arts || exit
# # ls -l
# total 304
# drwxr-xr-- 2 root root 135168 Aug 13 08:12 U000XXX
# drwxr-xr-- 2 root root 135168 Aug 14 01:09 U001XXX
# drwxr-xr-- 2 root root  36864 Aug 14 10:50 U002XXX

# Все файлы разложены по их тысячным единствам (ищи unity в коде)
# То есть в одной папке ровно 4000 картинок (1000 картин по 4 размера - F, M, S, XS
# /var/artchitect/arts/U000XXX/ ls -l
# ...
#-rwxr-xr-- 1 root root 247610 Aug 13 08:12 art-999-f.jpg
#-rwxr-xr-- 1 root root  62156 Aug 13 08:12 art-999-m.jpg
#-rwxr-xr-- 1 root root  18432 Aug 13 08:12 art-999-s.jpg
#-rwxr-xr-- 1 root root   6400 Aug 13 08:12 art-999-xs.jpg
#-rwxr-xr-- 1 root root 423628 Aug 12 16:57 art-99-f.jpg
#-rwxr-xr-- 1 root root  94749 Aug 12 16:57 art-99-m.jpg
#-rwxr-xr-- 1 root root  23002 Aug 12 16:57 art-99-s.jpg
#-rwxr-xr-- 1 root root   6673 Aug 12 16:57 art-99-xs.jpg
#-rwxr-xr-- 1 root root 362836 Aug 12 15:25 art-9-f.jpg
#-rwxr-xr-- 1 root root  95421 Aug 12 15:25 art-9-m.jpg
#-rwxr-xr-- 1 root root  28513 Aug 12 15:25 art-9-s.jpg
#-rwxr-xr-- 1 root root   9033 Aug 12 15:25 art-9-xs.jpg

echo "### [НАЧИНАЮ СОХРАНЕНИЕ КАРТИН] ###"
lastFolder=$(ls -tr | tail -1)
preFolder=$(ls -tr | tail -2 | head -n 1)
echo "- игнорирую папку $lastFolder, так как она последняя и еще не заполнена полностью"
echo "- работаю с папкой $preFolder, так как она предпоследняя заполнена полностью."
echo "Остальные архивы должны были собраться ранее."

filename="<home_path>/dumps/arts/$preFolder.tar.gz"
if [ -e "$filename" ]; then
    echo "[!!!] АРХИВ $filename УЖЕ СОБРАН. ПРОПУСКАЮ"
  else
    echo "[!!!] НАЧИНАЮ СОБИРАТЬ АРХИВ $filename"
    tempfile="<home_path>/tmp/$preFolder.tar.gz"
    tar -zcvf $tempfile $preFolder
    mv $tempfile $filename
    echo "[!!!] АРХИВ $filename СОБРАН"
  fi
