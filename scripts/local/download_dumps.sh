  # запускается локально для стягивания дампов на какой-то локальный компьютер
  LOCAL_DUMPS_PATH=/mnt/storage/artchitect
  REMOTE_DUMPS_PATH=/root/dumps
  mkdir -p $LOCAL_DUMPS_PATH/
  mkdir -p $LOCAL_DUMPS_PATH/db
  mkdir -p $LOCAL_DUMPS_PATH/arts
  mkdir -p $LOCAL_DUMPS_PATH/origin

  echo "НАЧИНАЮ СКАЧИВАНИЕ ДАМПА БД"
  lastDump=$(ssh memory.artchitect "cd $REMOTE_DUMPS_PATH/db && ls -tr | tail -1")
  echo "ПОСЛЕДНИЙ БД-ДАМП - $lastDump"
  path=$LOCAL_DUMPS_PATH/db/$lastDump
  if [ -e "$path" ]; then
    echo "ДАМП $lastDump УЖЕ СКАЧАН. ПРОПУСК"
  else
    scp memory.artchitect:$REMOTE_DUMPS_PATH/db/$lastDump $path
    echo "УСПЕШНО СКАЧАН ФАЙЛ $lastDump to $path"
  fi

  echo "НАЧИНАЮ СКАЧИВАНИЕ ДАМПА ARTS"
  lastDump=$(ssh memory.artchitect "cd $REMOTE_DUMPS_PATH/arts && ls -tr | tail -1")
  echo "ПОСЛЕДНИЙ ARTS-ДАМП - $lastDump"
  path=$LOCAL_DUMPS_PATH/arts/$lastDump
  if [ -e "$path" ]; then
    echo "ДАМП $lastDump УЖЕ СКАЧАН. ПРОПУСК"
  else
    scp memory.artchitect:$REMOTE_DUMPS_PATH/arts/$lastDump $path
    echo "УСПЕШНО СКАЧАН ФАЙЛ $lastDump to $path"
  fi

  echo "НАЧИНАЮ СКАЧИВАНИЕ ДАМПА ОРИГИНАЛОВ"
  lastDump=$(ssh storage.artchitect "cd $REMOTE_DUMPS_PATH/origin && ls -tr | tail -1")
  echo "ПОСЛЕДНИЙ ORIGIN-ДАМП - $lastDump"
  path=$LOCAL_DUMPS_PATH/origin/$lastDump
  if [ -e "$path" ]; then
    echo "ДАМП $path УЖЕ СКАЧАН. ПРОПУСК"
  else
    scp storage.artchitect:$REMOTE_DUMPS_PATH/origin/$lastDump $path
    echo "УСПЕШНО СКАЧАН ФАЙЛ $lastDump to $path"
  fi