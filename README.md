# sneaker
![test](https://github.com/terakoya76/sneaker/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/terakoya76/sneaker/badge.svg?branch=master)](https://coveralls.io/github/terakoya76/sneaker?branch=master)

Cron execution schedule visualizer

> **Warning**
> This repository has been archived. Please use https://crontab-visualizer.terakoya76.dev/ instead.

## How to use
sneaker help you to find when job processes are triggered and when not

This example describes some job processes will be triggered on 0,5,7,10,15,20,25,27,30,40,45,47,50 minute on every hour, and we filter it by grep.

```bash
$ crontab -l | sneaker | grep Dec | grep 01,
Dec 01, 00H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 01H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 02H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 03H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 04H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 05H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 06H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 07H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 08H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 09H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 10H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 11H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 12H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 13H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 14H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 15H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 16H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 17H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 18H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 19H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 20H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 21H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 22H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
Dec 01, 23H: ■□□□□■□■□□■□□□□■□□□□■□□□□■□■□□■□□□□□□□□□■□□□□■□■□□■□□□□□□□□□
```

## Note
This tool currently only evaluate minute,hour,day,month part of cron expression.
So what it describes is not an actual execution schedule.
It contains some false-positive execution schedule.
