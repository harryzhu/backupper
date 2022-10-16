# backupper

## Config

use [confctl](https://github.com/harryzhu/confctl) to edit conf.db

```
./confctl set --name=url_list --val=http://localhost/downloads.txt
./confctl set --name=dir_save_root --val=/Volumes/SSD2/temp
```

## Build
```
go build -o ./backupper main.go 
```

## Run
```
./backupper download
./backupper multidownload
```