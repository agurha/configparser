## ConfigParser package 

## You can just give a file with some configuration

[default]
port = 1337

[prod]
redis :
	port = 6397
	password = 12345
	host = prodserver.com

[dev]
redis :
	port = 6397
	password = 12345
	host = localhost