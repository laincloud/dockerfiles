#!/bin/bash

Dirs="www/data module config tmp"
PermanentDir="/data"
AppDir="/var/www/html/zentaopms"
UserCfg="${PermanentDir}/config/my.php"
InstallFile="${AppDir}/www/install.php"
UpgradeFile="${AppDir}/www/upgrade.php"

[ ! -d $PermanentDir/www ] && mkdir $PermanentDir/www

for d in $Dirs
do
  if [ ! -d ${PermanentDir}/${d} ] ;then

    if [ "$d" == "www/data" ];then
      [ -d ${AppDir}/${d} ] && mkdir -pv ${PermanentDir}/www && mv ${AppDir}/${d} ${PermanentDir}/www
    fi

    [ -d ${AppDir}/${d} ] && mv ${AppDir}/${d} ${PermanentDir}/${d} || mkdir -pv ${PermanentDir}/${d}
  else
    mv ${AppDir}/${d} ${AppDir}/${d}.bak
  fi

  ln -s ${PermanentDir}/${d} ${AppDir}/${d}
done

if [ -f $UserCfg ];then
  [ -f $InstallFile ] && rm -f $InstallFile
  [ -f $UpgradeFile ] && rm -f $UpgradeFile
fi

# run apache
exec docker-php-entrypoint apache2-foreground