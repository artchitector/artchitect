  # запускается локально для стягивания дампов на какой-то локальный компьютер
  mkdir -p ~/dumps/
  mkdir -p ~/dumps/db
  mkdir -p ~/dumps/arts
  mkdir -p ~/dumps/origin

  echo "НАЧИНАЮ СКАЧИВАНИЕ ДАМПА БД"
  lastDump=$(ssh memory "cd ~/dumps/db && ls -tr | tail -1")
  echo "ПОСЛЕДНИЙ БД-ДАМП - $lastDump"
  path=~/dumps/db/$lastDump
  if [ -e "$path" ]; then
    echo "ДАМП $lastDump УЖЕ СКАЧАН. ПРОПУСК"
  else
    scp memory:~/dumps/db/$lastDump $path
    echo "УСПЕШНО СКАЧАН ФАЙЛ $lastDump to $path"
  fi

  echo "НАЧИНАЮ СКАЧИВАНИЕ ДАМПА ARTS"
  lastDump=$(ssh memory "cd ~/dumps/arts && ls -tr | tail -1")
  echo "ПОСЛЕДНИЙ ARTS-ДАМП - $lastDump"
  path=~/dumps/arts/$lastDump
  if [ -e "$path" ]; then
    echo "ДАМП $lastDump УЖЕ СКАЧАН. ПРОПУСК"
  else
    scp memory:~/dumps/arts/$lastDump $path
    echo "УСПЕШНО СКАЧАН ФАЙЛ $lastDump to $path"
  fi

  echo "НАЧИНАЮ СКАЧИВАНИЕ ДАМПА ОРИГИНАЛОВ"
  lastDump=$(ssh storage "cd ~/dumps/origin && ls -tr | tail -1")
  echo "ПОСЛЕДНИЙ ORIGIN-ДАМП - $lastDump"
  path=~/dumps/origin/$lastDump
  if [ -e "$path" ]; then
    echo "ДАМП $path УЖЕ СКАЧАН. ПРОПУСК"
  else
    scp storage:~/dumps/origin/$lastDump $path
    echo "УСПЕШНО СКАЧАН ФАЙЛ $lastDump to $path"
  fi