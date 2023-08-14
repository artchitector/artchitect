#!/bin/bash

# Собираем jpeg-файлы оригиналов в архив, чтобы создать резервную копию данных
cd /var/artchitect/origin || exit
# # ls -l
# total 304
#drwxr-xr-- 2 root root 36864 Aug 13 08:12 U000XXX
#drwxr-xr-- 2 root root 36864 Aug 14 01:08 U001XXX
#drwxr-xr-- 2 root root 12288 Aug 14 11:15 U002XXX

# Все файлы разложены по их тысячным единствам (ищи unity в коде)
# То есть в одной папке ровно 1000 картинок в оригинальном размере
# /var/artchitect/origin/U000XXX/ ls -l
# ...
#-rwxr-xr-- 1 root root 5147284 Aug 13 08:10 art-997.jpg
#-rwxr-xr-- 1 root root 8216553 Aug 13 08:11 art-998.jpg
#-rwxr-xr-- 1 root root 3712821 Aug 13 08:12 art-999.jpg
#-rwxr-xr-- 1 root root 4895604 Aug 12 16:56 art-99.jpg
#-rwxr-xr-- 1 root root 4510420 Aug 12 15:25 art-9.jpg

echo "### [НАЧИНАЮ СОХРАНЕНИЕ ОРИГИНАЛОВ] ###"
lastFolder=$(ls -tr | tail -1)
preFolder=$(ls -tr | tail -2 | head -n 1)
echo "- игнорирую папку $lastFolder, так как она последняя и еще не заполнена полностью"
echo "- работаю с папкой $preFolder, так как она предпоследняя заполнена полностью."
echo "Остальные архивы должны были собраться ранее."

filename="<home_path>/dumps/origin/$preFolder.tar.gz"
if [ -e "$filename" ]; then
    echo "[!!!] АРХИВ $filename УЖЕ СОБРАН. ПРОПУСКАЮ"
  else
    echo "[!!!] НАЧИНАЮ СОБИРАТЬ АРХИВ $filename"
    tempfile="<home_path>/tmp/$preFolder.tar.gz"
    tar -zcvf $tempfile $preFolder
    mv $tempfile $filename
    echo "[!!!] АРХИВ $filename СОБРАН"
  fi
