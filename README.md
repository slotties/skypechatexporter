skypechatexporter
=================

Exports skype chat histories

Preparation
===========

Execute setup.bat or
	go get github.com/mattn/go-sqlite3
before you build it.	

Usage
=====
	
	skypeexport -chatname [what chat to export]

Optionally you can define a specific main.db if it's not placed in the same directory you execute skypeexport from:
	skypeepxort -chatname [...] -db /path/to/main.db