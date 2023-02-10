# dosee
Dosee Is a simple GO tool for Check Dosi Account by Session Lists

Created By viloid (github.com/xargs1)

> *** NOTE : USE AT YOUR OWN RISK! ***

> DONT SELL THIS SCRIPT! YOU'RE REALLY FUCKING POOR STUPID DOG!

<p align="center"> 
    <a href="https://goreportcard.com/report/github.com/xargs1/dosee">
        <img src="https://goreportcard.com/badge/github.com/xargs1/dosee">
    </a> 
    <a href="https://github.com/xargs1/dosee/issues">
        <img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat">
    </a> 
    <a href="https://twitter.com/verry__d">
        <img src="https://img.shields.io/twitter/follow/verry__d.svg?logo=twitter">
    </a> 
</p>

```bash
Dosee (Dosi Checker)

Is a simple GO tool for Check Dosi Account by Session Lists

Coded By : github.com/vsec7

Basic Usage :
 ▶ cat session_list.txt | dosee
 ▶ dosee < session_list.txt
Advanced Usage :
 ▶ cat session_list.txt | dosee -tg all -o result.txt

Options :
  -c, --conf <config.yaml>              Set file config.yaml (default: config.yaml)
  -tg, --telegram <all|active|expired>  Set Notification to Telegram
  -o, --output <file>                   Set Output File
```

## • Features
- Check dosi session by stdin
- Output : file & Telegram notification

## • Requirement
> go version: go1.18+ 

## • Installation
```bash
go install -v github.com/xargs1/dosee@latest
```

## • Configuration Template for telegram notification
```yaml
BOT_TOKEN: 606xxxx:AAFgmR6nDxxxxxxxxxxxxxxxxxxxxx
CHAT_ID: 173666xxxxxx
```

## • Donate

SOL Address : viloid.sol

BSC Address : 0xd3de361b186cc2Fc0C77764E30103F104a6d6D07
